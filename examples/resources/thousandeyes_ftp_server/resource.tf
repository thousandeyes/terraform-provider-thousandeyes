resource "thousandeyes_ftp_server" "example_ftp_server_test" {
  test_name      = "Example FTP Server test set from Terraform provider"
  interval       = 120
  alerts_enabled = false


  url = "www.thousandeyes.com"
  request_type = "Download"

  network_measurements = false #This flag is set to false to disable bandwidth measurements which is only possible when using Enterprise Agents

  username = "admin"
  password = "welcome"

  agents {
    agent_id = 3 # Singapur Cloud Agent
  }
}
