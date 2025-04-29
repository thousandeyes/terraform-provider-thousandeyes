data "thousandeyes_account_group" "sre" {
  name = "SRE"
}

resource "thousandeyes_http_server" "thousandeyes_http_test" {
  test_name            = "Example HTTP test with account group sharing set from Terraform provider"
  interval             = 120
  alerts_enabled       = true
  url                  = "https://www.thousandeyes.com"
  shared_with_accounts = [data.thousandeyes_account_group.sre.aid]
  agents               = ["3"] # Singapore
}
