resource "aws_iam_user" "tfer--AIDATFW6MDCYBGFOSDHKV" {
  force_destroy = "false"
  name          = "anto_user"
  path          = "/"
}

resource "aws_iam_user" "tfer--AIDATFW6MDCYFRTVHFNRC" {
  force_destroy = "false"
  name          = "anto_dev"
  path          = "/"

  tags = {
    type = "admin"
  }

  tags_all = {
    type = "admin"
  }
}

resource "aws_iam_user" "tfer--AIDATFW6MDCYH4RBIJPQN" {
  force_destroy = "false"
  name          = "github-actions-deploy"
  path          = "/"
}
