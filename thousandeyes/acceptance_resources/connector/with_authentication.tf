resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Basic"
  target = "https://example.com/webhooks/thousandeyes/auth-basic"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
