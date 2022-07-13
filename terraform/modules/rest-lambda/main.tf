data "archive_file" "zip" {
  type        = "zip"
  source_file = "../go/bin/${var.entity}-${var.operation}"
  output_path = "${var.entity}-${var.operation}.zip"
}

resource "aws_lambda_function" "lambda" {
  # If the file is not in the current working directory you will need to include a 
  # path.module in the filename.
  filename      = "${var.entity}-${var.operation}.zip"
  function_name = "${var.entity}-${var.operation}"
  role          = var.iam_role_arn
  handler       = "${var.entity}-${var.operation}"

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = data.archive_file.zip.output_base64sha256

  runtime = "go1.x"

  # environment {
  #   TODO
  # }

  tags = {
    entity = var.entity
    operation = var.operation
  }
}