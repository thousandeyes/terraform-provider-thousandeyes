resource "thousandeyes_connector" "example" {
  name   = "Example Connector"
  target = "https://webhook.site/your-id"

  headers {
    name  = "Content-Type"
    value = "application/json"
  }

  authentication {
    type     = "basic"
    username = "user"
    password = "pass"
  }
}
