resource "thousandeyes_webhook_operation" "test" {
  name     = "Test Webhook Operation"
  category = "alerts"
  enabled  = true

  path    = "/custom/alerts"
  payload = jsonencode({
    alert_id   = "{{alertId}}"
    test_name  = "{{testName}}"
    alert_type = "{{alertType}}"
  })

  headers = [
    {
      name  = "Content-Type"
      value = "application/json"
    }
  ]

  query_params = jsonencode({
    source = "thousandeyes"
  })
}
