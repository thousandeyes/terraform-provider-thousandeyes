resource "thousandeyes_dashboard" "test_dashboard_alert_list_widget" {
  description = "Test Dashboard with Alert List Widget (Updated)"
  title       = "Test Dashboard Alert List Widget (Updated)"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Alert List"
    title       = "Alert List Widget (Updated)"
    visual_mode = "Full"
    data_source = "ALERTS"

    alert_list_config {
      alert_types = ["Web - HTTP Server"]
      limit_to    = 20

      active_within_value = 14
      active_within_unit  = "Days"
    }
  }
}
