resource "thousandeyes_webhook_operation" "test" {
  name     = "Test Webhook Operation"
  category = "alerts"
  status   = "pending"

  path    = "/custom/alerts"
  payload = jsonencode({
    alert_id   = "{{alert.id}}"
    test_name  = "{{alert.test.name}}"
    alert_type = "{{type.id}}"
  })

  headers {
    name  = "Content-Type"
    value = "application/json"
  }

  query_params = jsonencode({
    source = "thousandeyes"
  })
}
