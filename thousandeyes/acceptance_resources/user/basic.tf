data "thousandeyes_permission" "emails_update" {
  permission_name = "EMAILS_UPDATE"
}

data "thousandeyes_account_group" "acc_group" {
  name = "API Team TF Providers"
}

resource "thousandeyes_role" "test_role" {
  name        = "User Acceptance Test - Role"
  permissions = [data.thousandeyes_permission.emails_update.id]
}

resource "thousandeyes_user" "test" {
  name                   = "User Acceptance Test - User"
  email                  = "example@test.com"
  login_account_group_id = data.thousandeyes_account_group.acc_group.id
  account_group_roles {
    account_group_id = data.thousandeyes_account_group.acc_group.id
    role_ids         = [thousandeyes_role.test_role.id]
  }
  all_account_group_role_ids = [thousandeyes_role.test_role.id]
}
