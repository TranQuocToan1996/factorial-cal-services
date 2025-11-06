resource "aws_iam_policy" "tfer--S3_READ_ONLY" {
  name = "S3_READ_ONLY"
  path = "/"

  policy = <<POLICY
{
  "Statement": [
    {
      "Action": [
        "s3:GetObjectVersionTagging",
        "s3:GetStorageLensConfigurationTagging",
        "s3:GetObjectAcl",
        "s3:GetBucketObjectLockConfiguration",
        "s3:GetIntelligentTieringConfiguration",
        "s3:GetStorageLensGroup",
        "s3:GetAccessGrantsInstanceForPrefix",
        "s3:GetObjectVersionAcl",
        "s3:GetBucketPolicyStatus",
        "s3:GetAccessGrantsLocation",
        "s3:GetObjectRetention",
        "s3:GetBucketWebsite",
        "s3:GetJobTagging",
        "s3:GetMultiRegionAccessPoint",
        "s3:GetObjectAttributes",
        "s3:GetAccessGrantsInstanceResourcePolicy",
        "s3:GetObjectLegalHold",
        "s3:GetBucketNotification",
        "s3:DescribeMultiRegionAccessPointOperation",
        "s3:GetReplicationConfiguration",
        "s3:GetObject",
        "s3:GetBucketMetadataTableConfiguration",
        "s3:DescribeJob",
        "s3:GetAnalyticsConfiguration",
        "s3:GetObjectVersionForReplication",
        "s3:GetAccessPointForObjectLambda",
        "s3:GetStorageLensDashboard",
        "s3:GetLifecycleConfiguration",
        "s3:GetAccessPoint",
        "s3:GetInventoryConfiguration",
        "s3:GetBucketTagging",
        "s3:GetAccessPointPolicyForObjectLambda",
        "s3:GetBucketLogging",
        "s3:GetAccessGrant",
        "s3:GetAccelerateConfiguration",
        "s3:GetObjectVersionAttributes",
        "s3:GetBucketPolicy",
        "s3:GetEncryptionConfiguration",
        "s3:GetObjectVersionTorrent",
        "s3:GetBucketRequestPayment",
        "s3:GetAccessPointPolicyStatus",
        "s3:GetAccessGrantsInstance",
        "s3:GetObjectTagging",
        "s3:GetMetricsConfiguration",
        "s3:GetBucketOwnershipControls",
        "s3:GetBucketPublicAccessBlock",
        "s3:GetMultiRegionAccessPointPolicyStatus",
        "s3:GetMultiRegionAccessPointPolicy",
        "s3:GetAccessPointPolicyStatusForObjectLambda",
        "s3:GetDataAccess",
        "s3:GetBucketVersioning",
        "s3:GetBucketAcl",
        "s3:GetAccessPointConfigurationForObjectLambda",
        "s3:GetObjectTorrent",
        "s3:GetMultiRegionAccessPointRoutes",
        "s3:GetStorageLensConfiguration",
        "s3:GetAccountPublicAccessBlock",
        "s3:GetBucketCORS",
        "s3:GetBucketLocation",
        "s3:GetAccessPointPolicy",
        "s3:GetObjectVersion"
      ],
      "Effect": "Allow",
      "Resource": "*",
      "Sid": "VisualEditor0"
    }
  ],
  "Version": "2012-10-17"
}
POLICY
}

resource "aws_iam_policy" "tfer--User_Billing" {
  description = "User_Billing"
  name        = "User_Billing"
  path        = "/"

  policy = <<POLICY
{
  "Statement": [
    {
      "Action": [
        "aws-portal:ViewBilling",
        "aws-portal:ViewUsage",
        "aws-portal:ViewAccount",
        "aws-portal:ModifyAccount",
        "budgets:ViewBudget",
        "budgets:ModifyBudget",
        "ce:Get*",
        "ce:List*",
        "ce:Create*",
        "ce:Update*",
        "ce:Delete*",
        "pricing:GetProducts"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ],
  "Version": "2012-10-17"
}
POLICY
}
