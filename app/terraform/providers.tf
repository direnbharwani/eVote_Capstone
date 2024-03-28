terraform {
  required_version = ">= 1.5.5"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 5.15.0"
    }
  }

  backend "s3" {
    bucket         = "evote-poc-terraform-state"
    region         = "ap-southeast-1"
    key            = "state/terraform.tfstate"
    encrypt        = true
    dynamodb_table = "terraform-locks"
  }
}

provider "aws" {
  region = var.aws_region
}