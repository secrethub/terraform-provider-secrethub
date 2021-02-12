---
layout: "secrethub"
page_title: "secrethub_dir"
sidebar_current: "docs-secrethub-resource-dir"
description: |-
  Creates a directory at a given path.
---

# secrethub_dir Resource

This resource allows you to create a directory at a given path.

## Example Usage

```terraform
data "secrethub_dir" "project" {
    path = "company/project"
}

resource "secrethub_dir" "environment" {
    path = "${data.secrethub_dir.project.path}/${var.environment}"
}

resource "secrethub_dir" "database" {
    path = "${secrethub_dir.environment.path}/db"
}

resource "secrethub_secret" "db_user" {
    path = "${secrethub_dir.database.path}/username"
    value = "db_user"
}

resource "secrethub_secret" "db_password" {
    path = "${secrethub_dir.database.path}/password"
    generate {
        length = 32
    }
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The path of the directory.
* `force_destroy` - (Optional) Whether to allow deleting this directory if it's not empty. When set to `false`, you'll get an error when trying to delete the directory if it still contains directories or secrets.
