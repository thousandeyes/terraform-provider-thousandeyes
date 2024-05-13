data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  severity                  = "MAJOR"
  rule_name                 = "Agent To Server Alert Rule Test"
  alert_type                = "End-to-End (Server)"
  expression                = "((loss >= 50%) || (probDetail != \"\") || (avgLatency >= 200 ms))"
  minimum_sources           = 2
  rounds_violating_required = 3
  rounds_violating_out_of   = 4
}

resource "thousandeyes_agent_to_server" "agent_to_server_test" {
  test_name        = "Agent To Server Test"
  interval         = 300
  alerts_enabled   = true
  server           = "api.stg.thousandeyes.com"
  protocol         = "TCP"
  port             = 443
  enabled          = true
  bgp_measurements = true
  use_public_bgp   = true
  mtu_measurements = true

  agents {
    agent_id = data.thousandeyes_agent.amsterdam.agent_id
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.test.id
  }
}
