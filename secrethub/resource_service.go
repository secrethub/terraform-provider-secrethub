package secrethub

import (
	"github.com/secrethub/secrethub-go/internals/api"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceService() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceCreate,
		Read:   resourceServiceRead,
		Delete: resourceServiceDelete,
		Schema: map[string]*schema.Schema{
			"repo": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path of the repository on which the service operates.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A description of the service so others will recognize it.",
			},
			"credential": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The credential of the service account.",
			},
			"id": {
				Type: schema.TypeString,
				Computed: true,
				Description: "A unique identifier for the service.",
			},
		},
	}
}

func resourceServiceCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	repo := d.Get("repo").(string)
	description := d.Get("description").(string)

	credential := credentials.CreateKey()

	service, err := client.Services().Create(repo, description, credential)
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)

	exported, err := credential.Export()
	if err != nil {
		return err
	}

	err = d.Set("credential", string(exported))
	if err != nil {
		return err
	}

	err = d.Set("id", service.ServiceID)

	return nil
}

func resourceServiceRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	remote, err := client.Services().Get(d.Id())
	if err == api.ErrServiceNotFound {
		// The service account was deleted outside of the current Terraform workspace, so invalidate this resource
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("description", remote.Description)
	if err != nil {
		return err
	}

	return nil
}

func resourceServiceDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	_, err := client.Services().Delete(d.Id())
	if err != nil {
		return err
	}

	return nil
}
