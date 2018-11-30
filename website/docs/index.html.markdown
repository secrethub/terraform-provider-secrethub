---
layout: "secrethub"
page_title: "Provider: SecretHub"
sidebar_current: "docs-secrethub-index"
description: |-
  The SecretHub provider is used to interact with the resources supported by SecretHub. The provider needs to be properly configured before it can be used.
---

# SecretHub Provider

The [SecretHub](https://www.secrethub.io) provider is used to interact with the
resources supported by SecretHub. The provider needs to be properly configured before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the SecretHub provider
provider "secrethub" {
  version = "latest"
  organization = "myOrg"
  repository = "myRepo"
}

# Generate random secret
resource "secrethub_generate" "db_password" {
  # ...
}

# Write a new secret
resource "secrethub_write" "api_key" {
  # ...
}
```

## Argument Reference

The following arguments are supported:

* `version` - (Optional) The provider version. Default value: `latest`.
* `organization` - (Required) The organization to use.
* `repository` - (Required) The repository to use.
* `config_dir` - (Optional) The directory where to find the SecretHub client configuration. Conflicts with `credential`. Default value: `~/.secrethub`
* `credential` - (Optional) Specify the encrypted credentials to use for authentication. Conflicts with `config_dir`.
* `credential_passphrase` - (Optional) Passphrase required to unlock the authentication passed in `credential`.