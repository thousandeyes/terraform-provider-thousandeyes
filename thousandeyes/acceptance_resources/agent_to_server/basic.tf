data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Agent To Server Alert Rule"
  alert_type                = "End-to-End (Server)"
  expression                = "((loss >= 10%) || (probDetail != \"\"))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_agent_to_server" "test" {
  test_name      = "User Acceptance Test - Agent To Server"
  interval       = 120
  alerts_enabled = true
  server         = "api.stg.thousandeyes.com"
  protocol       = "TCP"
  port           = 443
  probe_mode     = "SACK"

  agents {
    agent_id = data.thousandeyes_agent.amsterdam.agent_id
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.test.id
  }

  alert_rules {
    rule_id = 921610 #Agent-To-Server Default Alert Rule
  }
}
