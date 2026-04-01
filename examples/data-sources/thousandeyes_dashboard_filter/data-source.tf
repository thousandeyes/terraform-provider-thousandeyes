data "thousandeyes_dashboard_filter" "operations_core_services" {
  name = "Operations - Core Services"
}

resource "thousandeyes_dashboard" "operations_overview" {
  title              = "Operations Overview"
  description        = "Dashboard generated from Terraform using an existing dashboard filter"
  is_private         = false
  global_filter_id   = data.thousandeyes_dashboard_filter.operations_core_services.filter_id
  is_global_override = true

  default_timespan {
    duration = 7200
  }
}
