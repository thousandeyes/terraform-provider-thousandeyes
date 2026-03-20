resource "thousandeyes_dashboard" "test_dashboard_list_widget" {
  description = "Test Dashboard with List Widget"
  title       = "Test Dashboard List Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "List"
    title       = "Test List Widget"
    visual_mode = "Full"
    data_source = "EVENT_DETECTION"
    direction   = "TO_TARGET"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    list_config {
      active_within_value = 7
      active_within_unit = "Days"
    }
  }
}
