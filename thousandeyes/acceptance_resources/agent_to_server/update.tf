data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Agent To Server Alert Rule (Updated)"
  alert_type                = "end-to-end-server"
  expression                = "((loss >= 10%) || (probDetail != \"\"))"
  minimum_sources           = 1
  rounds_violating_required = 3
  rounds_violating_out_of   = 3
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "User Acceptance Test - Agent To Server (Updated)"
  interval       = 300
  alerts_enabled = true
  server         = "api.stg.thousandeyes.com"
  protocol       = "tcp"
  port           = 443
  probe_mode     = "sack"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, "921610"]
}
