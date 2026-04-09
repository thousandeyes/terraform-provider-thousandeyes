resource "thousandeyes_dashboard" "test_dashboard_stacked_bar_chart_widget" {
  description = "Test Dashboard with Stacked Bar Chart Widget (Updated)"
  title       = "Test Dashboard Stacked Bar Chart Widget (Updated)"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Bar Chart: Stacked"
    title        = "Stacked Bar Chart Widget (Updated)"
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

    stacked_bar_chart_config {
      axis_group_by           = "CLOUD_NATIVE_MONITORING-ACCOUNT"
      sort_by                 = "Alphabetical"
      sort_direction          = "Ascending"
      limit                   = 5
      show_labels             = false
      is_horizontal_bar_chart = false
    }
  }
}
