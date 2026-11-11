data "thousandeyes_endpoint_agent" "test" {
  agent_name = "__ENDPOINT_AGENT_NAME__"
}

resource "thousandeyes_endpoint_agent_to_server" "test" {
  test_name           = "Endpoint Agent-to-Server acceptance test updated"
  server              = "www.thousandeyes.com"
  interval            = 300
  protocol            = "icmp"
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.test.agent_id]
}
