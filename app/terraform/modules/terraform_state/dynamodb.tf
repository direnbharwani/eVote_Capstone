# ------------------------------------------------------------
# Terraform State Lock
# ------------------------------------------------------------
resource "aws_dynamodb_table" "terraform_locks_table" {
  name          = var.terraform_state_lock_table_name
  billing_mode  = "PAY_PER_REQUEST"
  hash_key      = "LockID"

  attribute {
    name = "LockID"
    type = "S" # string
  }
}