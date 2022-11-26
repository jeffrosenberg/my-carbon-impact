data "archive_file" "zip" {
  type        = "zip"
  source_file = "../go/bin/web/${var.entity}-${var.operation}"
  output_path = "../go/bin/web/${var.entity}-${var.operation}.zip"
}

data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  ent_oper = "${var.entity}-${var.operation}"
  name = "${var.entity}-${var.operation}-web"
}

resource "aws_lambda_function" "lambda" {
  # If the file is not in the current working directory you will need to include a 
  # path.module in the filename.
  filename      = "../go/bin/web/${local.ent_oper}.zip"
  function_name = local.name
  role          = var.iam_role_arn
  handler       = "main"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = data.archive_file.zip.output_base64sha256

  runtime = "go1.x"

  environment {
    variables = {
      zerolog_level = var.log_level
      region = data.aws_region.current.name
    }
  }

  tags = {
    entity = var.entity
    operation = var.operation
    tier = "web"
  }
}

resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${local.name}"
  retention_in_days = 90
}