data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT DNSSEC Alert Rule"
  alert_type                = "dnssec"
  expression                = "((probDetail != \"\"))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_dnssec" "test" {
  test_name      = "User Acceptance Test - DNSSEC"
  interval       = 120
  alerts_enabled = true
  domain         = "thousandeyes.com A"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, "921613"]
}
