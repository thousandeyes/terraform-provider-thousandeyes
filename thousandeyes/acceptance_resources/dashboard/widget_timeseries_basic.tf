resource "thousandeyes_dashboard" "test_dashboard_timeseries_widget" {
  description = "Test Dashboard with Timeseries Widget"
  title       = "Test Dashboard Timeseries Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Time Series: Line"
    title       = "Test Timeseries Widget"
    visual_mode = "Full"
    data_source = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "TOTAL"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    timeseries_config {
      group_by = "AGENT"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }
}
