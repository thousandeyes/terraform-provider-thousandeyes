data "thousandeyes_endpoint_agent" "selected" {
  agent_name = "TE_ENDPOINT_AGENT_NAME"
}

resource "thousandeyes_endpoint_http_server" "test" {
  test_name           = "User Acceptance Test - Endpoint HTTP Server (Updated)"
  url                 = "https://www.thousandeyes.com"
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.selected.agent_id]
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 120
  is_enabled          = true
  protocol            = "tcp"
}
