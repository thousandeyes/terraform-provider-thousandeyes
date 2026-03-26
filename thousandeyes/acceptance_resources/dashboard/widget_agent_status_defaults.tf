resource "thousandeyes_dashboard" "test_dashboard_agent_status_defaults" {
  description = "Test Dashboard Agent Status Defaults"
  title       = "Test Dashboard Agent Status Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Agent Status"
    title       = "Agent Status With Defaults"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"
  }
}
