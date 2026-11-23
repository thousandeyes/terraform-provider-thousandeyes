data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

resource "thousandeyes_http_server" "test" {
  test_name      = "User Acceptance Test - HTTP Default Alert Rules"
  description    = "Updated without explicit alert rules"
  interval       = 120
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.arg_frankfurt.agent_id]
}
