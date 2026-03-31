resource "thousandeyes_dashboard" "test_dashboard_multi_metric_table_widget" {
  description = "Test Dashboard with Multi Metric Table Widget (Updated)"
  title       = "Test Dashboard Multi Metric Table Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "Multi Metric Table"
    title       = "Test Multi Metric Table Widget (Updated)"
    visual_mode = "Full"
    data_source = "ALERTS"

    measure {
      type = "MEAN"
    }

    multi_metric_table_config {
      compare_to_previous_value = false
      row_group_by              = "TESTS"
      limit                     = 20
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"
      measure {
        type = "MEAN"
      }
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ACTIVE_ALERT_COUNT"
      measure {
        type = "MEAN"
      }
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"
      measure {
        type = "TOTAL"
      }
    }
  }
}
