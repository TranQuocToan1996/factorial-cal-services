resource "aws_iam_instance_profile" "tfer--ecsInstanceRole" {
  name = "ecsInstanceRole"
  path = "/"
  role = "ecsInstanceRole"
}
