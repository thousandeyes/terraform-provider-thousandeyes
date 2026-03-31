resource "thousandeyes_dashboard" "test_dashboard_preserve_unmanaged" {
  description = "Test Dashboard with Unmanaged Widget Preservation"
  title       = "Test Dashboard Preserve Unmanaged"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Time Series: Line"
    title       = "Managed Timeseries Widget"
    visual_mode = "Full"
    data_source = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "TOTAL"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    timeseries_config {
      group_by                          = "AGENT"
      show_timeseries_overall_baseline  = false
      is_timeseries_one_chart_per_line  = false
    }
  }
}
