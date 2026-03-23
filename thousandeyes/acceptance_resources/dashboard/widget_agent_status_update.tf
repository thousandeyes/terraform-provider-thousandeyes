resource "thousandeyes_dashboard" "test_dashboard_agent_status_widget" {
  description = "Test Dashboard with Agent Status Widget (Updated)"
  title       = "Test Dashboard Agent Status Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }

  widgets {
    type        = "Agent Status"
    title       = "Agent Status Widget (Updated)"
    visual_mode = "Full"
    data_source = "ENDPOINT_AGENTS"

    agent_status_config {
      show = "Owned Agents"
      agent_type = "Endpoint Agents"
    }
  }
}
