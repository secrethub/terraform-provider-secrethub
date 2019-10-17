---
layout: "secrethub"
page_title: "Resource: secrethub_service_aws"
sidebar_current: "docs-secrethub-resource-service-aws"
description: |-
  Creates and manages service accounts tied to an AWS IAM role.
---

# Resource: secrethub_service_aws

This resource allows you to manage a service account that is tied to an AWS IAM role.

The native AWS identity provider uses a combination of AWS IAM and AWS KMS to provide access to SecretHub for any service running on AWS (e.g. EC2, Lambda or ECS) without needing a SecretHub credential.

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the service so others will recognize it.
* `kms_key` - (Required) The ARN of the KMS-key to be used for encrypting the service's account key.
* `repo` - (Required) The path of the repository on which the service operates.
* `role` - (Required) The role name or ARN of the IAM role that should have access to this service account.

