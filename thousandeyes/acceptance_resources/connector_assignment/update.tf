resource "thousandeyes_connector_assignment" "test" {
  connector_id = local.connector_id
  operation_ids = [
    local.operation_id_1,
    local.operation_id_2
  ]
}
