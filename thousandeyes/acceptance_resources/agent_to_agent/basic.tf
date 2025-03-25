provider "thousandeyes" {
  api_endpoint = "https://api.stg.thousandeyes.com/v7"
}

data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_bgp_monitor" "ams_monitor" {
  monitor_name = "Amsterdam-40"
}

# resource "thousandeyes_alert_rule" "alert-rule-agent_to_agent" {
#   rule_name                 = "Custom UAT Agent To Agent Alert Rule"
#   alert_type                = "end-to-end-agent"
#   expression                = "((bothWaysLoss >= 10%) && (bothWaysProbDetail != \"\"))"
#   direction                 = "bidirectional"
#   rounds_violating_out_of   = 1
#   rounds_violating_required = 1
#   minimum_sources           = 1
# }

resource "thousandeyes_agent_to_agent" "test" {
  test_name      = "User Acceptance Test - Aget To Agent"
  interval       = 120
  alerts_enabled = false

  direction       = "bidirectional"
  protocol        = "tcp"
  target_agent_id = "2334" #Frankfurt, Germany
  monitors = [data.thousandeyes_bgp_monitor.ams_monitor.monitor_id]
  # alert_rules = ["921617", thousandeyes_alert_rule.alert-rule-agent_to_agent.id]

  agents {
    agent_id = data.thousandeyes_agent.arg_amsterdam.agent_id
  }

  # alert_rules {
  #   rule_id = 921617 #Agent-To-Agent Default Alert Rule
  # }

  # alert_rules {
  #   rule_id = thousandeyes_alert_rule.alert-rule-agent_to_agent.id
  # }
}
