resource "thousandeyes_dashboard" "test_dashboard_number_widget" {
  description = "Test Dashboard with Number Widget"
  title       = "Test Dashboard Number Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Number"
    title       = "Test Number Widget"
    visual_mode = "Full"
    data_source = "ALERTS"

    number_cards {
      description  = "Alert Count"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 1
        unit  = "Days"
      }
    }

    number_cards {
      description  = "Active Alerts"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ACTIVE_ALERT_COUNT"

      measure {
        type = "MEAN"
      }
    }
  }
}
