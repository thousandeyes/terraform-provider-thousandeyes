data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Voice Alert Rule"
  alert_type                = "voice"
  expression                = "((probDetail != \"\") && (discards >= 1%) && (Auto(latency >= Medium sensitivity)))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_voice" "test" {
  test_name        = "User Acceptance Test - Voice"
  interval         = 120
  alerts_enabled   = true
  target_agent_id  = "2334" #Frankfurt, Germany
  bgp_measurements = true
  use_public_bgp   = true
  agents = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  alert_rules = [thousandeyes_alert_rule.test.id, "921609"]
}