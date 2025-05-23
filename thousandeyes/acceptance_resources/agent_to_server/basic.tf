data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default Network Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Agent To Server Alert Rule"
  alert_type                = "end-to-end-server"
  expression                = "((loss >= 10%) || (probDetail != \"\"))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "User Acceptance Test - Agent To Server"
  interval       = 120
  alerts_enabled = true
  server         = "api.stg.thousandeyes.com"
  protocol       = "tcp"
  port           = 81   
  probe_mode     = "sack"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
