resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Remove Auth"
  target = "https://example.com/webhooks/thousandeyes/remove-auth"
}
