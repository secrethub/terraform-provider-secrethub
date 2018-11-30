---
layout: "secrethub"
page_title: "SecretHub: secrethub_read"
sidebar_current: "docs-secrethub-datasource-read"
description: |-
  Read a secret
---

# Data Source: secrethub_read

Use this data source to read secrets already in SecretHub

## Example Usage

```hcl
data "secrethub_read" "db_password" {
  path = "/db"
}
```

## Argument Reference

* `path` - (Required) The path where the secret is stored
* `version` - (Optional) The current secret version

## Attributes Reference

* `data` - The secret retrieved
* `version` - The version of the secret read