resource "thousandeyes_dashboard" "test_dashboard_box_and_whiskers_widget" {
  description = "Test Dashboard with Box and Whiskers Widget (Updated)"
  title       = "Test Dashboard Box and Whiskers Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "Box and Whiskers"
    title       = "Test Box and Whiskers Widget (Updated)"
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
      group_by = "CONTINENT"
    }
  }
}
