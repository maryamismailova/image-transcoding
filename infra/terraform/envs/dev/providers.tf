terraform {
  backend "s3" {
    bucket  = "rb-project-maryam"
    region  = "eu-central-1"
    encrypt = true
    key     = "terraform/dev/terraform.tfstate"
  }

  required_version = "~>1.2"

  required_providers {
    aws = {
      source = "hashicorp/aws"
    }
  }
}

provider "aws" {
  region = "eu-central-1"
}
