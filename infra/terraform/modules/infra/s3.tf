module "s3_source" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.4.0"
  bucket  = var.source_s3_bucket_name
  acl     = "private"
  versioning = {
    enabled = false
  }
}

module "s3_destination" {
  source  = "terraform-aws-modules/s3-bucket/aws"
  version = "3.4.0"
  bucket  = var.destination_s3_bucket_name
  acl     = "private"
  versioning = {
    enabled = false
  }
}
