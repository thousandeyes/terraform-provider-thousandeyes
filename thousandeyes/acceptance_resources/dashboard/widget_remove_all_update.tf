resource "thousandeyes_dashboard" "test_dashboard_remove_all_widgets" {
  description = "Test Dashboard for removing all widgets"
  title       = "Test Dashboard Remove All Widgets"
  is_private  = false
  default_timespan {
    duration = 3600
  }
}
