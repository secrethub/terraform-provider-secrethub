---
layout: "secrethub"
page_title: "secrethub_service_gcp"
sidebar_current: "docs-secrethub-resource-service-gcp"
description: |-
  Creates and manages SecretHub service accounts tied to a GCP Service Account.
---

# secrethub_service_gcp Resource

This resource allows you to manage a SecretHub service account that is tied to a GCP Service Account.

The GCP identity provider uses a combination of Cloud IAM and Cloud KMS to read secrets from SecretHub from any app running on Google Cloud (GCE, GKE, etc.) without needing to manage another key.

## GCP Project link

Before you can use this resource, you first have to link your SecretHub namespace with your GCP project.
You only have to do this once for your namespace and GCP project.

Because the linking process uses OAuth and therefore needs a web browser login, it cannot be Terraformed and needs the SecretHub CLI:

```
secrethub service gcp link <namespace> <project id>
```

## Argument Reference

The following arguments are supported:

* `service_account_email` - (Required) The email of the Google Service Account that provides the identity of the SecretHub service account.
* `kms_key_id` - (Required) The Resource ID of the Cloud KMS key to use to encrypt and decrypt your SecretHub key material.
* `repo` - (Required) The path of the repository on which the service operates.
* `description` - (Optional) A description of the service so others will recognize it.
