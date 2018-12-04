variable "aws_access_key" {}

variable "aws_secret_key" {}

variable "aws_region" {
  default = "us-east-1"
}

variable "environment" {
  default = "dev"
}

provider "aws" {
  access_key = "${var.aws_access_key}"
  secret_key = "${var.aws_secret_key}"
  region     = "${var.aws_region}"
}

provider "secrethub" {
  version      = "latest"
  organization = "myOrg"
  repository   = "myRepo"
}

resource "secrethub_generate" "db_password" {
  path    = "/${var.environment}/db/password"
  length  = 22
  symbols = false
}

resource "secrethub_write" "db_username" {
  path    = "/${var.environment}/db/username"
  data    = "dbUser"
}

resource "aws_db_instance" "default" {
  allocated_storage    = 10
  storage_type         = "gp2"
  engine               = "mysql"
  engine_version       = "5.7"
  instance_class       = "db.t2.micro"
  name                 = "mydb"
  username             = "${secrethub_write.db_username.data}"
  password             = "${secrethub_generate.db_password.data}"
  parameter_group_name = "default.mysql5.7"
}
