# Input variable definitions

variable "entity" {
  description = "Name of the REST entity/resource."
  type        = string
}

variable "operation" {
  description = "Name of the REST operation/method."
  type        = string
}

variable "environment" {
  description = "Map of environment variables"
  type = map
  nullable = true
  default = {}
}

variable "iam_role_arn" {
  description = "ARN of the basic IAM role for Lambdas"
  type = string
}