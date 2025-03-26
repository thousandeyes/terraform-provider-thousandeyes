data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "alert-rule-http-test" {
  rule_name                 = "Custom UAT Page Load Alert Rule"
  alert_type                = "Page Load"
  expression                = "((errorCount >= 50) && (pageLoadHasError == true) && (Auto(timeToFirstByte >= Medium sensitivity)) && (Auto(onContentLoadTime >= Medium sensitivity)))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources           = 1
}

resource "thousandeyes_page_load" "test" {
  test_name      = "User Acceptance Test - Page Load"
  interval       = 120
  http_interval  = 120
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"

  agents {
    agent_id = data.thousandeyes_agent.arg_amsterdam.agent_id
  }

  alert_rules {
    rule_id = 921620 #Page Load Default Alert Rule
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.alert-rule-http-test.id
  }
}
