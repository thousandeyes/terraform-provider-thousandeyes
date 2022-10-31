resource "thousandeyes_agent_to_server" "example_agent_to_server_test" {
  test_name      = "Example Agent to Server test set from Terraform provider"
  interval       = 120
  alerts_enabled = false
  use_public_bgp = false

  server = "www.thousandeyes.com"
  port   = 443

  agents {
    agent_id = 3 # Singapore
  }
}
