data "thousandeyes_agent" "arg_frankfurt" {
  agent_name = "Frankfurt, Germany"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default Page Load Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Page Load Alert Rule"
  alert_type                = "page-load"
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
  agents         = [data.thousandeyes_agent.arg_frankfurt.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
