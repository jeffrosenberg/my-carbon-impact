resource "aws_dynamodb_table" "profiles" {
  name           = "profiles"
  billing_mode   = "PAY_PER_REQUEST"
  hash_key       = "id"

  attribute {
    name = "id"
    type = "B"
  }

  attribute {
    name = "email"
    type = "S"
  }

  # Skipping declaration of TTL due to this issue:
  # https://github.com/hashicorp/terraform-provider-aws/issues/10304
  # ttl {
  #   attribute_name = "time_to_exist"
  #   enabled        = false
  # }

  global_secondary_index {
    name               = "profiles_by_email"
    hash_key           = "email"
    projection_type    = "INCLUDE"
    non_key_attributes = ["name"] # TODO: come back to this
  }
}