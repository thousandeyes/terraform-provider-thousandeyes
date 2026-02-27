resource "thousandeyes_connector" "test" {
  name   = "UAT - Connector Basic"
  target = "https://example.com/webhooks/thousandeyes"
}
