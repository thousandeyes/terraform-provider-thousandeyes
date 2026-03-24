resource "thousandeyes_dashboard" "test_dashboard" {
  description      = "Test Dashboard with All Widget Types"
  title            = "Test Dashboard - All Widgets"
  is_private       = true
  default_timespan {
    duration = 7200
  }
}