resource "thousandeyes_tag" "endpoint" {
  key         = "User Acceptance Test - Endpoint HTTP Server"
  value       = "selected"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = "and"

  filters {
    key    = "platform"
    mode   = "in"
    values = ["Mac"]
  }
}

resource "thousandeyes_endpoint_http_server" "test" {
  test_name           = "User Acceptance Test - Endpoint HTTP Server"
  url                 = "https://example.com"
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.endpoint.id]
  max_machines        = 25
  ip_version          = "V4_ONLY"
  interval            = 60
  is_enabled          = false
  protocol            = "icmp"
}
