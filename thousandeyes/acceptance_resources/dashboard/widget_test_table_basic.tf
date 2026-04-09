resource "thousandeyes_dashboard" "test_dashboard_test_table_widget" {
  description = "Test Dashboard with Test Table Widget"
  title       = "Test Dashboard Test Table Widget"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Test Table"
    title       = "Test Table Widget"
    visual_mode = "Full"
    data_source = "ALERTS"

    test_table_config {
      filter {
        type = "all"

        filters {
          key   = "Test Name"
          value = "API"
        }
      }

      exclude {
        type = "any"

        filters {
          key   = "Test ID"
          value = "123"
        }
      }
    }
  }
}
