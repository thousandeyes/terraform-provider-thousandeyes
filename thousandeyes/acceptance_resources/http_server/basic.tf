data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "alert-rule-http-test" {
  rule_name                 = "Custom UAT HTTP Alert Rule"
  alert_type                = "HTTP Server"
  expression                = "(((responseCode >= 400) || (responseCode == 0)))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_http_server" "test" {
  test_name      = "User Acceptance Test - HTTP"
  interval       = 120
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"

  agents {
    agent_id = data.thousandeyes_agent.arg_amsterdam.agent_id
  }

  alert_rules {
    rule_id = 921621 #HTTP Default Alert Rule
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.alert-rule-http-test.id
  }
}
