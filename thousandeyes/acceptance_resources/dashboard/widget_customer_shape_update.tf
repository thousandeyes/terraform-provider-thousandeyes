resource "thousandeyes_dashboard" "test_dashboard_customer_shape_widget_types" {
  title       = "Customer Dashboard A Updated"
  description = ""
  is_private  = false

  default_timespan {
    duration = 86400
  }

  widgets {
    type        = "Number"
    title       = "Network Status Summary"
    visual_mode = "Full"

    number_cards {
      description = "SDWAN Packet Loss"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "NET_LOSS"
      metric_group = "AGENT_TO_SERVER"
      max_scale   = 10

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }

    number_cards {
      description = "Internet Packet Loss"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "NET_LOSS"
      metric_group = "AGENT_TO_SERVER"
      max_scale   = 10

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }

    number_cards {
      description = "DNS Resolution Availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "DNS_SERVER_AVAILABILITY"
      metric_group = "DNS"

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }
  }

  widgets {
    type        = "Number"
    title       = "Application Status Summary"
    visual_mode = "Full"

    number_cards {
      description = "App A Availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "WEB_AVAILABILITY"
      metric_group = "HTTP_SERVER"
      min_scale   = 90

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }

    number_cards {
      description = "App B Availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "WEB_AVAILABILITY"
      metric_group = "HTTP_SERVER"
      min_scale   = 90

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }

    number_cards {
      description = "App C Availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric      = "WEB_AVAILABILITY"
      metric_group = "HTTP_SERVER"
      min_scale   = 90

      measure {
        type = "MEAN"
      }

      fixed_timespan {
        value = 2
        unit  = "Hours"
      }
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "Application Availability Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "WEB_AVAILABILITY"
    metric_group = "HTTP_SERVER"

    measure {
      type = "MEAN"
    }


    timeseries_config {
      group_by = "TEST_TAG"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "Network Packet Loss Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "NET_LOSS"
    metric_group = "AGENT_TO_SERVER"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 2
      unit  = "Hours"
    }


    timeseries_config {
      group_by = "TEST_TAG"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "Application Latency Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "NET_LATENCY"
    metric_group = "AGENT_TO_SERVER"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 2
      unit  = "Hours"
    }


    timeseries_config {
      group_by = "TEST_TAG"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "DNS Availability Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "DNS_SERVER_AVAILABILITY"
    metric_group = "DNS"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 2
      unit  = "Hours"
    }

    timeseries_config {
      group_by = "SERVER"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }

  widgets {
    type        = "Color Grid"
    title       = "Network Packet Loss by Test"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "NET_LOSS"
    metric_group = "AGENT_TO_SERVER"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 2
      unit  = "Hours"
    }


    color_grid_config {
      max_scale      = 10
      cards          = "TEST"
      group_cards_by = "TEST_TAG"
      columns        = 2
    }
  }

  widgets {
    type        = "Agent Status"
    title       = "Agent Status"
    visual_mode = "Full"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"

    agent_status_config {
      show       = "All Assigned Agents"
      agent_type = "Enterprise Agents"
    }
  }

  widgets {
    type        = "List"
    title       = "Recent Events List"
    visual_mode = "Full"
    data_source = "EVENT_DETECTION"

    measure {
      type = "MEAN"
    }

    list_config {
      active_within_value = 7
      active_within_unit  = "Days"
    }
  }

  widgets {
    type        = "Alert List"
    title       = "Recent Alerts"
    visual_mode = "Full"
    data_source = "ALERTS"

    alert_list_config {
      alert_types = ["API", "DNS Server"]
      limit_to            = 15
      active_within_value = 7
      active_within_unit  = "Days"
    }
  }
}