resource "thousandeyes_tag" "example_tag" {
  key         = "Example Tag Key"
  value       = "Example Tag Value"
  object_type = "test"
  color       = "#b3de69"
  access_type = "all"
  icon        = "LABEL"
}