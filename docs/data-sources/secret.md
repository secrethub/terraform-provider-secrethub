---
layout: "secrethub"
page_title: "secrethub_secret"
sidebar_current: "docs-secrethub-datasource-secret"
description: |-
  Read a secret
---

# secrethub_secret Data Source

Use this data source to read secrets already in SecretHub

## Example Usage

```terraform
data "secrethub_secret" "db_password" {
  path = "company/repo/db/password"
}
```

## Argument Reference

* `path` - (Required) The path where the secret is stored. To use a specific version, append the version number to the path, separated by a colon (path:version). Defaults to the latest version.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `value` - The secret contents.
* `version` - The version of the secret.
