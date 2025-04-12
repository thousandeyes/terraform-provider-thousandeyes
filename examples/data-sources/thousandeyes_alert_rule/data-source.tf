data "thousandeyes_alert_rule" "example_alert_rule" {
  rule_name = "Alert rule name"
}


resource "thousandeyes_http_server" "www_thousandeyes_http_test" {
  test_name      = "Example HTTP test set from Terraform provider"
  interval       = 120
  url            = "https://www.thousandeyes.com"
  alerts_enabled = true
  alert_rules    = [data.thousandeyes_alert_rule.example_alert_rule.id]
}
