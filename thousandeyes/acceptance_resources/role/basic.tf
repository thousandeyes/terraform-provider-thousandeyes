data "thousandeyes_permission" "reports_read" {
  permission_name = "REPORTS_READ"
}

resource "thousandeyes_role" "test" {
  name        = "User Acceptance Test - Role"
  permissions = [data.thousandeyes_permission.reports_read.permission_id]
}
