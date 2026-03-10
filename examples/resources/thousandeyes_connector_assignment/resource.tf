resource "thousandeyes_connector_assignment" "example" {
  connector_id = "<existing_connector_id>"
  operation_ids = [
    "<existing_webhook_operation_id_1>",
    "<existing_webhook_operation_id_2>"
  ]
}
