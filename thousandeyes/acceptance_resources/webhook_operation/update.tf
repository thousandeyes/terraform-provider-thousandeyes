resource "thousandeyes_webhook_operation" "test" {
  name     = "Test Webhook Operation Updated"
  enabled  = true
  category = "alerts"
  status   = "connected"

  path    = "/custom/alerts/v2"
  payload = jsonencode({
    alert_id   = "{{alert.id}}"
    test_name  = "{{alert.test.name}}"
    alert_type = "{{type.id}}"
    severity   = "{{alert.severity.id}}"
  })

  headers {
    name  = "Content-Type"
    value = "application/json"
  }

  headers {
    name  = "X-Custom-Header"
    value = "custom-value"
  }

  query_params = jsonencode({
    source = "thousandeyes"
    env    = "production"
  })

  type = "webhook"
}
