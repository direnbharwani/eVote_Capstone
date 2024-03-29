# ------------------------------------------------------------
# Voter Credentials
# ------------------------------------------------------------
resource "aws_dynamodb_table" "voter_credentials_table" {
  name         = "voter-credentials"
  billing_mode = var.dynamodb_billing_mode
  hash_key     = "nric"         # partition key: unique to each user
  range_key    = "electionID"   # sort key: unique to each election. Allows query for all users for election.

  attribute {
    name = "nric"
    type = var.dynamodb_attribute_type_string
  }

  attribute {
    name = "electionID"
    type = var.dynamodb_attribute_type_number
  }
}