resource "thousandeyes_voice" "example_voice_test" {
  test_name      = "Example RTP stream test set from Terraform provider"
  interval       = 120
  alerts_enabled = false

  target_agent_id = 4 # Tokyo

  agents {
    agent_id = 3 # Singapur
  }
}
