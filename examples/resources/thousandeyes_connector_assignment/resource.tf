resource "thousandeyes_connector_assignment" "example" {
  connector_id = "connector_id"
  operation_ids = [
    "webhook_operation_id_1",
    "webhook_operation_id_2"
  ]
}
