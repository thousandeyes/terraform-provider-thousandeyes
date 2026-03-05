resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Other Token"
  target = "https://example.com/webhooks/thousandeyes/auth-other"

  authentication {
    type  = "other-token"
    token = "test-other-token-value"
  }
}
