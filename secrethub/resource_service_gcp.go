package secrethub

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"
)

func resourceServiceGCP() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceGCPCreate,
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
			"service_account_email": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The email of the Google Service Account that provides the identity of the SecretHub service account.",
			},
			"kms_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The Resource ID of the Cloud KMS key to use to encrypt and decrypt your SecretHub key material.",
			},
		},
	}
}

func resourceServiceGCPCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	repo := d.Get("repo").(string)
	description := d.Get("description").(string)
	kmsKey := d.Get("kms_key_id").(string)
	serviceAccount := d.Get("service_account_email").(string)

	service, err := client.Services().Create(repo, description, credentials.CreateGCPServiceAccount(serviceAccount, kmsKey))
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)

	return resourceServiceRead(d, m)
}
