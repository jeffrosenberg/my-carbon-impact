# Input variable definitions

variable "entity" {
  description = "Name of the REST entity/resource."
  type        = string
}

variable "operation" {
  description = "Name of the REST operation/method."
  type        = string
}

variable "log_levels" {
  description = "Map of log levels to Zerolog ints"
  type = map
}

variable "iam_role_arn" {
  description = "ARN of the basic IAM role for Lambdas"
  type = string
}

variable "log_level" {
  description = "Logging level"
  type = string
  default = "TRACE"
}