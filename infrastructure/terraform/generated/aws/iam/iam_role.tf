resource "aws_iam_role" "tfer--AWSServiceRoleForAmazonEKS" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows EKS to manage clusters on your behalf."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AmazonEKSServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForAmazonEKS"
  path                 = "/aws-service-role/eks.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForAmazonMQ" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "mq.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows Amazon MQ to call AWS services on your behalf"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AmazonMQServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForAmazonMQ"
  path                 = "/aws-service-role/mq.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForApplicationAutoScaling_ECSService" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.application-autoscaling.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSApplicationAutoscalingECSServicePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForApplicationAutoScaling_ECSService"
  path                 = "/aws-service-role/ecs.application-autoscaling.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForCostOptimizationHub" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "cost-optimization-hub.bcm.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows Cost Optimization Hub to retrieve organization information and collect optimization-related data and metadata."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/CostOptimizationHubServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForCostOptimizationHub"
  path                 = "/aws-service-role/cost-optimization-hub.bcm.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForECS" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Policy to enable Amazon ECS to manage your EC2 instances and related resources."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AmazonECSServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForECS"
  path                 = "/aws-service-role/ecs.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForECSCompute" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-compute.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Policy to enable Amazon ECS Compute to manage your EC2 instances and related resources as part of ECS managed instance."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AmazonECSComputeServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForECSCompute"
  path                 = "/aws-service-role/ecs-compute.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForElasticLoadBalancing" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "elasticloadbalancing.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows ELB to call AWS services on your behalf."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSElasticLoadBalancingServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForElasticLoadBalancing"
  path                 = "/aws-service-role/elasticloadbalancing.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForGlobalAccelerator" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "globalaccelerator.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows Global Accelerator to call AWS services on customer's behalf"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSGlobalAcceleratorSLRPolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForGlobalAccelerator"
  path                 = "/aws-service-role/globalaccelerator.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForRDS" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "rds.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows Amazon RDS to manage AWS resources on your behalf"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AmazonRDSServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForRDS"
  path                 = "/aws-service-role/rds.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForResourceExplorer" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "resource-explorer-2.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSResourceExplorerServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForResourceExplorer"
  path                 = "/aws-service-role/resource-explorer-2.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForSupport" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "support.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Enables resource access for AWS to provide billing, administrative and support services"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForSupport"
  path                 = "/aws-service-role/support.amazonaws.com/"
}

resource "aws_iam_role" "tfer--AWSServiceRoleForTrustedAdvisor" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "trustedadvisor.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Access for the AWS Trusted Advisor Service to help reduce cost, increase performance, and improve security of your AWS environment."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/aws-service-role/AWSTrustedAdvisorServiceRolePolicy"]
  max_session_duration = "3600"
  name                 = "AWSServiceRoleForTrustedAdvisor"
  path                 = "/aws-service-role/trustedadvisor.amazonaws.com/"
}

resource "aws_iam_role" "tfer--EKS_cluster_default_role" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": [
        "sts:AssumeRole",
        "sts:TagSession"
      ],
      "Effect": "Allow",
      "Principal": {
        "Service": "eks.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows access to other AWS service resources that are required to operate Auto Mode clusters managed by EKS."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy", "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy", "arn:aws:iam::aws:policy/AmazonEKSComputePolicy", "arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy", "arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy"]
  max_session_duration = "3600"
  name                 = "EKS_cluster_default_role"
  path                 = "/"
}

resource "aws_iam_role" "tfer--EKS_node_default_role" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows EKS nodes to connect to EKS Auto Mode clusters and to pull container images from ECR."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly", "arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy"]
  max_session_duration = "3600"
  name                 = "EKS_node_default_role"
  path                 = "/"
}

