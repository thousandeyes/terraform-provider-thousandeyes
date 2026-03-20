resource "thousandeyes_dashboard" "test_dashboard" {
  description        = "Test Dashboard Description"
  title              = "Test Dashboard"
  is_private         = false
  is_global_override = false
  default_timespan {
    duration = 3600
  }
}