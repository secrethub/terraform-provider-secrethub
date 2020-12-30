package secrethub

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/internals/api"
)

func resourceDir() *schema.Resource {
	return &schema.Resource{
		Create: resourceDirCreate,
		Read:   resourceDirRead,
		Delete: resourceDirDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDirImport,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path of the directory.",
			},
		},
	}
}

func resourceDirCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("path").(string)

	_, err := client.Dirs().Create(path)
	if err != nil {
		return err
	}

	d.SetId(path)

	return resourceDirRead(d, m)
}

func resourceDirRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	_, err := client.Dirs().GetTree(path, 0, true)
	if api.IsErrNotFound(err) {
		// The directory was deleted outside of the current Terraform workspace, so invalidate this resource
		d.SetId("")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error fetching directory: %s", err)
	}

	return nil
}

func resourceDirDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	tree, err := client.Dirs().GetTree(path, 1, false)
	if api.IsErrNotFound(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if len(tree.Dirs) > 1 || len(tree.Secrets) > 0 {
		return fmt.Errorf("cannot remove directory %s: it is not empty", path)
	}

	return client.Dirs().Delete(path)
}

func resourceDirImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	path := d.Id()

	err := api.ValidateDirPath(path)
	if err != nil {
		return nil, err
	}

	err = d.Set("path", path)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
