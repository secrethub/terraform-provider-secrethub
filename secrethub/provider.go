package secrethub

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/keylockerbv/secrethub-go/pkg/secrethub"
)

// Provider returns the ScretHub Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Credential to use for SecretHub authentication.",
			},
			"credential_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Passphrase to unlock the authentication passed in `credential`.",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default path prefix of the secret resources and data sources. If left blank, every secret requires the path to be absolute (namespace/repository[/dir]/secret_name).",
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"secrethub_secret": resourceSecret(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	credRaw := d.Get("credential").(string)
	parser := secrethub.NewCredentialParser(secrethub.DefaultCredentialDecoders)
	parsed, err := parser.Parse(credRaw)
	if err != nil {
		return nil, err
	}

	cred, err := parsed.Decode()
	if err != nil {
		return nil, err
	}

	client := secrethub.NewClient(cred, nil)
	pathPrefix := d.Get("path_prefix").(string)
	return providerMeta{&client, pathPrefix}, nil
}

type providerMeta struct {
	client     *secrethub.Client
	pathPrefix string
}
