---
layout: "secrethub"
page_title: "Data Source: secrethub_secret"
sidebar_current: "docs-secrethub-datasource-secret"
description: |-
  Read a secret
---

# Data Source: secrethub_read

Use this data source to read secrets already in SecretHub

## Example Usage

```hcl
data "secrethub_secret" "db_password" {
  path = "db-password"
}
```

## Argument Reference

* `path` - (Required) The path where the secret is stored.
* `path_prefix` - (Optional) Overrides the `path_prefix` defined in the provider.
* `version` - (Optional) The version of the secret read. Defaults to the latest.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `data` - The secret contents.
