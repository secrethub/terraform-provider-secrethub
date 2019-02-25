package secrethub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/keylockerbv/secrethub-go/pkg/randchar"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path where the secret will be stored.",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Overrides the `path_prefix` defined in the provider.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the secret.",
			},
			"data": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"generate"},
				Description:   "The secret contents. Either `data` or `generate` must be defined.",
			},
			"generate": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Description:   "Settings for autogenerating a secret. Either `data` or `generate` must be defined.",
				ConflictsWith: []string{"data"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secret to generate.",
						},
						"symbols": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Whether the secret should contain symbols.",
						},
					},
				},
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	dataStr := d.Get("data").(string)
	generateList := d.Get("generate").([]interface{})
	if dataStr == "" && len(generateList) == 0 {
		return fmt.Errorf("either 'data' or 'generate' must be specified")
	}

	var data []byte
	if dataStr != "" {
		data = []byte(dataStr)
	}

	if len(generateList) > 0 {
		settings := generateList[0].(map[string]interface{})
		symbols := settings["symbols"].(bool)
		length := settings["length"].(int)
		var err error
		data, err = randchar.NewGenerator(symbols).Generate(length)
		if err != nil {
			return err
		}
	}

	path := getSecretPath(d, &provider)

	res, err := client.Secrets().Write(path, data)
	if err != nil {
		return err
	}

	d.SetId(string(path))
	d.Set("data", string(data))
	d.Set("version", res.Version)

	return resourceSecretRead(d, m)
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	remote, err := client.Secrets().Get(path)
	if err != nil {
		return err
	}

	prev := d.Get("version")
	if prev != remote.LatestVersion {
		// The secret has been updated outside of the current terraform scope, so the new secret version has to be fetched
		updated, err := client.Secrets().Versions().GetWithData(path)
		if err != nil {
			return err
		}
		d.Set("data", string(updated.Data))
		d.Set("version", updated.Version)
	}

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSecretCreate(d, m)
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	client.Secrets().Delete(path)

	return nil
}

// getSecretPath finds the full path of a secret, combining the specified path with the provider's path prefix
func getSecretPath(d *schema.ResourceData, provider *providerMeta) string {
	prefix := d.Get("path_prefix").(string)
	if prefix == "" {
		// Fall back to the provider prefix
		prefix = provider.pathPrefix
	}
	pathStr := d.Get("path").(string)
	return newCompoundSecretPath(prefix, pathStr)
}

const pathSeparator = "/"

// newCompoundSecretPath returns a SecretPath that combines multiple path components into a single secret path
func newCompoundSecretPath(components ...string) string {
	var processed []string
	for _, c := range components {
		trimmed := strings.Trim(c, pathSeparator)
		if trimmed != "" {
			processed = append(processed, trimmed)
		}
	}
	return strings.Join(processed, pathSeparator)
}
