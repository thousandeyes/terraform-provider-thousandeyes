data "thousandeyes_permission" "example_permission" {
  permission_name = "EMAILS_UPDATE"
}

resource "thousandeyes_role" "example_role" {
  name        = "Example Role"
  permissions = [data.thousandeyes_permission.example_permission.permission_id]
}
