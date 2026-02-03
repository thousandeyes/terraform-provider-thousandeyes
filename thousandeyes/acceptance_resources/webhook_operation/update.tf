resource "thousandeyes_webhook_operation" "test" {
  name     = "Test Webhook Operation Updated"
  category = "alerts"
  enabled  = false

  path    = "/custom/alerts/v2"
  payload = jsonencode({
    alert_id   = "{{alertId}}"
    test_name  = "{{testName}}"
    alert_type = "{{alertType}}"
    severity   = "{{severity}}"
  })

  headers = [
    {
      name  = "Content-Type"
      value = "application/json"
    },
    {
      name  = "X-Custom-Header"
      value = "custom-value"
    }
  ]

  query_params = jsonencode({
    source = "thousandeyes"
    env    = "production"
  })
}
