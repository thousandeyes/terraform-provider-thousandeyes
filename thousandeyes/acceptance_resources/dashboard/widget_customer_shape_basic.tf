resource "thousandeyes_dashboard" "test_dashboard_customer_shape_widget_types" {
  title       = "Customer Dashboard A"
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

      filter {
        property = "TEST"
        values   = ["1978846", "5819450", "6856213"]
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

      filter {
        property = "TEST"
        values   = ["8859181", "2241483", "3091431"]
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

      filter {
        property = "TEST"
        values = ["2667547", "7184517", "226817"]
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


      filter {
        property = "TEST"
        values   = ["2241483"]
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

      filter {
        property = "TEST"
        values   = ["8859181"]
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

      filter {
        property = "TEST"
        values   = ["3091431"]
      }
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "Application Availability Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric_group = "HTTP_SERVER"
    metric      = "WEB_AVAILABILITY"

    measure {
      type = "MEAN"
    }

    filter {
      property = "TEST_TAG"
      values   = ["876442940518191", "226261155961349", "400304130254252"]
    }

    filter {
      property = "TEST"
      values   = ["8859181", "2241483", "3091431"]
    }

    timeseries_config {
      group_by = "TEST_TAG"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }

  widgets {
    type        = "Time Series: Line"
    title       = "HTTP Response Time Timeline"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "WEB_TTFB"
    metric_group = "HTTP_SERVER"


    measure {
      type = "MEAN"
    }

    filter {
      property = "TEST_TAG"
      values   = ["876442940518191", "226261155961349", "400304130254252"]
    }

    filter {
      property = "TEST"
      values   = ["8859181", "2241483", "3091431"]
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

    filter {
      property = "TEST_TAG"
      values = [
        "265519578772930",
        "640697525195518",
        "876442940518191",
        "226261155961349",
        "400304130254252",
        "146839052746727",
      ]
    }

    filter {
      property = "TEST"
      values = [
        "8859181",
        "2241483",
        "3091431",
        "2667547",
        "7184517",
        "226817",
        "1978846",
        "5819450",
        "6856213",
        "5974575",
      ]
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

    filter {
      property = "TEST_TAG"
      values   = ["876442940518191", "226261155961349", "400304130254252"]
    }

    filter {
      property = "TEST"
      values   = ["8859181", "2241483", "3091431"]
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

    filter {
      property = "TEST"
      values = ["2667547", "7184517", "226817"]
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

    filter {
      property = "TEST_TAG"
      values   = ["876442940518191", "640697525195518", "888980457574315", "146839052746727"]
    }

    filter {
      property = "TEST"
      values = [
        "8859181",
        "2241483",
        "3091431",
        "2667547",
        "7184517",
        "226817",
        "1978846",
        "5819450",
        "6856213",
        "5974575",
      ]
    }

    color_grid_config {
      max_scale      = 10
      cards          = "TEST"
      group_cards_by = "TEST_TAG"
      columns        = 2
    }
  }

  widgets {
    type        = "Color Grid"
    title       = "Application Availability by Test"
    visual_mode = "Full"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric      = "WEB_AVAILABILITY"
    metric_group = "HTTP_SERVER"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 2
      unit  = "Hours"
    }

    filter {
      property = "TEST_TAG"
      values   = ["876442940518191", "226261155961349", "400304130254252"]
    }

    filter {
      property = "TEST"
      values   = ["8859181", "2241483", "3091431"]
    }

    color_grid_config {
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

    filter {
      property = "AGENT"
      values   = ["9381026"]
    }

    agent_status_config {
      show       = "All Assigned Agents"
      agent_type = "Enterprise Agents"
    }
  }
}