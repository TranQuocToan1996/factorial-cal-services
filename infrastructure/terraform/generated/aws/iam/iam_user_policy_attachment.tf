resource "aws_iam_user_policy_attachment" "tfer--anto_dev_IAMUserChangePassword" {
  policy_arn = "arn:aws:iam::aws:policy/IAMUserChangePassword"
  user       = "anto_dev"
}
