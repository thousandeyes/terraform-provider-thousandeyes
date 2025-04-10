data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
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
  test_name        = "User Acceptance Test - API"
  interval         = 120
  alerts_enabled   = true
  url              = "https://www.thousandeyes.com"
  agents           = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  alert_rules      = [data.thousandeyes_alert_rule.def_alert_rule.id, thousandeyes_alert_rule.alert-rule-http-test.id]
  bgp_measurements = false
  target_time      = 30
  time_limit       = 90
  requests {
    name                  = "Step 1 - GET Request"
    url                   = "https://api.stg.thousandeyes.com/v6/status.json"
    method                = "get"
    auth_type             = "basic"
    username              = "your_username"
    password              = "your_password"
    client_authentication = "basic-auth-header"

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
      value    = "username"
    }
  }
  requests {
    name                  = "Step 2 - Update Profile"
    url                   = "https://api.example.com/profile/update"
    method                = "post"
    auth_type             = "basic"
    username              = "your_username"
    password              = "your_password"
    client_authentication = "basic-auth-header"

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
      value    = "200"
    }

    assertions {
      name     = "response-body"
      operator = "includes"
      value    = "update successful"
    }

    wait_time_ms = 1000 # wait 1 second before next step (if any)
  }

}
