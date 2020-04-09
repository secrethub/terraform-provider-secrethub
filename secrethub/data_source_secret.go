package secrethub

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSecret() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSecretRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path where the secret is stored. To use a specific version, append the version number to the path, separated by a colon (path:version). Defaults to the latest version.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the secret.",
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The secret contents.",
			},
		},
	}
}

func dataSourceSecretRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("path").(string)

	secret, err := client.Secrets().Versions().GetWithData(path)
	if err != nil {
		return err
	}

	err = d.Set("value", string(secret.Data))
	if err != nil {
		return err
	}
	err = d.Set("version", secret.Version)
	if err != nil {
		return err
	}

	d.SetId(path)

	return nil
}
