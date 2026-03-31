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
      description  = "CEA Availability (Updated)"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 7
        unit  = "Days"
      }
    }

    number_cards {
      description  = "Agent Alerts (Updated)"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"

      measure {
        type = "TOTAL"
      }
    }
  }
}
