data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default HTTP Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "alert-rule-http-test" {
  rule_name                 = "Custom UAT HTTP Alert Rule (Updated)"
  alert_type                = "http-server"
  expression                = "(((responseCode >= 400) || (responseCode == 0)))"
  rounds_violating_out_of   = 3
  rounds_violating_required = 3
  minimum_sources           = 1
}

resource "thousandeyes_http_server" "test" {
  test_name      = "A User Acceptance Test - HTTP (Updated)"
  interval       = 300
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  alert_rules    = [data.thousandeyes_alert_rule.def_alert_rule.id, thousandeyes_alert_rule.alert-rule-http-test.id]
}
