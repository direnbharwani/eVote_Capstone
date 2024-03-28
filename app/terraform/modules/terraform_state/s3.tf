resource "aws_s3_bucket" "terraform_state_bucket" {
  bucket        = var.terraform_state_bucket_name
  force_destroy = true
}

resource "aws_s3_bucket_ownership_controls" "terraform_state_ownership_controls" {
  bucket = aws_s3_bucket.terraform_state_bucket.id
  rule {
    object_ownership = "BucketOwnerPreferred"
  }
}

resource "aws_s3_bucket_acl" "terraform_state_acl" {
  depends_on  = [aws_s3_bucket_ownership_controls.terraform_state_ownership_controls]

  bucket      = aws_s3_bucket.terraform_state_bucket.id
  acl         = "private"
}

resource "aws_s3_bucket_versioning" "terraform_state_versioning"{
  bucket = aws_s3_bucket.terraform_state_bucket.id
  versioning_configuration {
    status = "Enabled"
  }
}

resource "aws_s3_bucket_server_side_encryption_configuration" "terraform_state_encryption" {
  bucket = aws_s3_bucket.terraform_state_bucket.id
  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm = "AES256"
    }
  }
}