---
layout: "secrethub"
page_title: "SecretHub: secrethub_write"
sidebar_current: "docs-secrethub-resource-write"
description: |-
  Writes a secret at a given path.
---

# secrethub_write

This resource allows to write secrets at a given path, if the path is exists already then the resource will write a new version of it.

## Example Usage

```hcl
resource "secrethub_write" "ssh_key" {
  path = "/ssh_key"
  data = "${file("/path/to/ssh/key")}."
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) The path used for storing the secret.
* `data` - (Required) The secret to store.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `version` - The current version of the secret.