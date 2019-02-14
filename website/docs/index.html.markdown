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
provider "secrethub" {
  credential = "${file("~/.secrethub/credential")}"
  path_prefix = "my_org/my_repo"
}

resource "secrethub_secret" "db_password" {
  path = "db_password"
  data = "mypassword"
}
```

## Argument Reference

The following arguments are supported:

* `credential` - (Required) Credential to use for SecretHub authentication.
* `credential_passphrase` - (Optional) Passphrase to unlock the authentication passed in `credential`.
* `path_prefix` - (Optional) The default path prefix of the secret resources and data sources. Specifying it will reduce redundancy in the secret resources as it enables the use of relative paths.
