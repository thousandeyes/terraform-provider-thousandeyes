resource "thousandeyes_endpoint_agent_to_server" "test" {
  test_name           = "Endpoint Agent-to-Server acceptance test"
  server              = "www.thousandeyes.com"
  interval            = 300
  protocol            = "icmp"
  agent_selector_type = "all-agents"
}
