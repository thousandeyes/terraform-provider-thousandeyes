data "thousandeyes_agent" "arg_ny" {
  agent_name = "New York, NY"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default API Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "alert-rule-http-test" {
  rule_name                 = "Custom UAT API Alert Rule"
  alert_type                = "api"
  expression                = "((apiCompletion <= 100))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_api" "test" {
  test_name            = "User Acceptance Test - API"
  interval             = 120
  alerts_enabled       = true
  url                  = "https://www.thousandeyes.com"
  agents               = [data.thousandeyes_agent.arg_ny.agent_id]
  alert_rules          = [data.thousandeyes_alert_rule.def_alert_rule.id, thousandeyes_alert_rule.alert-rule-http-test.id]
  network_measurements = false
  bgp_measurements     = false
  target_time          = 30
  time_limit           = 90
  requests {
    name                  = "Step 1 - GET Request"
    url                   = "https://www.thousandeyes.com"
    method                = "get"
    auth_type             = "basic"
    username              = "test_username"
    password              = "test_password"
    client_authentication = "in-body"

    headers {
      key   = "Accept"
      value = "application/json"
    }

    assertions {
      name     = "status-code"
      operator = "is"
      value    = "200"
    }

    assertions {
      name     = "response-body"
      operator = "includes"
      value    = "timestamp"
    }
  }
  requests {
    name                  = "Step 2 - POST OAuth2 request"
    url                   = "https://www.thousandeyes.com"
    method                = "post"
    auth_type             = "oauth2"
    client_id             = "test_client"
    client_secret         = "test_client_secret"
    scope                 = "test_scope"
    token_url             = "https://www.thousandeyes.com"
    client_authentication = "in-body"

    body = jsonencode({
      firstName = "John"
      lastName  = "Doe"
    })

    headers {
      key   = "Content-Type"
      value = "application/json"
    }

    assertions {
      name     = "status-code"
      operator = "is"
      value    = "401"
    }

    assertions {
      name     = "response-body"
      operator = "includes"
      value    = "error"
    }

    wait_time_ms = 1000 
  }

}
