resource "thousandeyes_dashboard" "test_dashboard_multi_metric_table_defaults" {
  description = "Test Dashboard Multi Metric Table Defaults"
  title       = "Test Dashboard Multi Metric Table Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Multi Metric Table"
    title       = "Multi Metric Table With Defaults"
    data_source = "ALERTS"
    measure {
      type = "MEAN"
    }
  }
}
