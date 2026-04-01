resource "thousandeyes_dashboard" "example" {
  description      = "Example Dashboard description"
  title            = "Example Dashboard"
  is_private       = true
  default_timespan {
    duration = 7200
  }

  widgets {
    type        = "Agent Status"
    title       = "Agent Status Widget"
    visual_mode = "Full"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"

    agent_status_config {
      show       = "All Assigned Agents"
      agent_type = "Enterprise Agents"
    }
  }

  widgets {
    type         = "Map"
    title        = "Map Widget"
    visual_mode  = "Full"
    is_embedded  = false
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_BGP"
    data_source  = "ALERTS"

    measure {
      type = "MEAN"
    }

    geo_map_config {
      min_scale           = 0
      max_scale           = 100
      group_by            = "COUNTRY"
      is_geo_map_per_test = true
    }
  }

  widgets {
    type         = "Time Series: Line"
    title        = "Timeseries Line Widget"
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
      group_by                         = "AGENT"
      show_timeseries_overall_baseline = true
      is_timeseries_one_chart_per_line = true
    }
  }

  widgets {
    type         = "Time Series: Stacked Area"
    title        = "Stacked Area Widget"
    visual_mode  = "Full"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    stacked_area_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }

  widgets {
    type         = "Pie Chart"
    title        = "Pie Chart Widget"
    visual_mode  = "Full"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    pie_chart_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }

  widgets {
    type         = "Box and Whiskers"
    title        = "Box and Whiskers Widget"
    visual_mode  = "Full"
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

    box_and_whiskers_config {
      group_by = "COUNTRY"
    }
  }
}
