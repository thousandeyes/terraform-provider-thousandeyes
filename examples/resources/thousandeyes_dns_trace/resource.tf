resource "thousandeyes_dns_trace" "example_dns_trace_test" {
  test_name      = "Example DNS Trace test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  domain = "www.thousandeyes.com ANY"

  agents {
    agent_id = 3 # Singapur
  }
}
