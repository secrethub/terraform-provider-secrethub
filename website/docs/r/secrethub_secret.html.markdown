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
  path = "/ssh_key"
  data = "${file("/path/to/ssh/key")}."
}
```

To generate a new secret:

```hcl
resource "secrethub_secret" "db_password" {
  path = "/db_password"
  generate {
    length = 20
    symbols = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The path used for storing
 the secret.
* `data` - (Optional) The secret to store. Either `data` or `generate` must be specified.
* `generate` - (Optional) The settings block for autogenerating a secret. Either `data` or `generate` must be specified.

### Generate block

* `length` - (Optional) How many characters long the secret should be.
* `symbols` - (Optional) Specifies if symbols can be used to generate the secret.
* `force_new` - (Optional) Force a new secret generation at every run.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The current version of the secret.
