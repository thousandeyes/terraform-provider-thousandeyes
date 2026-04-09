resource "thousandeyes_dashboard" "test_dashboard_table_widget" {
  description = "Test Dashboard with Table Widget"
  title       = "Test Dashboard Table Widget"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Table"
    title        = "Test Table Widget"
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

    table_config {
      compare_to_previous_value = true
      row_group_by              = "AGENT"
      column_group_by           = "TEST"
      sort_by                   = "Alphabetical"
      sort_direction            = "Ascending"
      limit                     = 10
    }
  }
}
