resource "thousandeyes_webhook_operation" "example_webhook_operation" {
  name     = "Example Webhook Operation"
  category = "alerts"
  status   = "pending"
  enabled  = true
  type     = "webhook"

  path    = "/custom/alerts/endpoint"
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
    name  = "Authorization"
    value = "Bearer YOUR_TOKEN_HERE"
  }

  query_params = jsonencode({
    source = "thousandeyes"
    env    = "production"
  })
}
