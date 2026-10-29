data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_api" "test" {
  test_name            = "User Acceptance Test - API Without Sharing"
  interval             = 600
  agents               = [data.thousandeyes_agent.amsterdam.agent_id]
  url                  = "https://api.thousandeyes.com/status"
  alerts_enabled       = false
  network_measurements = false
  bgp_measurements     = false

  requests {
    name                 = "step1"
    url                  = "https://api.thousandeyes.com/status"
    method               = "get"
    collect_api_response = true

    headers {
      key   = "Content-Type"
      value = "application/json"
    }

    assertions {
      name     = "status-code"
      operator = "is"
      value    = "200"
    }
  }
}
