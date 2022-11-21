terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.16"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region  = "us-west-2"
  profile = "terraform"

  default_tags {
    tags = {
      project = "my-carbon-impact"
      env = var.env
    }
  }
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}