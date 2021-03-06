---
layout: "secrethub"
page_title: "secrethub_access_rule"
sidebar_current: "docs-secrethub-resource-access-rule"
description: |-
  Creates and manages access rules.
---

# secrethub_access_rule Resource

This resource allows you to create and manage access rules, to give users and/or service accounts permissions on directories.

## Example Usage

```terraform
data "secrethub_dir" "demo" {
    path = "workspace/demo"
}

resource "secrethub_access_rule" "demo_app" {
  account_name = secrethub_service_aws.demo_app.id
  dir          = data.secrethub_dir.demo.path
  permission   = "read"
}
```

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) The name of the account (username or service ID) for which the permission holds.
* `dir` - (Required) The path of the directory on which the permission holds.
* `permission` - (Required) The permission that the account has on the given directory: read, write or admin

## Import

Access rules can be imported using the id: `<path>:<account_name>` e.g. `path/to/dir:alice`.
