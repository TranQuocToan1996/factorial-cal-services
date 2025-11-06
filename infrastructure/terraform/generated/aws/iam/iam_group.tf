resource "aws_iam_group" "tfer--admin_user" {
  name = "admin_user"
  path = "/"
}

resource "aws_iam_group" "tfer--dev_example" {
  name = "dev_example"
  path = "/"
}

resource "aws_iam_group" "tfer--dev_ops_example" {
  name = "dev_ops_example"
  path = "/"
}

resource "aws_iam_group" "tfer--github_deploy_bot" {
  name = "github_deploy_bot"
  path = "/"
}
