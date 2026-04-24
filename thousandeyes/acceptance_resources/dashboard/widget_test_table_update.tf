resource "thousandeyes_dashboard" "test_dashboard_test_table_widget" {
  description = "Test Dashboard with Test Table Widget (Updated)"
  title       = "Test Dashboard Test Table Widget (Updated)"
  is_private  = false

  default_timespan {
    duration = 3600
  }

  widgets {
    type        = "Test Table"
    title       = "Test Table Widget (Updated)"
    visual_mode = "Full"

    test_table_config {
      filter {
        type = "any"

        filters {
          key   = "Target"
          value = "example.com"
        }
      }

      exclude {
        type = "all"

        filters {
          key   = "Tag ID"
          value = "456"
        }
      }
    }
  }
}
