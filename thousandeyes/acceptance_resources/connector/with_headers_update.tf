resource "thousandeyes_connector" "test_headers" {
  name   = "UAT - Connector With Headers (Updated)"
  target = "https://example.com/webhooks/thousandeyes/headers"

  headers {
    name  = "X-Updated-Header"
    value = "updated-value"
  }
}
