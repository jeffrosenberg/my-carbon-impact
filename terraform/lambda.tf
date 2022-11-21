resource "aws_iam_role" "mci_lambda" {
  name = "mci-lambda"

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

resource "aws_iam_role_policy_attachment" "terraform_lambda_policy" {
  role       = aws_iam_role.mci_lambda.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSLambdaBasicExecutionRole"
}

resource "aws_iam_policy" "lambda_data_access" {
  name        = "mci-data-access-policy"
  description = "Provides DynamoDb access to the mci-lambda role"

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "DynamoDBIndexAndStreamAccess",
            "Effect": "Allow",
            "Action": [
                "dynamodb:GetShardIterator",
                "dynamodb:Scan",
                "dynamodb:Query",
                "dynamodb:DescribeStream",
                "dynamodb:GetRecords",
                "dynamodb:ListStreams"
            ],
            "Resource": [
                "${aws_dynamodb_table.profiles.arn}/index/*",
                "${aws_dynamodb_table.profiles.arn}/stream/*"
            ]
        },
        {
            "Sid": "DynamoDBTableAccess",
            "Effect": "Allow",
            "Action": [
                "dynamodb:BatchGetItem",
                "dynamodb:BatchWriteItem",
                "dynamodb:ConditionCheckItem",
                "dynamodb:PutItem",
                "dynamodb:DescribeTable",
                "dynamodb:DeleteItem",
                "dynamodb:GetItem",
                "dynamodb:Scan",
                "dynamodb:Query",
                "dynamodb:UpdateItem"
            ],
            "Resource": "${aws_dynamodb_table.profiles.arn}"
        }
    ]
}
EOF
}

resource "aws_iam_policy_attachment" "lambda_data_access" {
  name       = "lambda-data-access-attachment"
  roles      = [aws_iam_role.mci_lambda.name]
  policy_arn = aws_iam_policy.lambda_data_access.arn
}

module "rest_lambda" {
  source    = "./modules/rest-lambda"
  for_each  = var.entities
  entity    = each.value.entity
  operation = each.value.operation
  iam_role_arn = aws_iam_role.mci_lambda.arn
  log_level = var.log_level
}