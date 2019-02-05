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
  path = "/db"
}
```

## Argument Reference

* `path` - (Required) The path where the secret is stored, optionally including a version number.

## Attributes Reference

* `data` - The secret retrieved.
* `version` - The version of the secret read.
