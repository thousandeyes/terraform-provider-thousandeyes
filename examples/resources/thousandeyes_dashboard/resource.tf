resource "thousandeyes_dashboard" "example" {
  description      = "Example Dashboard description"
  title            = "Example Dashboard"
  is_private       = true
  default_timespan {
    duration = 7200
  }
}
