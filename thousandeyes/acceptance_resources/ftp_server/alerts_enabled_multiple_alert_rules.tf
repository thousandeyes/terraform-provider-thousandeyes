data "thousandeyes_agent" "test" {
  agent_name = "Vancouver, Canada"
}

resource "thousandeyes_alert_rule" "alert-rule-ftp-test" {
  rule_name                 = "ftp server test"
  alert_type                = "FTP"
  expression                = "((ftpErrorType != \"None\"))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources_pct       = 80
}

resource "thousandeyes_ftp_server" "test" {
  password             = "test_password"
  username             = "test_username"
  test_name            = "Acceptance Test - FTP Multiple Alert Rules"
  description          = "description"
  request_type         = "Download"
  ftp_time_limit       = 10
  ftp_target_time      = 1000
  interval             = 900
  alerts_enabled       = true
  network_measurements = false
  url                  = "ftp://speedtest.tele2.net/"

  agents {
    agent_id = data.thousandeyes_agent.test.agent_id
  }

  alert_rules {
    rule_id = 921623 #FTP Default Alert Rule
  }

  alert_rules {
    rule_id = thousandeyes_alert_rule.alert-rule-ftp-test.id
  }
}
