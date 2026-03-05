resource "thousandeyes_connector" "test_full" {
  name   = "UAT - Connector Full"
  target = "https://example.com/webhooks/thousandeyes/full"

  headers {
    name  = "X-Correlation-ID"
    value = "abc-123"
  }

  authentication {
    type     = "basic"
    username = "admin"
    password = "secret"
  }
}
