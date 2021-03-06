---
layout: "secrethub"
page_title: "secrethub_secret"
sidebar_current: "docs-secrethub-resource-secret"
description: |-
  Writes a secret at a given path.
---

# secrethub_secret Resource

This resource allows you to write secrets at a given path, if the path already exists then the resource will write a new version of it.

## Example Usage

To write a secret:

```terraform
data "secrethub_dir" "repo" {
    path = "company/repo"
}

resource "secrethub_secret" "ssh_key" {
  path  = "${data.secrethub_dir.repo.path}/ssh_key"
  value = file("/path/to/ssh/key")
}
```

To generate a new, 20 characters long secret made of alphanumeric characters:

```terraform
data "secrethub_dir" "repo" {
    path = "company/repo"
}

resource "secrethub_secret" "db_password" {
  path = "${data.secrethub_dir.repo.path}/db_password"

  generate {
    length   = 20
  }
}
```

To generate a new secret made of lowercase letters and symbols, with minimum 5 symbols:

```terraform
data "secrethub_dir" "repo" {
    path = "company/repo"
}

resource "secrethub_secret" "db_password" {
  path = "${data.secrethub_dir.repo.path}/db_password"

  generate {
    length   = 20
    charsets = ["lowercase", "symbols"]
    min      = {
        symbols = 5
    }
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
* `charsets` - (Optional) List of charset names defining the set of characters to randomly generate a secret from. The supported charsets are: all, alphanumeric, numeric, lowercase, uppercase, letters, symbols and human-readable. Defaults to alphanumeric.
* `min` - (Optional) A map defining lower bounds on the number of characters to use from any specific charsets.

~> Adding constraints reduces the strength of the secret. When possible avoid adding any constraints.
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The version of the secret.
