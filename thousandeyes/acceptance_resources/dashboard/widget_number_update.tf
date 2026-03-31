resource "thousandeyes_dashboard" "test_dashboard_number_widget" {
  description = "Test Dashboard with Number Widget (Updated)"
  title       = "Test Dashboard Number Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "Number"
    title       = "Test Number Widget (Updated)"
    visual_mode = "Full"

    number_cards {
      description  = "Alert Count (Updated)"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 7
        unit  = "Days"
      }
    }

    number_cards {
      description  = "Active Alerts (Updated)"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ACTIVE_ALERT_COUNT"

      measure {
        type = "TOTAL"
      }
    }
  }
}
