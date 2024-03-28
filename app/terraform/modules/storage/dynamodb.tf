# ------------------------------------------------------------
# Voter Table
# ------------------------------------------------------------
resource "aws_dynamodb_table" "voter_table" {
  name          = "voter-table"
  billing_mode  = var.dynamodb_billing_mode
  hash_key      = var.voter_id_attribute
  range_key     = var.ballot_id_attribute 

  attribute {
    name = var.voter_id_attribute
    type = var.dynamodb_attribute_type_string
  }

  attribute {
    name = var.ballot_id_attribute
    type = var.dynamodb_attribute_type_string
  }
}