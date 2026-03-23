resource "thousandeyes_dashboard" "test_dashboard_timeseries_widget" {
  description = "Test Dashboard with Timeseries Widget (Updated)"
  title       = "Test Dashboard Timeseries Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }

  widgets {
    type        = "Time Series: Line"
    title       = "Test Timeseries Widget (Updated)"
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
      group_by = "TEST"
      show_timeseries_overall_baseline = true
      is_timeseries_one_chart_per_line = true
    }
  }
}
