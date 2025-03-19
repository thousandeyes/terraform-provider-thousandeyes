data "thousandeyes_agent" "amsterdam" {
  agent_name = "Amsterdam, Netherlands"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "Custom UAT SIP Server Alert Rule"
  alert_type                = "SIP Server"
  expression                = "((sipErrorType != \"None\") && (Auto(connectTime >= Medium sensitivity)) && (Auto(sipResponseTime >= Medium sensitivity)))"
  minimum_sources           = 1
  rounds_violating_required = 1
  rounds_violating_out_of   = 1
}

resource "thousandeyes_sip_server" "test" {
  test_name      = "User Acceptance Test - SIP Server"
  interval       = 120
  alerts_enabled = true
  probe_mode     = "SACK"
  target_sip_credentials {
    auth_user     = ""
    protocol      = "TCP"
    port          = 5060
    sip_registrar = "thousandeyes.com"
  }

  agents {
    agent_id = data.thousandeyes_agent.amsterdam.agent_id
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.test.id
  }

  alert_rules {
    rule_id = 921618 #SIP Server Default Alert Rule
  }
}
