locals {
  stage = "dev"
}

resource "aws_api_gateway_rest_api" "api" {
  name = "My Carbon Impact API"
}

resource "aws_api_gateway_resource" "profile" {
  path_part   = "profile"
  parent_id   = aws_api_gateway_rest_api.api.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.api.id
}

resource "aws_api_gateway_resource" "profile_id" {
  path_part   = "{id}"
  parent_id   = aws_api_gateway_resource.profile.id
  rest_api_id = aws_api_gateway_rest_api.api.id
}

# resource "aws_api_gateway_model" "profile" {
#   rest_api_id  = aws_api_gateway_rest_api.api.id
#   name         = "profile"
#   description  = "a JSON schema"
#   content_type = "application/json"
#   schema = <<EOF
# {
#   "id": "string",
#   "name": "string",
#   "email": "string"
# }
# EOF
# }

resource "aws_api_gateway_method" "profile_get" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.profile_id.id
  http_method   = "GET"
  authorization = "NONE"
  request_parameters = {
    "method.request.path.id" = true
  }
}

resource "aws_api_gateway_method" "profile_create" {
  rest_api_id   = aws_api_gateway_rest_api.api.id
  resource_id   = aws_api_gateway_resource.profile.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "profile_get" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.profile_id.id
  http_method             = aws_api_gateway_method.profile_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.rest_lambda["profile-get"].invoke_arn
}

resource "aws_api_gateway_integration" "profile_create" {
  rest_api_id             = aws_api_gateway_rest_api.api.id
  resource_id             = aws_api_gateway_resource.profile.id
  http_method             = aws_api_gateway_method.profile_create.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.rest_lambda["profile-create"].invoke_arn
}

resource "aws_api_gateway_stage" "mci" {
  deployment_id = aws_api_gateway_deployment.mci.id
  rest_api_id   = aws_api_gateway_rest_api.api.id
  stage_name    = local.stage
  depends_on    = [aws_api_gateway_account.ApiGatewayAccountSetting]

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api.arn
    format          = "{ \"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\",\"routeKey\":\"$context.routeKey\", \"status\":\"$context.status\",\"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\" }"
  }
}

resource "aws_api_gateway_deployment" "mci" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.profile,
      aws_api_gateway_resource.profile_id,
      aws_api_gateway_method.profile_get,
      aws_api_gateway_method.profile_create,
      aws_api_gateway_integration.profile_get,
      aws_api_gateway_integration.profile_create,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

output "profile_url" {
  value = "${aws_api_gateway_stage.mci.invoke_url}${aws_api_gateway_resource.profile.path}"
}

# Permissions and logging

resource "aws_lambda_permission" "apigw_profile_create" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.rest_lambda["profile-create"].function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

resource "aws_lambda_permission" "apigw_profile_get" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = module.rest_lambda["profile-get"].function_name
  principal     = "apigateway.amazonaws.com"

  source_arn = "${aws_api_gateway_rest_api.api.execution_arn}/*/*/*"
}

# Create a Log Group for API Gateway to push logs to
resource "aws_cloudwatch_log_group" "api" {
  name_prefix = "/aws/APIGW/mci-api"
}

# Create a Log Policy to allow Cloudwatch to Create log streams and put logs
resource "aws_cloudwatch_log_resource_policy" "api" {
  policy_name     = "Terraform-CloudWatchLogPolicy-${data.aws_caller_identity.current.account_id}"
  policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Id": "CWLogsPolicy",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": [ 
          "apigateway.amazonaws.com",
          "delivery.logs.amazonaws.com"
          ]
      },
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
        ],
      "Resource": "${aws_cloudwatch_log_group.api.arn}",
      "Condition": {
        "ArnEquals": {
          "aws:SourceArn": "${aws_api_gateway_rest_api.api.arn}"
        }
      }
    }
  ]
}
EOF  
}

resource "aws_api_gateway_account" "ApiGatewayAccountSetting" {
  cloudwatch_role_arn = aws_iam_role.APIGatewayCloudWatchRole.arn
}

resource "aws_iam_role" "APIGatewayCloudWatchRole" {
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "",
      "Effect": "Allow",
      "Principal": {
        "Service": "apigateway.amazonaws.com"
      },
      "Action": "sts:AssumeRole"
    }
  ]
}
EOF
}

resource "aws_iam_role_policy" "APIGatewayCloudWatchPolicy" {
  role = aws_iam_role.APIGatewayCloudWatchRole.id

  policy = <<EOF
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": [
                "logs:CreateLogGroup",
                "logs:CreateLogStream",
                "logs:DescribeLogGroups",
                "logs:DescribeLogStreams",
                "logs:PutLogEvents",
                "logs:GetLogEvents",
                "logs:FilterLogEvents"
            ],
            "Resource": "*"
        }
    ]
}
EOF
}

# Configure API Gateway to push all logs to CloudWatch Logs
resource "aws_api_gateway_method_settings" "MyApiGatewaySetting" {
  rest_api_id = aws_api_gateway_rest_api.api.id
  stage_name  = aws_api_gateway_stage.mci.stage_name
  method_path = "*/*"

  settings {
    # Enable CloudWatch logging and metrics
    metrics_enabled = true
    logging_level   = "INFO"
  }
}