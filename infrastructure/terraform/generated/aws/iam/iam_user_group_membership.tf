resource "aws_iam_user_group_membership" "tfer--anto_dev-002F-admin_user" {
  groups = ["admin_user"]
  user   = "anto_dev"
}

resource "aws_iam_user_group_membership" "tfer--anto_user-002F-dev_ops_example" {
  groups = ["dev_ops_example"]
  user   = "anto_user"
}

resource "aws_iam_user_group_membership" "tfer--github-actions-deploy-002F-github_deploy_bot" {
  groups = ["github_deploy_bot"]
  user   = "github-actions-deploy"
}
