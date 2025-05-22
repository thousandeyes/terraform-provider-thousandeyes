resource "thousandeyes_api" "example_api_test" {
  test_name            = "Example for API Test Resource set from Terraform provider"
  interval             = 120
  alerts_enabled       = false
  url                  = "https://www.thousandeyes.com"
  agents               = [3] #Singapore
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
