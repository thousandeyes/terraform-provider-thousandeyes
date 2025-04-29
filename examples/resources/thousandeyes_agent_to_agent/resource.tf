resource "thousandeyes_agent_to_agent" "example_agent_to_agent_test" {
  test_name       = "Example Agent to Agent test set from Terraform provider"
  interval        = 120
  alerts_enabled  = false
  direction       = "to-target"
  protocol        = "tcp"
  target_agent_id = "5"
  agents          = ["3"] # Singapore
}
