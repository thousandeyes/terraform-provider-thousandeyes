// Minimal dashboard with no default_timespan block — API still returns a default timespan on read.
resource "thousandeyes_dashboard" "test_omit_default_timespan" {
  description        = "Omit default_timespan"
  title              = "Test Dashboard Omit Default Timespan"
  is_private         = false
  is_global_override = false
}
