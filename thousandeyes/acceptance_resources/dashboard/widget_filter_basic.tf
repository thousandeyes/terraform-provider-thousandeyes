resource "thousandeyes_dashboard" "test_dashboard_filter_widget" {
  description = "Test Dashboard with Widget Filter"
  title       = "Test Dashboard Filter Widget"
  is_private  = false
  default_timespan {
    duration = 3600
  }
  widgets {
    type        = "Time Series: Line"
    title       = "Filter test"
    visual_mode = "Full"
    data_source = "INTERNET_INSIGHTS"
    metric_group = "NETWORK_OUTAGES"
    metric       = "NETWORK_OUTAGES_LOCATIONS"

    measure {
      type = "MAXIMUM"
    }

    fixed_timespan {
      value = 1
      unit = "Days"
    }

    filter {
      property = "INSIGHTS_NETWORK"
      values   = [32133, 262287, 46606]
    }

    timeseries_config {
      group_by = "ALL"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }
}
