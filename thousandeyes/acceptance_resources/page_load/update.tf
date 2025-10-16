data "thousandeyes_agent" "arg_amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default Page Load Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT Page Load Alert Rule (Updated)"
  alert_type                = "page-load"
  expression                = "((errorCount >= 50) && (pageLoadHasError == true) && (Auto(timeToFirstByte >= Medium sensitivity)) && (Auto(onContentLoadTime >= Medium sensitivity)))"
  rounds_violating_out_of   = 3
  rounds_violating_required = 3
  minimum_sources           = 1
}

resource "thousandeyes_page_load" "test" {
  test_name      = "User Acceptance Test - Page Load (Updated)"
  interval       = 300
  http_interval  = 300
  alerts_enabled = true
  url            = "https://www.thousandeyes.com"
  agents         = [data.thousandeyes_agent.arg_amsterdam.agent_id]
  alert_rules    = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
