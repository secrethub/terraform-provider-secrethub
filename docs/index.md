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
