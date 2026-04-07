resource "thousandeyes_dashboard" "test_dashboard_map_widget" {
  description = "Test Dashboard with Map Widget (Updated)"
  title       = "Test Dashboard Map Widget (Updated)"
  is_private  = false
  default_timespan {
    duration = 7200
  }

  widgets {
    type        = "Map"
    title       = "Test Map Widget (Updated)"
    visual_mode = "Full"
    metric_group= "ALERTS"
    metric      = "ALERT_COUNT_BGP"
    measure     {
      type = "MEAN"
    }
    data_source = "ALERTS"

    geo_map_config {
      min_scale = 10
      max_scale = 200
      group_by = "CONTINENT"
      is_geo_map_per_test = true
    }
  }
}
