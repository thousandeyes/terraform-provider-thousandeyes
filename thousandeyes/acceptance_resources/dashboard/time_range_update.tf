resource "thousandeyes_dashboard" "test_dashboard_time_range" {
  description = "Updated Test Dashboard with Time Range"
  title       = "Test Dashboard Time Range (Updated)"
  is_private  = true
  default_timespan {
    start = "2026-02-01T00:00:00Z"
    end   = "2026-03-01T00:00:00Z"
  }
}
