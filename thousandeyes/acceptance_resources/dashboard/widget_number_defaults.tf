resource "thousandeyes_dashboard" "test_dashboard_number_defaults" {
  description = "Test Dashboard Number Defaults"
  title       = "Test Dashboard Number Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Number"
    title       = "Number With Defaults"
    data_source = "ALERTS"

    number_cards {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"
      measure {
        type = "MEAN"
      }
    }
  }
}
