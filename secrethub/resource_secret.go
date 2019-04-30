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
			"value": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"generate"},
				Description:   "The secret contents. Either `value` or `generate` must be defined.",
			},
			"generate": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Description:   "Settings for autogenerating a secret. Either `value` or `generate` must be defined.",
				ConflictsWith: []string{"value"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secret to generate.",
						},
						"use_symbols": {
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

	valueStr := d.Get("value").(string)
	generateList := d.Get("generate").([]interface{})
	if valueStr == "" && len(generateList) == 0 {
		return fmt.Errorf("either 'value' or 'generate' must be specified")
	}

	var value []byte
	if valueStr != "" {
		value = []byte(valueStr)
	}

	if len(generateList) > 0 {
		settings := generateList[0].(map[string]interface{})
		useSymbols := settings["use_symbols"].(bool)
		length := settings["length"].(int)
		var err error
		value, err = randchar.NewGenerator(useSymbols).Generate(length)
		if err != nil {
			return err
		}
	}

	path := getSecretPath(d, &provider)

	res, err := client.Secrets().Write(path, value)
	if err != nil {
		return err
	}

	d.SetId(string(path))
	d.Set("value", string(value))
	d.Set("version", res.Version)

	return resourceSecretRead(d, m)
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

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
		d.Set("value", string(updated.Data))
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
