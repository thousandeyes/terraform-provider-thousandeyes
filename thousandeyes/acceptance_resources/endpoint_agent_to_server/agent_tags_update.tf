resource "thousandeyes_tag" "endpoint_agents" {
  key         = "User Acceptance Test - Endpoint Agent To Server"
  value       = "selected"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = "and"

  filters {
    key    = "platform"
    values = ["Windows"]
    mode   = "in"
  }
}

resource "thousandeyes_endpoint_agent_to_server" "test" {
  test_name           = "Endpoint Agent-to-Server acceptance test updated"
  server              = "www.thousandeyes.com"
  interval            = 300
  protocol            = "icmp"
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.endpoint_agents.id]
}
