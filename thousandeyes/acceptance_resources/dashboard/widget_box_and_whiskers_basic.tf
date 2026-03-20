resource "thousandeyes_dashboard" "test_dashboard_box_and_whiskers_widget" {
  description = "Test Dashboard with Box and Whiskers Widget"
  title       = "Test Dashboard Box and Whiskers Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Box and Whiskers"
    title       = "Test Box and Whiskers Widget"
    visual_mode = "Full"
    data_source = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    box_and_whiskers_config {
      group_by = "COUNTRY"
    }
  }
}
