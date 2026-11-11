resource "thousandeyes_endpoint_agent_to_server_dynamic" "test" {
  test_name           = "User Acceptance Test - Endpoint Agent To Server Dynamic (Updated)"
  application         = "zoom"
  agent_selector_type = "all-agents"
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 120
  is_enabled          = true
  protocol            = "tcp"
}
