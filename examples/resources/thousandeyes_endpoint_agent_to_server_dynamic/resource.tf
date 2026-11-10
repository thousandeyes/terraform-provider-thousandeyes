data "thousandeyes_endpoint_agent" "example" {
  agent_name = "Example laptop"
}

resource "thousandeyes_endpoint_agent_to_server_dynamic" "example" {
  test_name           = "Endpoint Webex"
  application         = "webex"
  interval            = 300
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.example.agent_id]
  max_machines        = 25
}
