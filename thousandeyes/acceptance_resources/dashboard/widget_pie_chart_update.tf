resource "thousandeyes_dashboard" "test_dashboard_pie_chart_widget" {
  description = "Test Dashboard with Pie Chart Widget (Updated)"
  title       = "Test Dashboard Pie Chart Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "Pie Chart"
    title       = "Test Pie Chart Widget (Updated)"
    visual_mode = "Full"
    data_source = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    pie_chart_config {
      group_by = "CLOUD_NATIVE_MONITORING-ACCOUNT"
    }
  }
}
