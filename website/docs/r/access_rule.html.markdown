---
layout: "secrethub"
page_title: "Resource: secrethub_access_rule"
sidebar_current: "docs-secrethub-resource-access-rule"
description: |-
  Creates and manages access rules.
---

# Resource: secrethub_access_rule

This resource allows you to create and manage access rules, to give users and/or service accounts permissions on directories.

## Argument Reference

The following arguments are supported:

* `account_name` - (Required) The name of the account (username or service ID) for which the permission holds.
* `dir` - (Required) The path of the directory on which the permission holds.
* `permission` - (Required) The permission that the account has on the given directory: read, write or admin
