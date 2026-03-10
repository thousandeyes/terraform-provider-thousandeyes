resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Bearer"
  target = "https://example.com/webhooks/thousandeyes/auth-bearer"

  authentication {
    type  = "bearer-token"
    token = "test-bearer-token-value"
  }
}
