resource "thousandeyes_dashboard" "test_dashboard_agent_status_widget" {
  description = "Test Dashboard with Agent Status Widget"
  title       = "Test Dashboard Agent Status Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Agent Status"
    title       = "Test Agent Status Widget"
    visual_mode = "Full"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"

    agent_status_config {
      show = "Owned Agents"
      agent_type = "Enterprise Agents"
    }
  }
}
