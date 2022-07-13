resource "aws_iam_role" "iam_for_lambda" {
  name = "iam_for_lambda"

  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Principal": {
        "Service": "lambda.amazonaws.com"
      },
      "Effect": "Allow",
      "Sid": ""
    }
  ]
}
EOF
}

module "rest_lambda" {
  source    = "./modules/rest-lambda"
  for_each  = var.entities
  entity    = each.value.entity
  operation = each.value.operation
  iam_role_arn = aws_iam_role.iam_for_lambda.arn
}