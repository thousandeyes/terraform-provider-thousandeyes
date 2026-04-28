resource "thousandeyes_dashboard" "test_dashboard_map_widget" {
  description = "Test Dashboard with Map Widget"
  title       = "Test Dashboard Map Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Map"
    title       = "Test Map Widget"
    visual_mode = "Full"
    metric_group= "ALERTS"
    metric      = "ALERT_COUNT_BGP"
    measure     {
      type = "MEAN"
    }
    data_source = "ALERTS"

    geo_map_config {
      group_by = "COUNTRY"
    }
  }
}
