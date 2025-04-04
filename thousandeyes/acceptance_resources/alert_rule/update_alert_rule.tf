data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  severity                  = "major"
  rule_name                 = "Agent To Server Alert Rule Test"
  alert_type                = "end-to-end-server"
  expression                = "((loss >= 50%) || (probDetail != \"\") || (avgLatency >= 200 ms))"
  minimum_sources           = 2
  rounds_violating_required = 4
  rounds_violating_out_of   = 4
}

resource "thousandeyes_agent_to_server" "agent_to_server_test" {
  test_name        = "Agent To Server Test"
  interval         = 300
  alerts_enabled   = true
  server           = "api.stg.thousandeyes.com"
  protocol         = "tcp"
  enabled          = true
  bgp_measurements = true
  use_public_bgp   = true
  mtu_measurements = true
  agents           = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules      = [thousandeyes_alert_rule.test.id]
}
