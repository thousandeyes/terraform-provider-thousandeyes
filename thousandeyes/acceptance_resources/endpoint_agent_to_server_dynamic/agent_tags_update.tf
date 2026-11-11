resource "thousandeyes_tag" "endpoint" {
  key         = "User Acceptance Test - Endpoint Agent To Server Dynamic"
  value       = "selected"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = "and"

  filters {
    key    = "platform"
    mode   = "in"
    values = ["Windows"]
  }
}

resource "thousandeyes_endpoint_agent_to_server_dynamic" "test" {
  test_name           = "User Acceptance Test - Endpoint Agent To Server Dynamic (Updated)"
  application         = "zoom"
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.endpoint.id]
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 120
  is_enabled          = true
  protocol            = "tcp"
}
