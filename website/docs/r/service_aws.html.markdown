---
layout: "secrethub"
page_title: "Resource: secrethub_service_aws"
sidebar_current: "docs-secrethub-resource-service-aws"
description: |-
  Creates and manages service accounts tied to an AWS IAM role.
---

# Resource: secrethub_service_aws

This resource allows you to manage a service account that is tied to an AWS IAM role.

The AWS identity provider uses a combination of AWS IAM and AWS KMS to read secrets from SecretHub from any app running on AWS (EC2, ECS, Lambda, etc.) without needing to manage another key.

## Example Usage

```terraform
resource "secrethub_service_aws" "your_application" {
  repo        = "workspace/repo"
  role        = "${aws_iam_role.your_application.name}"
  kms_key_arn = "${aws_kms_key.secrethub_e2e.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description of the service so others will recognize it.
* `kms_key_arn` - (Required) The ARN of the KMS-key to be used for encrypting the service's account key.
* `repo` - (Required) The path of the repository on which the service operates.
* `role` - (Required) The role name or ARN of the IAM role that should have access to this service account.

## See also

- [AWS Integration](https://secrethub.io/docs/reference/aws/)
- [AWS EC2 Guide](https://secrethub.io/docs/guides/aws-ec2/)
- [AWS ECS Guide](https://secrethub.io/docs/guides/aws-ecs/)
