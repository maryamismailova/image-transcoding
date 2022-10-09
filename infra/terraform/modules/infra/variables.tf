variable "source_s3_bucket_name" {
  default = "rb-input"
  type    = string
}

variable "destination_s3_bucket_name" {
  default = "rb-output"
  type    = string
}


variable "env" {
  type = string
}


# lambda
variable "lambda_function_name" {
  default = "image-scaling-trigger"
  type    = string
}

variable "lambda_handler_name" {
  default = "main"
  type    = string
}

variable "lambda_function_runtime" {
  default = "go1.x"
  type    = string
}

variable "lambda_memory_size" {
  default = 512
  type    = number
}

variable "lambda_architectures" {
  default = ["x86_64"]
  type    = list(string)
}
