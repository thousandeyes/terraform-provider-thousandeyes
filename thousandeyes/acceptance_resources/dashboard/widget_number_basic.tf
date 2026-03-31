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

    number_cards {
      description  = "CEA Availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 1
        unit  = "Days"
      }
    }

    number_cards {
      description  = "Agent Alerts"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"

      measure {
        type = "TOTAL"
      }
    }
  }
}
