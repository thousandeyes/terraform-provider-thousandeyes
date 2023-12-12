data "thousandeyes_agent" "test" {
  agent_name = "Vancouver, Canada"
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
  alerts_enabled       = true
  network_measurements = false
  url                  = "ftp://speedtest.tele2.net/"

  agents {
    agent_id = data.thousandeyes_agent.test.agent_id
  }
}
