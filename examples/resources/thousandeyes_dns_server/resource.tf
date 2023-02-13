resource "thousandeyes_dns_server" "example_dns_server_test" {
  test_name      = "Example DNS server test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  domain = "www.thousandeyes.com"

  dns_servers {
    server_name = "ns1.google.com"
  }

  agents {
    agent_id = 3 # Singapore
  }
}
