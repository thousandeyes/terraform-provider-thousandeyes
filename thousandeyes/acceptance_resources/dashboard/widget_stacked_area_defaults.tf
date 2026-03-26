resource "thousandeyes_dashboard" "test_dashboard_stacked_area_defaults" {
  description = "Test Dashboard Stacked Area Defaults"
  title       = "Test Dashboard Stacked Area Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type         = "Time Series: Stacked Area"
    title        = "Stacked Area With Defaults"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"
    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    stacked_area_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }
}
