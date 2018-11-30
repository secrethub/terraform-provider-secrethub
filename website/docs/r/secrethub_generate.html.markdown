---
layout: "secrethub"
page_title: "SecretHub: secrethub_generate"
sidebar_current: "docs-secrethub-resource-generate"
description: |-
  Creates a random secret.
---

# secrethub_generate

This resource allows the creation of random secrets, if a path already used is specified then the resource will generate a new version of it.

## Example Usage

```hcl
resource "secrethub_generate" "db_password" {
  path = "/database"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The path used for storing the secret
* `length` - (Optional) How many characters long the secret should be
* `symbols` - (Optional) Specifies if symbols can be used to generate the secret
* `force_new` - (Optional) Force a new secret generation at every run

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - The secret generated
* `version` - The current version of the secret generated