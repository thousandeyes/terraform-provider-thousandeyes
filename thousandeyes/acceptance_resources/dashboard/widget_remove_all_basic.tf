resource "thousandeyes_dashboard" "test_dashboard_remove_all_widgets" {
  description = "Test Dashboard for removing all widgets"
  title       = "Test Dashboard Remove All Widgets"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Time Series: Line"
    title       = "Widget One"
    visual_mode = "Full"
    data_source = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "TOTAL"
    }

    timeseries_config {
      group_by = "AGENT"
    }
  }
  widgets {
    type        = "Box and Whiskers"
    title       = "Widget Two"
    visual_mode = "Full"
    data_source = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "MEAN"
    }

    box_and_whiskers_config {
      group_by = "AGENT"
    }
  }
}
