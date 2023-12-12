data "thousandeyes_agent" "test" {
  agent_name = "Vancouver, Canada"
}

resource "thousandeyes_alert_rule" "alert-rule-ftp-test" {
  rule_name                 = "API Team: (${var.environment}) Alert Slack"
  alert_type                = "HTTP Server"
  expression                = "((probDetail != \"\"))"
  minimum_sources_pct       = 80
  notify_on_clear           = true
  rounds_violating_out_of   = 2
  rounds_violating_required = 2
  rounds_violating_mode     = "ANY"
}

resource "thousandeyes_ftp_server" "test" {
  password             = "test_password"
  username             = "test_username"
  test_name            = "Acceptance Test - FTP"
  description          = "description"
  request_type         = "Download"
  ftp_time_limit       = 10
  ftp_target_time      = 1000
  interval             = 900
  alerts_enabled       = false
  network_measurements = false
  url                  = "ftp://speedtest.tele2.net/"

  agents {
    agent_id = data.thousandeyes_agent.test.agent_id
  }
}
