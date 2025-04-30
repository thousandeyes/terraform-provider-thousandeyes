resource "thousandeyes_user" "example_user" {
  name                   = "Example User"
  email                  = "example@test.com"
  login_account_group_id = "1234"
  account_group_roles {
    account_group_id = "1234"
    role_ids         = ["1"]
  }
  all_account_group_role_ids = ["1"]
}
