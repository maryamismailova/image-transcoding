module "lambda" {
  source  = "terraform-aws-modules/lambda/aws"
  version = "4.0.2"

  function_name = var.lambda_function_name
  handler       = var.lambda_handler_name
  runtime       = var.lambda_function_runtime
  architectures = var.lambda_architectures
  memory_size   = var.lambda_memory_size

  create_package         = false
  local_existing_package = "${path.module}/empty.zip"

  create_unqualified_alias_allowed_triggers = false

  timeout             = 10
  attach_policy_json  = true
  attach_policy_jsons = true
  policy_json = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = [
          "s3:GetObject"
        ]
        Effect = "Allow"
        Resource = [
          "${module.s3_source.s3_bucket_arn}/*"
        ]
      },
      {
        Action = [
          "s3:PutObject"
        ]
        Effect = "Allow"
        Resource = [
          "${module.s3_destination.s3_bucket_arn}/*"
        ]
      }
    ]
  })

  environment_variables = {
    ENV = var.env
  }

}

resource "aws_s3_bucket_notification" "aws-lambda-trigger" {
  bucket = module.s3_source.s3_bucket_id
  lambda_function {
    lambda_function_arn = module.lambda_alias.lambda_alias_arn
    events              = ["s3:ObjectCreated:*"]

  }
}

module "lambda_alias" {
  source  = "terraform-aws-modules/lambda/aws//modules/alias"
  version = "4.0.2"

  create             = true
  use_existing_alias = false
  name               = var.env
  function_name      = module.lambda.lambda_function_arn

  allowed_triggers = {
    S3CreatedObject = {
      service    = "s3"
      source_arn = module.s3_source.s3_bucket_arn
    }
  }

}
