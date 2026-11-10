# TE_ENDPOINT_AGENT_TAG_ID must refer to a dynamic Endpoint Agent tag created
# with the CP-5555 `thousandeyes_tag` support.
resource "thousandeyes_endpoint_agent_to_server" "test" {
  test_name           = "Endpoint Agent-to-Server acceptance test updated"
  server              = "www.thousandeyes.com"
  interval            = 300
  protocol            = "icmp"
  agent_selector_type = "agent-tags"
  tag_ids             = ["__ENDPOINT_AGENT_TAG_ID__"]
}
