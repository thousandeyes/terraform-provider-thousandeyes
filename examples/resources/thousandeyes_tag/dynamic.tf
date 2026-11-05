resource "thousandeyes_tag" "endpoint_agents" {
  key         = "office"
  value       = "london-windows"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = "and"

  filters {
    key    = "platform"
    values = ["Windows"]
    mode   = "in"
    # scope defaults to "default" for built-in ThousandEyes agent fields.
  }

  filters {
    key    = "hostname"
    values = ["LON-*"]
    mode   = "in"
  }
}
