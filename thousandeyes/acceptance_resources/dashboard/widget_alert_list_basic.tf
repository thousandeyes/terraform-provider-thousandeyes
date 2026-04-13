resource "thousandeyes_dashboard" "test_dashboard_alert_list_widget" {
  description = "Test Dashboard with Alert List Widget"
  title       = "Test Dashboard Alert List Widget"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Alert List"
    title       = "Alert List Widget"
    visual_mode = "Full"
    data_source = "ALERTS"

    alert_list_config {
      alert_types = ["API", "DNS Server"]
      limit_to    = 15

      active_within_value = 7
      active_within_unit  = "Days"
    }
  }
}
