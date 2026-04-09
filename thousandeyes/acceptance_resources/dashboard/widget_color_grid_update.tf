resource "thousandeyes_dashboard" "test_dashboard_color_grid_widget" {
  description = "Test Dashboard with Color Grid Widget (Updated)"
  title       = "Test Dashboard Color Grid Widget (Updated)"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type         = "Color Grid"
    title        = "Color Grid Widget (Updated)"
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
      min_scale            = 10
      max_scale            = 200
      cards                = "CONTINENT"
      group_cards_by       = "AGENT"
      columns              = 1
      limit                = 4
      sort_by              = "Alphabetical"
      sort_direction       = "Ascending"
      sort_group_by        = "Value"
      sort_group_direction = "Descending"
    }
  }
}
