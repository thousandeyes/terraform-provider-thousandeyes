data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

data "thousandeyes_account_group" "share_target" {
  name = "API Team Non Default Test"
}

resource "thousandeyes_http_server" "test" {
  test_name      = "User Acceptance Test - HTTP Shared"
  interval       = 120
  alerts_enabled = false
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.arg_frankfurt.agent_id]
}
