resource "thousandeyes_bgp" "example_bgp_test" {
  test_name      = "Example BGP test set from Terraform provider"
  alerts_enabled = false

  prefix = "163.10.0.0/16"
}