resource "aws_iam_role" "tfer--GitHubActionsECSRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRoleWithWebIdentity",
      "Condition": {
        "StringLike": {
          "token.actions.githubusercontent.com:aud": "sts.amazonaws.com",
          "token.actions.githubusercontent.com:sub": [
            "repo:TranQuocToan1996/*:ref:refs/heads/main",
            "repo:TranQuocToan1996/*:ref:refs/heads/master",
            "repo:TranQuocToan1996/*:ref:refs/heads/dev",
            "repo:TranQuocToan1996/*:ref:refs/heads/staging"
          ]
        }
      },
      "Effect": "Allow",
      "Principal": {
        "Federated": "arn:aws:iam::218435950768:oidc-provider/token.actions.githubusercontent.com"
      }
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "GitHubActionsECSRole"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser", "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly", "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly", "arn:aws:iam::aws:policy/AmazonECS_FullAccess", "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess", "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
  max_session_duration = "3600"
  name                 = "GitHubActionsECSRole"
  path                 = "/"
}

resource "aws_iam_role" "tfer--ecs-task-execution-role" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Sid": "AllowAccessToECSForTaskExecutionRole"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description = "Allows access to other AWS service resources that are required to run Amazon ECS tasks."

  inline_policy {
    name   = "ecs-full-ecr-access"
    policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"ecr:GetAuthorizationToken\",\"ecr:BatchGetImage\",\"ecr:GetDownloadUrlForLayer\"],\"Effect\":\"Allow\",\"Resource\":\"*\"},{\"Action\":[\"logs:CreateLogStream\",\"logs:PutLogEvents\"],\"Effect\":\"Allow\",\"Resource\":\"*\"},{\"Action\":[\"secretsmanager:GetSecretValue\"],\"Effect\":\"Allow\",\"Resource\":[\"arn:aws:secretsmanager:us-east-1:218435950768:secret:dev/simple-order-service/api-FCJNv4*\"]},{\"Action\":[\"kms:Decrypt\"],\"Condition\":{\"StringEquals\":{\"kms:ViaService\":[\"secretsmanager.us-east-1.amazonaws.com\",\"ecr.us-east-1.amazonaws.com\"]}},\"Effect\":\"Allow\",\"Resource\":\"*\"}]}"
  }

  inline_policy {
    name   = "ecs-get-secret-simple-order"
    policy = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Action\":[\"secretsmanager:GetSecretValue\"],\"Effect\":\"Allow\",\"Resource\":[\"arn:aws:secretsmanager:us-east-1:218435950768:secret:dev/factorial-service-PuUw9d*\"],\"Sid\":\"Statement1\"}]}"
  }

  managed_policy_arns  = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
  max_session_duration = "3600"
  name                 = "ecs-task-execution-role"
  path                 = "/"
}

resource "aws_iam_role" "tfer--ecs-task-role" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs-tasks.amazonaws.com"
      },
      "Sid": ""
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "ecs-task-role"
  managed_policy_arns  = ["arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"]
  max_session_duration = "3600"
  name                 = "ecs-task-role"
  path                 = "/"
}

resource "aws_iam_role" "tfer--ecsInfrastructureRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ecs.amazonaws.com"
      },
      "Sid": "AllowAccessToECSForInfrastructureManagement"
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows ECS to create and manage AWS EC2 resources on your behalf."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/AmazonECSInfrastructureRolePolicyForManagedInstances"]
  max_session_duration = "3600"
  name                 = "ecsInfrastructureRole"
  path                 = "/"
}

resource "aws_iam_role" "tfer--ecsInstanceRole" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "ec2.amazonaws.com"
      },
      "Sid": ""
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  description          = "Allows EC2 instances to communicate with ECS on your behalf."
  managed_policy_arns  = ["arn:aws:iam::aws:policy/AmazonECSInstanceRolePolicyForManagedInstances"]
  max_session_duration = "3600"
  name                 = "ecsInstanceRole"
  path                 = "/"
}

resource "aws_iam_role" "tfer--rds-monitoring-role" {
  assume_role_policy = <<POLICY
{
  "Statement": [
    {
      "Action": "sts:AssumeRole",
      "Effect": "Allow",
      "Principal": {
        "Service": "monitoring.rds.amazonaws.com"
      },
      "Sid": ""
    }
  ],
  "Version": "2012-10-17"
}
POLICY

  managed_policy_arns  = ["arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"]
  max_session_duration = "3600"
  name                 = "rds-monitoring-role"
  path                 = "/"
}
