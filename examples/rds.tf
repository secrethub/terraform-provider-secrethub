variable "aws_access_key" {}

variable "aws_secret_key" {}

variable "aws_region" {
  default = "us-east-1"
}

variable "environment" {
  default = "dev"
}

provider "aws" {
  access_key = var.aws_access_key
  secret_key = var.aws_secret_key
  region     = var.aws_region
}

provider "secrethub" {
  credential = file("~/.secrethub/credential")
}

data "secrethub_dir" "repo" {
    path = "company/repo"
}

resource "secrethub_dir" "environment" {
    path = "${data.secrethub_dir.repo.path}/${var.environment}"
}

resource "secrethub_secret" "db_password" {
  path = "${secrethub_dir.environment.path}/db/password"

  generate {
    length      = 22
    use_symbols = true
  }
}

resource "secrethub_secret" "db_username" {
  path  = "${secrethub_dir.environment.path}/db/username"
  value = "mysqluser"
}

resource "aws_db_instance" "default" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = secrethub_secret.db_username.value
  password             = secrethub_secret.db_password.value
  parameter_group_name = "default.mysql5.7"
}
