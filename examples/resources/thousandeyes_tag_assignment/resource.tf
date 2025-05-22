resource "thousandeyes_tag_assignment" "example_assignment" {
  tag_id = "example_tag_resource_id"
  assignments {
    id   = "123456" #Id of existing entity (Test, Dashboard, etc.)
    type = "test"
  }
  assignments {
    id   = "123457"
    type = "test"
  }
}
