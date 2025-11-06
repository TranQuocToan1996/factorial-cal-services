resource "aws_iam_group_policy_attachment" "tfer--admin_user_AWSBillingReadOnlyAccess" {
  group      = "admin_user"
  policy_arn = "arn:aws:iam::aws:policy/AWSBillingReadOnlyAccess"
}

resource "aws_iam_group_policy_attachment" "tfer--admin_user_AdministratorAccess" {
  group      = "admin_user"
  policy_arn = "arn:aws:iam::aws:policy/AdministratorAccess"
}
