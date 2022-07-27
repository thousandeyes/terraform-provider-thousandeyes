resource "thousandeyes_agent_to_agent" "example_agent_to_agent_test" {
  test_name      = "Example Agent to Agent test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  direction = "TO_TARGET"
  protocol = "TCP"
  target_agent_id = "5"
  throughput_rate = 500

  agents {
    agent_id = 3 # Singapur
  }
}
