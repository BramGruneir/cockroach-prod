variable "aws_access_key" {}
variable "aws_secret_key" {}

variable "sql_logic_instances" {
  default = 0
}

# AWS region to use. WARNING: changing this will break the AMI ID.
variable "aws_region" {
  default = "us-east-1"
}

# AWS availability zone. Make sure it exists for your account.
variable "aws_availability_zone" {
  default = "us-east-1a"
}

# AWS image ID. The default is valid for region "us-east-1".
variable "aws_ami_id" {
  default = "ami-408c7f28"
}

# Path to the cockroach repository.
variable "cockroach_repo" {
  default = "../../../../cockroach"
}

# Path to the sqllogictest repository.
variable "sqllogictest_repo" {
  default = "../../../../sqllogictest"
}

# Name of the ssh key pair for this AWS region. Your .pem file must be:
# ~/.ssh/<key_name>.pem
variable "key_name" {
  default = "cockroach"
}
