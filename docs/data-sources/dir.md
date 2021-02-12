---
layout: "secrethub"
page_title: "secrethub_dir"
sidebar_current: "docs-secrethub-datasource-dir"
description: |-
  Read a directory
---

# secrethub_dir Data Source

Use this data source to read a directory already in SecretHub.

## Example Usage

```terraform
data "secrethub_dir" "db" {
    path = "company/project/${var.environment}/db"
}

resource "secrethub_secret" "db_password" {
    path = "${data.secrethub_dir.db.path}/password"
    generate {
        length = 32
    }
}
```

## Argument Reference

* `path` - (Required) The path of the directory.
