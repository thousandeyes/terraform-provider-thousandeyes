resource "thousandeyes_ftp_server" "example_ftp_server_test" {
  test_name            = "Example FTP Server test set from Terraform provider"
  interval             = 120
  alerts_enabled       = false
  url                  = "www.thousandeyes.com"
  request_type         = "download"
  network_measurements = false #This flag is set to false to disable bandwidth measurements which is only possible when using Enterprise Agents
  username             = "admin"
  password             = "welcome" #Password won't be saved in state, so Terraform will notify that this field should be updated durring terraform plan
  agents               = ["3"]     # Singapore
}
