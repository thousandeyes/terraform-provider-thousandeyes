resource "thousandeyes_dashboard" "test_dashboard_pie_chart_defaults" {
  description = "Test Dashboard Pie Chart Defaults"
  title       = "Test Dashboard Pie Chart Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type         = "Pie Chart"
    title        = "Pie Chart With Defaults"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"
    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }
    pie_chart_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }
}
