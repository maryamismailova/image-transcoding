module "infra" {
  source                     = "../../modules/infra"
  env                        = "dev"
  source_s3_bucket_name      = "rb-input-dev"
  destination_s3_bucket_name = "rb-output-dev"
  lambda_function_name       = "image-scaling-trigger-dev"
  lambda_memory_size         = 512

}
