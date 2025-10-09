data "thousandeyes_agent" "test" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default FTP Alert Rule 2.0"
}

resource "thousandeyes_alert_rule" "test" {
  rule_name                 = "ftp server test"
  alert_type                = "ftp"
  expression                = "((ftpErrorType != \"None\"))"
  rounds_violating_out_of   = 1
  rounds_violating_required = 1
  minimum_sources_pct       = 80
}

resource "thousandeyes_ftp_server" "test" {
  password             = "test_password"
  username             = "test_username"
  test_name            = "New User Acceptance Test - FTP Server"
  description          = "description"
  request_type         = "download"
  ftp_time_limit       = 10
  ftp_target_time      = 1000
  interval             = 120
  alerts_enabled       = true
  network_measurements = false
  url                  = "ftp://speedtest.tele2.net/"
  agents               = [data.thousandeyes_agent.test.agent_id]
  alert_rules          = [thousandeyes_alert_rule.test.id, data.thousandeyes_alert_rule.def_alert_rule.id]
}
