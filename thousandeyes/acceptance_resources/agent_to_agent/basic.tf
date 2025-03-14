data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "alert-rule-agent_to_agent" {
  rule_name                 = "Custom UAT Agent To Agent Alert Rule"
  alert_type                = "End-to-End (Agent)"
  expression                = "((bothWaysLoss >= 10%) && (bothWaysProbDetail != \"\"))"
  direction                 = "BIDIRECTIONAL"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_agent_to_agent" "test" {
  test_name      = "User Acceptance Test - Aget To Agent"
  interval       = 120
  alerts_enabled = true

  direction       = "BIDIRECTIONAL"
  protocol        = "TCP"
  target_agent_id = "2334" #Frankfurt, Germany

  agents {
    agent_id = data.thousandeyes_agent.arg_amsterdam.agent_id
  }

  alert_rules {
    rule_id = 921617 #Agent-To-Agent Default Alert Rule
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.alert-rule-agent_to_agent.id
  }
}
