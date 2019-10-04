---
layout: "secrethub"
page_title: "Resource: secrethub_secret"
sidebar_current: "docs-secrethub-resource-secret"
description: |-
  Writes a secret at a given path.
---

# secrethub_secret

This resource allows to write secrets at a given path, if the path is exists already then the resource will write a new version of it.

## Example Usage

To write a secret:

```hcl
resource "secrethub_secret" "ssh_key" {
  path = "company/repo/ssh_key"
  value = "${file("/path/to/ssh/key")}"
}
```

To generate a new secret:

```hcl
resource "secrethub_secret" "db_password" {
  path = "company/repo/db_password"

  generate {
    length = 20
    use_symbols = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The path where the secret will be stored.
* `value` - (Optional) The secret contents. Either `value` or `generate` must be defined.
* `generate` - (Optional) Settings for autogenerating a secret. Either `value` or `generate` must be defined.

Nested `generate` blocks have the following structure:

* `length` - (Required) The length of the secret to generate.
* `use_symbols` - (Optional) Whether the secret should contain symbols.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The version of the secret.
