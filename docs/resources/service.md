---
layout: "secrethub"
page_title: "secrethub_service"
sidebar_current: "docs-secrethub-resource-service"
description: |-
  Creates and manages service accounts
---

# secrethub_service Resource

This resource allows you to manage a service account - an account for machines.

## Example Usage

```terraform
resource "secrethub_service" "demo_service_account" {
  repo = "workspace/repo"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the service so others will recognize it.
* `repo` - (Required) The path of the repository on which the service operates.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `credential` - The credential of the service account.
