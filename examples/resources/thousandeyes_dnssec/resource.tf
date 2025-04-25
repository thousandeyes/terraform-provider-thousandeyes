resource "thousandeyes_dnssec" "example_dnssec_test" {
  test_name      = "Example DNSSEC test set from Terraform provider"
  interval       = 120
  alerts_enabled = false
  domain         = "www.thousandeyes.com ANY"
  agents         = ["3"] # Singapore
}
