variable "dynamodb_attribute_type_boolean" {
  type        = string
  description = "Binary type for a dynamodb attribute"
  default     = "B"
}

variable "dynamodb_attribute_type_number" {
  type        = string
  description = "Numeric type for a dynamodb attribute"
  default     = "N"
}

variable "dynamodb_attribute_type_string" {
  type        = string
  description = "String type for a dynamodb attribute"
  default     = "S"
}

variable "dynamodb_billing_mode" {
  type        = string
  description = "The billing mode for the dynamoDB storage tables"
  default     = "PAY_PER_REQUEST"
}

# ------------------------------------------------------------
# Voter Table
# ------------------------------------------------------------

variable "voter_id_attribute" {
  type        = string
  description = "An attribute that stores a voter id."
  default     = "voter_id"
}

variable "ballot_id_attribute" {
  type        = string
  description = "An attribute that stores a ballot id."
  default     = "ballot_id"
}

