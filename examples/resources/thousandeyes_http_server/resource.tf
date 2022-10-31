resource "thousandeyes_http_server" "example_http_server_test" {
  test_name      = "Example HTTP test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  url = "https://www.thousandeyes.com"

  agents {
    agent_id = 3 # Singapore
  }
}
