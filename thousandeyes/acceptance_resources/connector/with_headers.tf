resource "thousandeyes_connector" "test_headers" {
  name   = "UAT - Connector With Headers"
  target = "https://example.com/webhooks/thousandeyes/headers"

  headers {
    name  = "X-Custom-Header"
    value = "custom-value-1"
  }

  headers {
    name  = "X-Another-Header"
    value = "custom-value-2"
  }
}
