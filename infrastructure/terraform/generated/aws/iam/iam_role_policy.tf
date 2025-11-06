resource "aws_iam_role_policy" "tfer--ecs-task-execution-role_ecs-full-ecr-access" {
  name = "ecs-full-ecr-access"

  policy = <<POLICY
{
  "Statement": [
    {
      "Action": [
        "ecr:GetAuthorizationToken",
        "ecr:BatchGetImage",
        "ecr:GetDownloadUrlForLayer"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "logs:CreateLogStream",
        "logs:PutLogEvents"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:secretsmanager:us-east-1:218435950768:secret:dev/simple-order-service/api-FCJNv4*"
      ]
    },
    {
      "Action": [
        "kms:Decrypt"
      ],
      "Condition": {
        "StringEquals": {
          "kms:ViaService": [
            "secretsmanager.us-east-1.amazonaws.com",
            "ecr.us-east-1.amazonaws.com"
          ]
        }
      },
      "Effect": "Allow",
      "Resource": "*"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  role = "ecs-task-execution-role"
}

resource "aws_iam_role_policy" "tfer--ecs-task-execution-role_ecs-get-secret-simple-order" {
  name = "ecs-get-secret-simple-order"

  policy = <<POLICY
{
  "Statement": [
    {
      "Action": [
        "secretsmanager:GetSecretValue"
      ],
      "Effect": "Allow",
      "Resource": [
        "arn:aws:secretsmanager:us-east-1:218435950768:secret:dev/factorial-service-PuUw9d*"
      ],
      "Sid": "Statement1"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  role = "ecs-task-execution-role"
}
