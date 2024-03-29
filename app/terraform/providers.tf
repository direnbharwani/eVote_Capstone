terraform {
    required_version = ">= 1.5.5"

    required_providers {
      aws = {
        source = "hashicorp/aws"
        version = ">=5.40.0"
      }
    }
}

provider "aws" {
    region = var.aws_region        
}