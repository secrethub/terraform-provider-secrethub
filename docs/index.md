---
layout: "secrethub"
page_title: "SecretHub Provider"
sidebar_current: "docs-secrethub-index"
description: |-
  The SecretHub provider is used to interact with the resources supported by SecretHub. The provider needs to be properly configured before it can be used.
---

# SecretHub Provider

The [SecretHub](https://www.secrethub.io) provider is used to interact with the
resources supported by SecretHub. The provider needs to be configured with a SecretHub credential before it can be used.

You can set environment variable `SECRETHUB_CREDENTIAL` or read it from disk using the `file()` function.

Use the navigation to the left to read about the available resources.

## Installation

### Terraform 0.13 and later

To install this provider, copy and paste this code into your Terraform configuration. Then, run `terraform init`.

```hcl
terraform {
  required_providers {
    secrethub = {
      source = "secrethub/secrethub"
      version = "1.2.2"
    }
  }
}
```

### Terraform 0.12 and earlier

For Terraform 0.12 and earlier, you have to download the binary for your operating system and architecture from [GitHub releases](https://github.com/secrethub/terraform-provider-secrethub/releases) yourselves and place it in `~/.terraform.d/plugins` (`%APPDATA%\terraform.d\plugins` on Windows).

For linux amd64, you can do so with the following command:
```sh
mkdir -p ~/.terraform.d/plugins && curl -SfL https://github.com/secrethub/terraform-provider-secrethub/releases/latest/download/terraform-provider-secrethub-linux-amd64.tar.gz | tar zxf - -C ~/.terraform.d/plugins
```

## Example Usage

```hcl
provider "secrethub" {
  credential = file("~/.secrethub/credential")
}
```

## Argument Reference

The following arguments are supported:

* `credential` - (Required) Credential to use for SecretHub authentication. Can also be sourced from `SECRETHUB_CREDENTIAL`.
* `credential_passphrase` - (Optional) Passphrase to unlock the authentication passed in `credential`.
