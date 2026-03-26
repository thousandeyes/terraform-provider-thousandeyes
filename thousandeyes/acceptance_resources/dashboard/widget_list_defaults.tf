resource "thousandeyes_dashboard" "test_dashboard_list_defaults" {
  description = "Test Dashboard List Defaults"
  title       = "Test Dashboard List Defaults"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "List"
    title       = "List With Defaults"
    data_source = "EVENT_DETECTION"
    direction   = "TO_TARGET"
    measure {
      type = "MEAN"
    }
  }
}
