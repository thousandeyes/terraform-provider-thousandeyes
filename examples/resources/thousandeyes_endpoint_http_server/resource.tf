data "thousandeyes_endpoint_agent" "example" {
  agent_name = "Example laptop"
}

resource "thousandeyes_endpoint_http_server" "example" {
  test_name           = "Endpoint HTTP Server"
  url                 = "https://example.com"
  interval            = 300
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.example.agent_id]
  max_machines        = 25
}
