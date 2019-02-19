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
				Description: "The path where the secret is stored.",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Overrides the `path_prefix` defined in the provider.",
			},
			"version": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The version of the secret. Defaults to the latest.",
			},
			"data": {
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

	path, err := getSecretPath(d, &provider)
	if err != nil {
		return err
	}

	remote, err := client.Secrets().Get(path)
	if err != nil {
		return err
	}

	version := d.Get("version").(int)
	if version == 0 {
		d.Set("version", remote.LatestVersion)
	}

	if d.Get("data") == "" || d.Get("version") != remote.LatestVersion {
		// Only fetch the secret contents if it hasn't been fetched before or if the version is out of sync
		updated, err := client.Secrets().Versions().GetWithData(path)
		if err != nil {
			return err
		}
		d.Set("data", string(updated.Data))
		d.Set("version", updated.Version)
	}

	d.SetId(string(path))

	return nil
}
