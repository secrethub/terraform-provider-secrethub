package secrethub

import "github.com/hashicorp/terraform/helper/schema"

func dataSourceDir() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDirRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path of the directory",
			},
		},
	}
}

func dataSourceDirRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("path").(string)

	_, err := client.Dirs().GetTree(path, 0, false)
	if err != nil {
		return err
	}

	d.SetId(path)

	return nil
}
