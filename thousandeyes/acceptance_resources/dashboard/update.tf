resource "thousandeyes_dashboard" "test_dashboard" {
  description = "Updated Test Dashboard Description"
  title       = "Test Dashboard (Updated)"
  is_private  = true

  default_timespan {
    duration = 7200
  }
}
