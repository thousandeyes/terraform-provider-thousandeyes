data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}
data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default Agent to Agent Network Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "alert-rule-agent_to_agent" {
  rule_name                 = "Custom UAT Agent To Agent Alert Rule (Updated)"
  alert_type                = "end-to-end-agent"
  expression                = "((bothWaysLoss >= 10%) && (bothWaysProbDetail != \"\"))"
  direction                 = "bidirectional"
  rounds_violating_out_of   = 3
  rounds_violating_required = 3
  minimum_sources           = 1
}

resource "thousandeyes_agent_to_agent" "test" {
  test_name        = "User Acceptance Test - Aget To Agent (Updated)"
  interval         = 300
  alerts_enabled   = true
  bgp_measurements = true
  direction        = "bidirectional"
  protocol         = "tcp"
  target_agent_id  = data.thousandeyes_agent.arg_frankfurt.agent_id
  agents           = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  use_public_bgp   = true
  alert_rules      = [data.thousandeyes_alert_rule.def_alert_rule.id, thousandeyes_alert_rule.alert-rule-agent_to_agent.id]
}
