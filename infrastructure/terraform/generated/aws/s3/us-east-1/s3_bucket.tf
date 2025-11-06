resource "aws_s3_bucket" "tfer--aws-cloudtrail-logs-218435950768-6da8860e" {
  bucket        = "aws-cloudtrail-logs-218435950768-6da8860e"
  force_destroy = "false"

  grant {
    id          = "c08b218b0b43459263e7f0f15f34bbe134322dc1dc88286dcad1fa1f3940d709"
    permissions = ["FULL_CONTROL"]
    type        = "CanonicalUser"
  }

  object_lock_enabled = "false"

  policy = <<POLICY
{
  "Statement": [
    {
      "Action": "s3:GetBucketAcl",
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudtrail:us-east-1:218435950768:trail/management-events"
        }
      },
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudtrail.amazonaws.com"
      },
      "Resource": "arn:aws:s3:::aws-cloudtrail-logs-218435950768-6da8860e",
      "Sid": "AWSCloudTrailAclCheck20150319-8841942a-e2fa-42e8-b516-54f3604fa8f9"
    },
    {
      "Action": "s3:PutObject",
      "Condition": {
        "StringEquals": {
          "AWS:SourceArn": "arn:aws:cloudtrail:us-east-1:218435950768:trail/management-events",
          "s3:x-amz-acl": "bucket-owner-full-control"
        }
      },
      "Effect": "Allow",
      "Principal": {
        "Service": "cloudtrail.amazonaws.com"
      },
      "Resource": "arn:aws:s3:::aws-cloudtrail-logs-218435950768-6da8860e/AWSLogs/218435950768/*",
      "Sid": "AWSCloudTrailWrite20150319-57aee438-f822-470a-ab14-945b300326db"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  region        = "us-east-1"
  request_payer = "BucketOwner"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }

      bucket_key_enabled = "false"
    }
  }

  versioning {
    enabled    = "false"
    mfa_delete = "false"
  }
}

resource "aws_s3_bucket" "tfer--factorial-calculator-service" {
  bucket        = "factorial-calculator-service"
  force_destroy = "false"

  grant {
    id          = "c08b218b0b43459263e7f0f15f34bbe134322dc1dc88286dcad1fa1f3940d709"
    permissions = ["FULL_CONTROL"]
    type        = "CanonicalUser"
  }

  object_lock_enabled = "false"
  region              = "us-east-1"
  request_payer       = "BucketOwner"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }

      bucket_key_enabled = "true"
    }
  }

  versioning {
    enabled    = "true"
    mfa_delete = "false"
  }
}

resource "aws_s3_bucket" "tfer--simple-order-service" {
  bucket        = "simple-order-service"
  force_destroy = "false"

  grant {
    id          = "c08b218b0b43459263e7f0f15f34bbe134322dc1dc88286dcad1fa1f3940d709"
    permissions = ["FULL_CONTROL"]
    type        = "CanonicalUser"
  }

  object_lock_enabled = "false"
  region              = "us-east-1"
  request_payer       = "BucketOwner"

  server_side_encryption_configuration {
    rule {
      apply_server_side_encryption_by_default {
        sse_algorithm = "AES256"
      }

      bucket_key_enabled = "true"
    }
  }

  versioning {
    enabled    = "true"
    mfa_delete = "false"
  }
}
