resource "thousandeyes_dashboard" "test_dashboard_timeseries_defaults" {
  description = "Test Dashboard Timeseries Defaults"
  title       = "Test Dashboard Timeseries Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type         = "Time Series: Line"
    title        = "Timeseries With Defaults"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"
    measure {
      type = "TOTAL"
    }
  }
}
