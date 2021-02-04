package secrethub

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/secrethub/secrethub-go/pkg/secrethub"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"
)

var version string

// Provider returns the SecretHub Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECRETHUB_CREDENTIAL", nil),
				Description: "Credential to use for SecretHub authentication. Can also be sourced from SECRETHUB_CREDENTIAL.",
			},
			"credential_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECRETHUB_CREDENTIAL_PASSPHRASE", nil),
				Description: "Passphrase to unlock the authentication passed in `credential`. Can also be sourced from SECRETHUB_CREDENTIAL_PASSPHRASE.",
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"secrethub_secret":      resourceSecret(),
			"secrethub_access_rule": resourceAccessRule(),
			"secrethub_service":     resourceService(),
			"secrethub_service_aws": resourceServiceAWS(),
			"secrethub_service_gcp": resourceServiceGCP(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"secrethub_secret": dataSourceSecret(),
		},
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	credRaw := d.Get("credential").(string)
	passphrase := d.Get("credential_passphrase").(string)

	options := []secrethub.ClientOption{
		secrethub.WithAppInfo(&secrethub.AppInfo{
			Name:    "terraform-provider-secrethub",
			Version: version,
		}),
	}

	if credRaw != "" {
		keyProvider := credentials.UseKey(credentials.FromString(credRaw))
		var provider credentials.Provider = keyProvider
		if passphrase != "" {
			provider = keyProvider.Passphrase(credentials.FromString(passphrase))
		}
		options = append(options, secrethub.WithCredentials(provider))
	}

	client, err := secrethub.NewClient(options...)
	if err != nil {
		return nil, err
	}

	return providerMeta{client}, nil
}

type providerMeta struct {
	client *secrethub.Client
}
