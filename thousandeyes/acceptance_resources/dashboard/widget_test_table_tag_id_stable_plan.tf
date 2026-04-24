resource "thousandeyes_dashboard" "test_dashboard_test_table_tag_id_stable_plan" {
  description = "Test Dashboard with Test Table Tag ID stable plan coverage"
  title       = "Test Dashboard Test Table Tag ID Stable Plan"
  is_private  = false

  widgets {
    type        = "Test Table"
    title       = "Test Table Tag ID Stable Plan"
    visual_mode = "Full"

    test_table_config {
      filter {
        type = "all"

        filters {
          key   = "Tag ID"
          value = "123"
        }
      }
    }
  }
}
