data "thousandeyes_dashboard_filter" "operations_core_services" {
  name = "Operations - Core Services"
}

resource "thousandeyes_dashboard" "example" {
  title              = "Operations Overview"
  description        = "Terraform-managed dashboard for shared service health views"
  is_private         = false
  global_filter_id   = data.thousandeyes_dashboard_filter.operations_core_services.filter_id
  is_global_override = true

  default_timespan {
    duration = 7200
  }

  widgets {
    type        = "Agent Status"
    title       = "Enterprise Agent Status"
    visual_mode = "Full"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"

    agent_status_config {
      show       = "Owned Agents"
      agent_type = "Enterprise Agents"
    }
  }

  widgets {
    type         = "Time Series: Line"
    title        = "Active Agent Alerts"
    visual_mode  = "Full"
    data_source  = "ALERTS"
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
      group_by = "AGENT"
    }
  }

  widgets {
    type        = "Number"
    title       = "Service Summary"
    visual_mode = "Full"

    number_cards {
      description  = "HTTP availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }
    }
  }
}
