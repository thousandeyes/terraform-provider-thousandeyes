data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "UAT - Agent To Server ICMP"
  interval       = 120
  alerts_enabled = false
  server         = "api.stg.thousandeyes.com"
  protocol       = "icmp"
  # port is not specified for ICMP
  agents         = [data.thousandeyes_agent.amsterdam.agent_id]
}

