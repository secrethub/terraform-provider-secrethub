package secrethub

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/secrethub/secrethub-go/pkg/secrethub"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"
)

var version string

// Provider returns the ScretHub Terraform provider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"credential": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECRETHUB_CREDENTIAL", nil),
				Description: "Credential to use for SecretHub authentication. Can also be sourced from SECRETHUB_CREDENTIAL.",
			},
			"credential_passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SECRETHUB_CREDENTIAL_PASSPHRASE", nil),
				Description: "Passphrase to unlock the authentication passed in `credential`. Can also be sourced from SECRETHUB_CREDENTIAL_PASSPHRASE.",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default value to prefix path values with. If set, paths for resources and data sources will be prefixed with the given prefix, allowing you to use relative paths instead. If left blank, every path must be absolute (namespace/repository/[dir/]secret_name).",
				Deprecated:  "Deprecated in favor of Terraform's native variables",
			},
		},
		ConfigureFunc: configureProvider,
		ResourcesMap: map[string]*schema.Resource{
			"secrethub_secret":      resourceSecret(),
			"secrethub_access_rule": resourceAccessRule(),
			"secrethub_service":     resourceService(),
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
		options = append(options, secrethub.WithCredentials(credentials.UseKey(credentials.FromString(credRaw)).Passphrase(credentials.FromString(passphrase))))
	}

	client, err := secrethub.NewClient(options...)
	if err != nil {
		return nil, err
	}

	pathPrefix := d.Get("path_prefix").(string)
	return providerMeta{client, pathPrefix}, nil
}

type providerMeta struct {
	client     *secrethub.Client
	pathPrefix string
}
