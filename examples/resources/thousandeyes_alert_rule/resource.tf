resource "thousandeyes_alert_rule" "example_alert_rule" {
  rule_name                 = "Example Alert Rule set from Terraform provider"
  alert_type                = "http-server"
  expression                = "((errorType != \"None\"))" # Error Type is ANY
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}
