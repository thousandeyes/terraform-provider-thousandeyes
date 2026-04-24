resource "thousandeyes_dashboard" "test_dashboard_alert_list_defaults" {
  description = "Test Dashboard with Alert List default alert types"
  title       = "Test Dashboard Alert List Default Alert Types"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Alert List"
    title       = "Alert List Default Alert Types"
    visual_mode = "Full"
    data_source = "ALERTS"

    alert_list_config {
      limit_to = 15

      active_within_value = 7
      active_within_unit  = "Days"
    }
  }
}
