data "thousandeyes_bgp_monitor" "example_bgp_monitor" {
  monitor_name = "test-monitor"
}

resource "thousandeyes_bgp" "example_bgp_test" {
  test_name      = "Example BGP test set from Terraform provider."
  alerts_enabled = false
  bgp_monitors   = [data.thousandeyes_bgp_monitor.example_bgp_monitor.monitor_id]
  prefix         = "163.10.0.0/16"
}
