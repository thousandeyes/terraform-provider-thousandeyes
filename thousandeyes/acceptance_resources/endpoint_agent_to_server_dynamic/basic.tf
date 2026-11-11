resource "thousandeyes_endpoint_agent_to_server_dynamic" "test" {
  test_name           = "User Acceptance Test - Endpoint Agent To Server Dynamic"
  application         = "webex"
  agent_selector_type = "all-agents"
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 60
  is_enabled          = false
  protocol            = "icmp"
}
