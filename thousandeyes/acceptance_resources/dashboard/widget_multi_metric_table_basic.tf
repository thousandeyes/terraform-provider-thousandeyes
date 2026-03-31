resource "thousandeyes_dashboard" "test_dashboard_multi_metric_table_widget" {
  description = "Test Dashboard with Multi Metric Table Widget"
  title       = "Test Dashboard Multi Metric Table Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Multi Metric Table"
    title       = "Test Multi Metric Table Widget"
    visual_mode = "Full"

    multi_metric_table_config {
      compare_to_previous_value = true
      row_group_by              = "TESTS"
      limit                     = 10
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT"
      measure {
        type = "MEAN"
      }
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT"
      measure {
        type = "TOTAL"
      }
    }
  }
}
