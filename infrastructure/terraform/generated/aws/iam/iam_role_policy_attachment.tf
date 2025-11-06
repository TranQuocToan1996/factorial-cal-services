resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForAmazonEKS_AmazonEKSServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AmazonEKSServiceRolePolicy"
  role       = "AWSServiceRoleForAmazonEKS"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForAmazonMQ_AmazonMQServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AmazonMQServiceRolePolicy"
  role       = "AWSServiceRoleForAmazonMQ"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForApplicationAutoScaling_ECSService_AWSApplicationAutoscalingECSServicePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSApplicationAutoscalingECSServicePolicy"
  role       = "AWSServiceRoleForApplicationAutoScaling_ECSService"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForCostOptimizationHub_CostOptimizationHubServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/CostOptimizationHubServiceRolePolicy"
  role       = "AWSServiceRoleForCostOptimizationHub"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForECSCompute_AmazonECSComputeServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AmazonECSComputeServiceRolePolicy"
  role       = "AWSServiceRoleForECSCompute"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForECS_AmazonECSServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AmazonECSServiceRolePolicy"
  role       = "AWSServiceRoleForECS"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForElasticLoadBalancing_AWSElasticLoadBalancingServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSElasticLoadBalancingServiceRolePolicy"
  role       = "AWSServiceRoleForElasticLoadBalancing"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForGlobalAccelerator_AWSGlobalAcceleratorSLRPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSGlobalAcceleratorSLRPolicy"
  role       = "AWSServiceRoleForGlobalAccelerator"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForRDS_AmazonRDSServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AmazonRDSServiceRolePolicy"
  role       = "AWSServiceRoleForRDS"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForResourceExplorer_AWSResourceExplorerServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSResourceExplorerServiceRolePolicy"
  role       = "AWSServiceRoleForResourceExplorer"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForSupport_AWSSupportServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSSupportServiceRolePolicy"
  role       = "AWSServiceRoleForSupport"
}

resource "aws_iam_role_policy_attachment" "tfer--AWSServiceRoleForTrustedAdvisor_AWSTrustedAdvisorServiceRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/aws-service-role/AWSTrustedAdvisorServiceRolePolicy"
  role       = "AWSServiceRoleForTrustedAdvisor"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_cluster_default_role_AmazonEKSBlockStoragePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSBlockStoragePolicy"
  role       = "EKS_cluster_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_cluster_default_role_AmazonEKSClusterPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = "EKS_cluster_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_cluster_default_role_AmazonEKSComputePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSComputePolicy"
  role       = "EKS_cluster_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_cluster_default_role_AmazonEKSLoadBalancingPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSLoadBalancingPolicy"
  role       = "EKS_cluster_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_cluster_default_role_AmazonEKSNetworkingPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSNetworkingPolicy"
  role       = "EKS_cluster_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_node_default_role_AmazonEC2ContainerRegistryPullOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly"
  role       = "EKS_node_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--EKS_node_default_role_AmazonEKSWorkerNodeMinimalPolicy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodeMinimalPolicy"
  role       = "EKS_node_default_role"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_AmazonEC2ContainerRegistryPowerUser" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_AmazonEC2ContainerRegistryPullOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPullOnly"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_AmazonEC2ContainerRegistryReadOnly" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_AmazonECSTaskExecutionRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_AmazonECS_FullAccess" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonECS_FullAccess"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--GitHubActionsECSRole_CloudWatchLogsFullAccess" {
  policy_arn = "arn:aws:iam::aws:policy/CloudWatchLogsFullAccess"
  role       = "GitHubActionsECSRole"
}

resource "aws_iam_role_policy_attachment" "tfer--ecs-task-execution-role_AmazonECSTaskExecutionRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = "ecs-task-execution-role"
}

resource "aws_iam_role_policy_attachment" "tfer--ecs-task-role_AmazonECSTaskExecutionRolePolicy" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonECSTaskExecutionRolePolicy"
  role       = "ecs-task-role"
}

resource "aws_iam_role_policy_attachment" "tfer--ecsInfrastructureRole_AmazonECSInfrastructureRolePolicyForManagedInstances" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonECSInfrastructureRolePolicyForManagedInstances"
  role       = "ecsInfrastructureRole"
}

resource "aws_iam_role_policy_attachment" "tfer--ecsInstanceRole_AmazonECSInstanceRolePolicyForManagedInstances" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonECSInstanceRolePolicyForManagedInstances"
  role       = "ecsInstanceRole"
}

resource "aws_iam_role_policy_attachment" "tfer--rds-monitoring-role_AmazonRDSEnhancedMonitoringRole" {
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonRDSEnhancedMonitoringRole"
  role       = "rds-monitoring-role"
}
