---
page_title: "Dashboard Filters and Timespans"
subcategory: ""
---

# Dashboard Filters and Timespans

Use this guide with [`thousandeyes_dashboard`](../resources/dashboard.md) when you need to:

- attach a saved dashboard filter with `global_filter_id`
- reuse an existing dashboard filter by name with [`thousandeyes_dashboard_filter`](../data-sources/dashboard_filter.md)
- understand when to use dashboard-level `default_timespan` versus widget-level `fixed_timespan`
- configure widget-level `filter` blocks correctly

For the product concepts behind dashboard filters, see the ThousandEyes documentation for [Dashboards](https://docs.thousandeyes.com/product-documentation/dashboards), [Dashboard Widgets](https://docs.thousandeyes.com/product-documentation/dashboards/dashboard-widgets), and [Dashboard Filters](https://docs.thousandeyes.com/product-documentation/dashboards/dashboard-filters).

## Saved dashboard filters with `global_filter_id`

`global_filter_id` sets the dashboard's saved default filter. In the ThousandEyes product, these are reusable dashboard-wide filter sets. The provider does not currently manage dashboard filters as a resource, so the Terraform workflow is:

1. Create or save the filter in the ThousandEyes UI or API.
2. Look it up by exact name with `data "thousandeyes_dashboard_filter"`.
3. Pass the returned `filter_id` to `thousandeyes_dashboard.global_filter_id`.

```terraform
data "thousandeyes_dashboard_filter" "operations_core_services" {
  name = "Operations - Core Services"
}

resource "thousandeyes_dashboard" "operations_overview" {
  title              = "Operations Overview"
  description        = "Dashboard generated from Terraform using an existing dashboard filter"
  global_filter_id   = data.thousandeyes_dashboard_filter.operations_core_services.filter_id
  is_global_override = true

  default_timespan {
    duration = 7200
  }
}
```

Use `global_filter_id` for a saved dashboard-wide starting view. Use widget `filter` blocks when you need a widget or number card to stay scoped to specific IDs in configuration.

## Dashboard filters versus widget filters

These are different mechanisms:

- `global_filter_id` references a saved dashboard filter set.
- `widgets.filter` stores widget-specific filter criteria directly in Terraform.
- `widgets.number_cards.filter` stores card-specific filter criteria for `Number` widgets.

The API models widget and number-card filters as a property-to-ID map. The provider exposes that map as repeated blocks with a `property` and a set of string `values`.

```terraform
data "thousandeyes_agent" "singapore" {
  agent_name = "Singapore"
}

resource "thousandeyes_http_server" "checkout" {
  test_name = "Checkout API"
  interval  = 300
  url       = "https://example.com/health"
  agents    = [data.thousandeyes_agent.singapore.agent_id]
}

resource "thousandeyes_dashboard" "service_health" {
  title = "Service Health"

  widgets {
    type         = "Time Series: Line"
    title        = "Checkout availability"
    data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
    metric_group = "HTTP_SERVER"
    metric       = "WEB_AVAILABILITY"

    measure {
      type = "MEAN"
    }

    filter {
      property = "TEST"
      values   = [thousandeyes_http_server.checkout.test_id]
    }

    filter {
      property = "AGENT"
      values   = [data.thousandeyes_agent.singapore.agent_id]
    }

    timeseries_config {
      group_by = "TEST"
    }
  }
}
```

### Filter value semantics

- Define `values` as strings in Terraform, even when the upstream API documentation shows numeric arrays for some filter properties.
- Prefer Terraform references when the provider already exposes the ID you need, such as `test_id` or `agent_id`.
- `filter` and `values` are sets, so order is not significant.

### Product limitations that matter in Terraform

The ThousandEyes product documentation calls out a few behaviors worth keeping in mind:

- Dashboard filters are scoped per data source. A saved filter for one data source does not affect widgets using another data source.
- Some widget types are not currently affected by dashboard filters: `Agent Status`, `Alert List`, and `Test Table`.
- The UI has a "Lock Widget Filters" capability, but that is not a Terraform attribute on `thousandeyes_dashboard`. Do not model it as if it were provider-managed behavior.

## Dashboard `default_timespan` versus widget `fixed_timespan`

Use `default_timespan` for the dashboard's default time range. Use `fixed_timespan` only when a specific widget, number card, or similar element should keep its own time range.

```terraform
resource "thousandeyes_dashboard" "latency_overview" {
  title              = "Latency Overview"
  is_global_override = false

  default_timespan {
    duration = 7200
  }

  widgets {
    type         = "Time Series: Line"
    title        = "7-day latency trend"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "TOTAL"
    }

    fixed_timespan {
      value = 7
      unit  = "Days"
    }

    timeseries_config {
      group_by = "AGENT"
    }
  }
}
```

### How `is_global_override` works

- `is_global_override = true`: the dashboard's `default_timespan` overrides widget-level timespans on compatible widgets.
- `is_global_override = false`: compatible widgets can continue using their own `fixed_timespan`.

### Important provider behavior

`default_timespan` is optional and computed. If you omit it, Terraform stores the API's default value in state to avoid perpetual drift.

If you explicitly set `default_timespan` and later remove the block from configuration, the provider typically shows no change and the previous timespan stays in place. To change the value again, set a new explicit `default_timespan`.

## Card filters and column direction

Use `number_cards.filter` when cards inside the same `Number` widget need different scopes.

Card example:

```terraform
resource "thousandeyes_dashboard" "regional_summary" {
  title = "Regional Summary"

  widgets {
    type  = "Number"
    title = "Regional KPIs"

    number_cards {
      description  = "EMEA availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }

      filter {
        property = "AGENT_LABEL"
        values   = ["emea"]
      }
    }

    number_cards {
      description  = "AMER availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }

      filter {
        property = "AGENT_LABEL"
        values   = ["amer"]
      }
    }
  }
}
```

`Multi Metric Table` columns do not support per-column `filter` blocks. Each column does support its own metric fields, measure, and, for certain data sources, `direction`.

```terraform
resource "thousandeyes_dashboard" "directional_summary" {
  title = "Directional Summary"

  widgets {
    type  = "Multi Metric Table"
    title = "Directional latency summary"

    multi_metric_table_config {
      row_group_by = "COUNTRY"
      limit        = 10
    }

    multi_metric_columns {
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "AGENT_TO_AGENT"
      direction    = "TO_TARGET"
      metric       = "NET_LATENCY"

      measure {
        type = "MEAN"
      }
    }

    multi_metric_columns {
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "AGENT_TO_AGENT"
      direction    = "FROM_TARGET"
      metric       = "NET_LATENCY"

      measure {
        type = "MEAN"
      }
    }
  }
}
```

Use the widget guides for complete examples of each supported widget type:

- [Dashboard widgets: status and summary](dashboard-widgets-status-and-summary.md)
- [Dashboard widgets: charts and maps](dashboard-widgets-charts-and-maps.md)
