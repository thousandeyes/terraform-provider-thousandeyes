resource "thousandeyes_dashboard" "test_dashboard_color_grid_widget" {
  description = "Test Dashboard with Color Grid Widget"
  title       = "Test Dashboard Color Grid Widget"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Color Grid"
    title        = "Color Grid Widget"
    visual_mode  = "Full"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    color_grid_config {
      min_scale            = 0
      max_scale            = 100
      cards                = "COUNTRY"
      group_cards_by       = "TEST"
      columns              = 2
      limit                = 6
      sort_by              = "Value"
      sort_direction       = "Descending"
      sort_group_by        = "Alphabetical"
      sort_group_direction = "Ascending"
    }
  }
}
