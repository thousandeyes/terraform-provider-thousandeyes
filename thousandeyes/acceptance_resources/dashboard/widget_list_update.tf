resource "thousandeyes_dashboard" "test_dashboard_list_widget" {
  description = "Test Dashboard with List Widget (Updated)"
  title       = "Test Dashboard List Widget (Updated)"
  is_private  = true
  default_timespan {
    duration = 7200
  }
  widgets {
    type        = "List"
    title       = "Test List Widget (Updated)"
    visual_mode = "Full"
    data_source = "ALERTS"
    direction   = "TO_TARGET"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    list_config {
      active_within_value = 14
      active_within_unit = "Days"
    }
  }
}
