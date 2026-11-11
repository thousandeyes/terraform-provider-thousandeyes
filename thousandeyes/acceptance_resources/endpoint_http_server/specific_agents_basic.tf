data "thousandeyes_endpoint_agent" "selected" {
  agent_name = "__ENDPOINT_AGENT_NAME__"
}

resource "thousandeyes_endpoint_http_server" "test" {
  test_name           = "User Acceptance Test - Endpoint HTTP Server"
  url                 = "https://example.com"
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.selected.agent_id]
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 60
  is_enabled          = false
  protocol            = "icmp"
}
