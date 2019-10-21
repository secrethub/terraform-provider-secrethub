package secrethub

import (
	"fmt"

	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceServiceAWS() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceAWSCreate,
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
			"kms_key_arn": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ARN of the KMS-key to be used for encrypting the service's account key.",
			},
			"role": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The role name or ARN of the IAM role that should have access to this service account.",
			},
		},
	}
}

func resourceServiceAWSCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	repo := d.Get("repo").(string)
	description := d.Get("description").(string)
	kmsKey := d.Get("kms_key_arn").(string)
	role := d.Get("role").(string)

	cfg := aws.NewConfig()

	kmsKeyARN, err := arn.Parse(kmsKey)
	if err != nil {
		return fmt.Errorf("the provider kms key is not a valid ARN: %s", err)
	}
	cfg = cfg.WithRegion(kmsKeyARN.Region)

	service, err := client.Services().Create(repo, description, credentials.CreateAWS(kmsKey, role, cfg))
	if err != nil {
		return err
	}

	d.SetId(service.ServiceID)

	return nil
}
