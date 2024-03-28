variable "terraform_state_bucket_name" {
  type        = string
  description = "The name of the S3 bucket that contains the terraform state."
  default     = "evote-poc-terraform-state"
}

variable "terraform_state_lock_table_name" {
  type        = string
  description = "The name of the DynamoDB table for locking the terraform state S3 bucket."
  default     = "terraform-locks"
}