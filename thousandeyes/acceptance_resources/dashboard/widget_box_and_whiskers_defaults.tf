resource "thousandeyes_dashboard" "test_dashboard_box_and_whiskers_defaults" {
  description = "Test Dashboard Box and Whiskers Defaults"
  title       = "Test Dashboard Box and Whiskers Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type         = "Box and Whiskers"
    title        = "Box and Whiskers With Defaults"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"
    measure {
      type = "MEAN"
    }
  }
}
