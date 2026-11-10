resource "thousandeyes_endpoint_agent_to_server" "example" {
  agent_selector_type = "all-agents"
  test_name           = "Example Endpoint Agent-to-Server test"
  server              = "www.thousandeyes.com"
  protocol            = "tcp"
  port                = 443
  interval            = 300
}
