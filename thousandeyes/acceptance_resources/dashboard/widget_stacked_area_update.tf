resource "thousandeyes_dashboard" "test_dashboard_stacked_area_widget" {
  description = "Test Dashboard with Stacked Area Widget (Updated)"
  title       = "Test Dashboard Stacked Area Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "Time Series: Stacked Area"
    title       = "Test Stacked Area Widget (Updated)"
    visual_mode = "Full"
    data_source = "CLOUD_NATIVE_MONITORING" // Cloud insights
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    stacked_area_config {
      group_by = "CLOUD_NATIVE_MONITORING-ACCOUNT"
    }
  }
}
