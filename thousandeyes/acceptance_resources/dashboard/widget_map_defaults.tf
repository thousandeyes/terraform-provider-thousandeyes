resource "thousandeyes_dashboard" "test_dashboard_map_defaults" {
  description = "Test Dashboard Map Defaults"
  title       = "Test Dashboard Map Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type         = "Map"
    title        = "Map With Defaults"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_BGP"
    measure {
      type = "MEAN"
    }
  }
}
