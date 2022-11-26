resource "aws_api_gateway_rest_api" "web" {
  name = "My Carbon Impact frontend"
}

resource "aws_api_gateway_resource" "profile_page" {
  path_part   = "profile/{id+}"
  parent_id   = aws_api_gateway_rest_api.web.root_resource_id
  rest_api_id = aws_api_gateway_rest_api.web.id
}

resource "aws_api_gateway_method" "profile_page_get" {
  rest_api_id   = aws_api_gateway_rest_api.web.id
  resource_id   = aws_api_gateway_resource.profile_page.id
  http_method   = "GET"
  authorization = "NONE"
  request_parameters = {
    "method.request.path.id" = true
  }
}

resource "aws_api_gateway_integration" "profile_page_get" {
  rest_api_id             = aws_api_gateway_rest_api.web.id
  resource_id             = aws_api_gateway_resource.profile_page.id
  http_method             = aws_api_gateway_method.profile_page_get.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = module.web_lambda["profile-get"].invoke_arn
}

resource "aws_api_gateway_stage" "web" {
  deployment_id = aws_api_gateway_deployment.web.id
  rest_api_id   = aws_api_gateway_rest_api.web.id
  stage_name    = var.api_stage
  depends_on    = [aws_api_gateway_account.ApiGatewayAccountSetting]

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api.arn
    format          = "{ \"requestId\":\"$context.requestId\", \"ip\": \"$context.identity.sourceIp\", \"requestTime\":\"$context.requestTime\", \"httpMethod\":\"$context.httpMethod\",\"routeKey\":\"$context.routeKey\", \"status\":\"$context.status\",\"protocol\":\"$context.protocol\", \"responseLength\":\"$context.responseLength\" }"
  }
}

resource "aws_api_gateway_deployment" "web" {
  rest_api_id = aws_api_gateway_rest_api.web.id
  
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.profile_page,
      aws_api_gateway_method.profile_page_get,
      aws_api_gateway_integration.profile_page_get,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_cloudfront_distribution" "profile_get" {
  origin {
    origin_id = "profile-get"
    domain_name = "my-carbon-impact.jeff-rosenberg.com"
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "My Carbon Impact profile display"
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "profile-get"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    lambda_function_association {
      event_type   = "viewer-request"
      lambda_arn   = module.web_lambda["profile-get"].invoke_arn
      include_body = true
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none" 
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

resource "aws_cloudfront_distribution" "profile_create" {
  origin {
    origin_id = "profile-create"
    domain_name = "my-carbon-impact.jeff-rosenberg.com"
  }

  enabled             = true
  is_ipv6_enabled     = true
  comment             = "My Carbon Impact profile display"
  default_root_object = "index.html"

  default_cache_behavior {
    allowed_methods  = ["GET", "HEAD", "OPTIONS"]
    cached_methods   = ["GET", "HEAD"]
    target_origin_id = "profile-create"

    forwarded_values {
      query_string = false

      cookies {
        forward = "none"
      }
    }

    lambda_function_association {
      event_type   = "viewer-request"
      lambda_arn   = module.web_lambda["profile-create"].invoke_arn
      include_body = true
    }

    viewer_protocol_policy = "allow-all"
    min_ttl                = 0
    default_ttl            = 3600
    max_ttl                = 86400
  }

  restrictions {
    geo_restriction {
      restriction_type = "none" 
    }
  }

  viewer_certificate {
    cloudfront_default_certificate = true
  }
}

output "web_url" {
  value = "${aws_api_gateway_stage.web.invoke_url}"
}