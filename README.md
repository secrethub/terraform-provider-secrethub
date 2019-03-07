# Terraform SecretHub Provider

> [SecretHub](https://secrethub.io) is a developer tool to help you keep database passwords, API tokens, and other secrets out of IT automation scripts.

The Terraform SecretHub Provider lets you manage your secrets using Terraform.

<br>

<p align="center">
  <img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" alt="Terraform" width="330px">
  <img width="50px"/>
  <img src="https://secrethub.io/img/secrethub-logo.png" alt="SecretHub" width="360px">
</p>


## Installation

Download and extract the [latest release](https://github.com/keylockerbv/terraform-provider-secrethub/releases/latest) of the Terraform SecretHub Provider and move it to your [Terraform plugin directory](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins) (`~/.terraform.d/plugins/`, or `%APPDATA%\terraform.d\plugins` on Windows).

Afterwards, you can use `terraform init` as you would normally.

## Usage

```hcl
provider "secrethub" {
  # pass in credential or set SECRETHUB_CREDENTIAL enviroment variable
  credential  = "${file("~/.secrethub/credential")}" 
  path_prefix = "my-org/my-repo"
}

resource "secrethub_secret" "db_password" {
  path = "db/password"

  generate {
    length  = 22
    symbols = true
  }
}

resource "secrethub_secret" "db_username" {
  path = "db/username"
  data = "db-user"
}

resource "aws_db_instance" "default" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "${secrethub_secret.db_username.data}"
  password             = "${secrethub_secret.db_password.data}"
  parameter_group_name = "default.mysql5.7"
}
```

Have a look at the [reference docs](/website/docs) for more information on the supported resources and data sources.

Have a look at the [examples](/examples) for more practical use-cases.

## Development

### Building

Go get the source code:

```
go get -u https://github.com/keylockerbv/terraform-provider-secrethub
```

Build it using:

```
make build
```

### Testing

To run the [acceptance tests](https://www.terraform.io/docs/extend/testing/acceptance-tests/index.html), the following environment variables need to be setup.

* `SECRETHUB_CREDENTIAL` - 
* `SECRETHUB_TF_ACC_NAMESPACE` - a namespace registered on SecretHub. Make sure `SECRETHUB_CREDENTIAL` has admin access.
* `SECRETHUB_TF_ACC_REPOSITORY` - a repository within `SECRETHUB_TF_ACC_NAMESPACE` to be used by the acceptance tests. Make sure `SECRETHUB_CREDENTIAL` has admin access.

With the environment variables properly set up, run:

```
make testacc
```
