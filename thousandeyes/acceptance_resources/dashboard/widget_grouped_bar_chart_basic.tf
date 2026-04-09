resource "thousandeyes_dashboard" "test_dashboard_grouped_bar_chart_widget" {
  description = "Test Dashboard with Grouped Bar Chart Widget"
  title       = "Test Dashboard Grouped Bar Chart Widget"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Bar Chart: Grouped"
    title        = "Grouped Bar Chart Widget"
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

    grouped_bar_chart_config {
      group_by                = "COUNTRY"
      axis_group_by           = "TEST"
      sort_by                 = "Alphabetical"
      sort_direction          = "Ascending"
      limit                   = 12
      show_labels             = true
      is_horizontal_bar_chart = false
    }
  }
}
