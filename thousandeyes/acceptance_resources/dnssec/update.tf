data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default DNSSEC Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT DNSSEC Alert Rule (Updated)"
  alert_type                = "dnssec"
  expression                = "((probDetail != \"\"))"
  minimum_sources           = 1
  rounds_violating_required = 3
  rounds_violating_out_of   = 3
}

resource "thousandeyes_dnssec" "test" {
  test_name      = "User Acceptance Test - DNSSEC (Updated)"
  interval       = 300
  alerts_enabled = true
  domain         = "thousandeyes.com A"
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
