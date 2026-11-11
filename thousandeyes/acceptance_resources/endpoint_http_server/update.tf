resource "thousandeyes_endpoint_http_server" "test" {
  test_name           = "User Acceptance Test - Endpoint HTTP Server (Updated)"
  url                 = "https://www.thousandeyes.com"
  agent_selector_type = "all-agents"
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 120
  is_enabled          = true
  protocol            = "tcp"
}
