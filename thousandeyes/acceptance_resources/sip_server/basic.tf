data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT SIP Server Alert Rule"
  alert_type                = "sip-server"
  expression                = "((sipErrorType != \"None\") && (Auto(connectTime >= Medium sensitivity)) && (Auto(sipResponseTime >= Medium sensitivity)))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_sip_server" "test" {
  test_name      = "User Acceptance Test - SIP Server"
  interval       = 120
  alerts_enabled = true
  probe_mode     = "sack"
  protocol      = "tcp"
  # target_sip_credentials {
  #   auth_user     = ""
  #   protocol      = "tcp"
  #   port          = 5060
  #   sip_registrar = "thousandeyes.com"
  # }
  agents     = [data.thousandeyes_agent.amsterdam.agent_id]
  alert_rules = [thousandeyes_alert_rule.test.id, "921618"]
}
