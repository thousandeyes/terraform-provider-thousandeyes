resource "thousandeyes_alert_rule" "example_alert_rule" {
  rule_name                 = "Example Alert Rule set from Terraform provider"
  alert_type                = "http-server"
  expression                = "((errorType != \"None\"))" # Error Type is ANY
  minimum_sources           = 2
  rounds_violating_required = 4
  rounds_violating_out_of   = 4
}
