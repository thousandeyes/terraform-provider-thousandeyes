data "thousandeyes_agent" "test" {
  agent_name = "Amsterdam, Netherlands"
}

data "thousandeyes_alert_rule" "def_alert_rule" {
  rule_name = "Default FTP Alert Rule 2.0"
}

resource "thousandeyes_ftp_server" "test" {
  password             = "test_password"
  username             = "test_username"
  test_name            = "Acceptance Test - FTP Alerts Enabled"
  description          = "description"
  request_type         = "download"
  ftp_time_limit       = 10
  ftp_target_time      = 1000
  interval             = 900
  alerts_enabled       = true
  network_measurements = false
  bgp_measurements     = false
  url                  = "ftp://speedtest.tele2.net/"
  agents               = [data.thousandeyes_agent.test.agent_id]
  alert_rules          = [data.thousandeyes_alert_rule.def_alert_rule]
}
