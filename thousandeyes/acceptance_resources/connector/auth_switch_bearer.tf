resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Switch"
  target = "https://example.com/webhooks/thousandeyes/auth-switch"

  authentication {
    type  = "bearer-token"
    token = "new-bearer-token"
  }
}
