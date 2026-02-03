resource "thousandeyes_webhook_operation" "example_webhook_operation" {
  name     = "Example Webhook Operation"
  category = "alerts"
  enabled  = true

  path    = "/custom/alerts/endpoint"
  payload = jsonencode({
    alert_id   = "{{alertId}}"
    test_name  = "{{testName}}"
    alert_type = "{{alertType}}"
    severity   = "{{severity}}"
    timestamp  = "{{timestamp}}"
  })

  headers = [
    {
      name  = "Content-Type"
      value = "application/json"
    },
    {
      name  = "Authorization"
      value = "Bearer YOUR_TOKEN_HERE"
    }
  ]

  query_params = jsonencode({
    source = "thousandeyes"
    env    = "production"
  })
}
