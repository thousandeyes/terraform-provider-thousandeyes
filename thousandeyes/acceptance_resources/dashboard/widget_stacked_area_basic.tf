resource "thousandeyes_dashboard" "test_dashboard_stacked_area_widget" {
  description = "Test Dashboard with Stacked Area Widget"
  title       = "Test Dashboard Stacked Area Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Time Series: Stacked Area"
    title       = "Test Stacked Area Widget"
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
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }
}
