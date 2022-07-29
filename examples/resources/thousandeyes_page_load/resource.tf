resource "thousandeyes_page_load" "test" {
  test_name      = "Example Page Load test set from Terraform provider"
  alerts_enabled = false
  url            = "https://www.thousandeyes.com"

  interval      = 120
  http_interval = 120

  agents {
    agent_id = 3 # Singapur
  }
}
