# Output variable definitions

output "arn" {
  description = "ARN of the Lambda function"
  value       = aws_lambda_function.lambda.arn
}

output "name" {
  description = "Name of the Lambda function"
  value       = aws_lambda_function.lambda.function_name
}
