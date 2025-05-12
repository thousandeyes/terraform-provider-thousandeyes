resource "thousandeyes_page_load" "example_page_load_test" {
  test_name      = "Example Page Load test set from Terraform provider"
  alerts_enabled = false
  url            = "https://www.thousandeyes.com"
  interval       = 120
  http_interval  = 120
  agents         = ["3"] # Singapore
}
