data "thousandeyes_permission" "reports_read" {
  permission_name = "REPORTS_READ"
}

data "thousandeyes_permission" "emails_update" {
  permission_name = "EMAILS_UPDATE"
}

resource "thousandeyes_role" "test" {
  name        = "User Acceptance Test - Role (Updated)"
  permissions           = [data.thousandeyes_permission.reports_read.permission_id, data.thousandeyes_permission.emails_update.permission_id]
}
