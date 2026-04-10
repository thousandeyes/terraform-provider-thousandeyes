resource "thousandeyes_dashboard" "test_dashboard_table_widget" {
  description = "Test Dashboard with Table Widget (Updated)"
  title       = "Test Dashboard Table Widget (Updated)"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Table"
    title        = "Test Table Widget (Updated)"
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
      compare_to_previous_value = false
      row_group_by              = "TEST"
      column_group_by           = "AGENT"
      limit                     = 20
    }
  }
}
