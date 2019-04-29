package secrethub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/internals/api"
	"github.com/secrethub/secrethub-go/pkg/randchar"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecretImport,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path where the secret will be stored.",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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

	synchronizePathPrefix(d, &provider)
	path := getSecretPath(d)

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

	synchronizePathPrefix(d, &provider)

	path := d.Id()
	remote, err := client.Secrets().Get(path)
	if err == api.ErrSecretNotFound {
		// The secret was deleted outside of the current Terraform workspace, so invalidate this resource
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	prev := d.Get("version")
	if prev != remote.LatestVersion {
		// The secret has been updated outside of the current Terraform workspace, so the new secret version has to be fetched
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

func resourceSecretImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	path := d.Id()

	provider := m.(providerMeta)
	if provider.pathPrefix != "" && !strings.HasPrefix(path, provider.pathPrefix) {
		return nil, fmt.Errorf("secret import path must be absolute")
	}

	err := api.ValidateSecretPath(path)
	if err != nil {
		return nil, err
	}

	if provider.pathPrefix != "" {
		relativePath := strings.TrimPrefix(path, provider.pathPrefix)
		path = trimPathComponent(relativePath)
	}

	d.Set("path", path)

	return []*schema.ResourceData{d}, nil
}

// getSecretPath finds the full path of a secret, combining the (relative) path with path prefix
func getSecretPath(d *schema.ResourceData) string {
	return newCompoundSecretPath(d.Get("path_prefix").(string), d.Get("path").(string))
}

// synchronizeSecretPath synchronizes the path prefix state of the resource with the fallback path prefix set on the provider
func synchronizePathPrefix(d *schema.ResourceData, provider *providerMeta) {
	if _, ok := d.GetOk("path_prefix"); !ok {
		// Fall back to the provider prefix
		d.Set("path_prefix", provider.pathPrefix)
	}
}

const pathSeparator = "/"

// newCompoundSecretPath returns a SecretPath that combines multiple path components into a single secret path
func newCompoundSecretPath(components ...string) string {
	var processed []string
	for _, c := range components {
		trimmed := trimPathComponent(c)
		if trimmed != "" {
			processed = append(processed, trimmed)
		}
	}
	return strings.Join(processed, pathSeparator)
}

func trimPathComponent(c string) string {
	return strings.Trim(c, pathSeparator)
}
