resource "thousandeyes_dashboard" "test_dashboard_time_range" {
  description = "Test Dashboard with Time Range"
  title       = "Test Dashboard Time Range"
  is_private  = false
  default_timespan {
    start = "2026-01-01T00:00:00Z"
    end   = "2026-02-01T00:00:00Z"
  }
}
