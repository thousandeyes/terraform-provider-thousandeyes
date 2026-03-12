resource "thousandeyes_dashboard" "test_dashboard" {
  description = "Test Dashboard Description"
  title       = "Test Dashboard"
  global_filter_id = 123
  is_private  = false
  default_timespan {
    duration = 3600
  }
}