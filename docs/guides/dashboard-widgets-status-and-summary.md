---
page_title: "Dashboard Widgets: Status and Summary"
subcategory: ""
---

# Dashboard Widgets: Status and Summary

This guide covers the provider-supported dashboard widgets that focus on current status, lists, tabular summaries, and number cards:

- `Agent Status`
- `Alert List`
- `List`
- `Number`
- `Table`
- `Multi Metric Table`
- `Test Table`

For the resource schema, see [`thousandeyes_dashboard`](../resources/dashboard.md). For dashboard-wide filters and timespan behavior, see [Dashboard filters and timespans](dashboard-filters-and-timespans.md).

## Shared patterns

Most summary widgets use some combination of:

- `data_source`
- `metric_group`
- `metric`
- `measure`
- a widget-specific config block such as `list_config` or `table_config`

The exceptions in this guide are:

- `Number`, which defines metrics inside each `number_cards` block
- `Test Table`, which uses `test_table_config` instead of `data_source`, `metric_group`, and `metric`
- `Agent Status`, which uses `agent_status_config`

## Agent Status

Use `Agent Status` for a live summary of enterprise or endpoint agent health.

Minimum shape:

- `type`
- `title`
- `data_source`

Example:

```terraform
resource "thousandeyes_dashboard" "agent_status" {
  title = "Agent Status"

  widgets {
    type        = "Agent Status"
    title       = "Enterprise Agents"
    visual_mode = "Full"
    data_source = "CLOUD_AND_ENTERPRISE_AGENTS"

    agent_status_config {
      show       = "Owned Agents"
      agent_type = "Enterprise Agents"
    }
  }
}
```

This widget does not use `metric_group`, `metric`, or `measure`. Product documentation notes that dashboard-wide saved filters do not currently affect `Agent Status` widgets.

## Alert List

Use `Alert List` when you want a live list of active or recently cleared alerts.

Minimum shape:

- `type`
- `title`
- `data_source`

Example:

```terraform
resource "thousandeyes_dashboard" "alert_list" {
  title = "Alerts"

  widgets {
    type        = "Alert List"
    title       = "Critical alerts"
    visual_mode = "Full"
    data_source = "ALERTS"

    alert_list_config {
      alert_types = ["API", "DNS Server"]
      limit_to    = 15

      active_within_value = 7
      active_within_unit  = "Days"
    }
  }
}
```

Use `alert_list_config` for `alert_types`, `limit_to`, and the active-within window. Product documentation notes that dashboard-wide saved filters do not currently affect `Alert List` widgets.

## List

Use `List` for event-driven list views such as Event Detection.

Minimum shape:

- `type`
- `title`
- `data_source`
- `measure`

Example:

```terraform
resource "thousandeyes_dashboard" "event_list" {
  title = "Event Detection"

  widgets {
    type        = "List"
    title       = "Recent events"
    visual_mode = "Full"
    data_source = "EVENT_DETECTION"

    measure {
      type = "MEAN"
    }

    list_config {
      active_within_value = 7
      active_within_unit  = "Days"
    }
  }
}
```

`list_config` controls the active-within window shown by the widget. When the API returns default widget settings, the provider stores them in state even if the block is omitted from configuration.

## Number

Use `Number` for KPI-style cards, each with its own metric and optional scope.

Minimum shape:

- `type`
- `title`
- at least one `number_cards` block

Example:

```terraform
resource "thousandeyes_dashboard" "kpis" {
  title = "Key Metrics"

  widgets {
    type        = "Number"
    title       = "Service KPIs"
    visual_mode = "Full"

    number_cards {
      description  = "HTTP availability"
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_AVAILABILITY"

      measure {
        type = "MEAN"
      }
    }

    number_cards {
      description  = "Agent alerts"
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT_AGENT"

      measure {
        type = "TOTAL"
      }

      fixed_timespan {
        value = 1
        unit  = "Days"
      }
    }
  }
}
```

Set `data_source` on each `number_cards` block, not on the parent `widgets` block. Use `number_cards.filter` when different cards need different IDs or labels.

## Table

Use `Table` when one metric should be broken down across rows and columns.

Minimum shape:

- `type`
- `title`
- `data_source`
- `metric_group`
- `metric`
- `measure`

Example:

```terraform
resource "thousandeyes_dashboard" "table_summary" {
  title = "Alert Summary"

  widgets {
    type         = "Table"
    title        = "Alerts by agent and test"
    visual_mode  = "Full"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_AGENT"

    measure {
      type = "MEAN"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    table_config {
      compare_to_previous_value = true
      row_group_by              = "AGENT"
      column_group_by           = "TEST"
      limit                     = 10
    }
  }
}
```

Use `table_config` to control grouping and the row limit. `row_group_by` and `column_group_by` follow the dashboard API aggregate-property values.

## Multi Metric Table

Use `Multi Metric Table` when one table should show multiple metrics side by side.

Minimum shape:

- `type`
- `title`
- `multi_metric_table_config`
- at least one `multi_metric_columns` block

Example:

```terraform
resource "thousandeyes_dashboard" "multi_metric" {
  title = "Multi-Metric Summary"

  widgets {
    type        = "Multi Metric Table"
    title       = "Availability and fetch time"
    visual_mode = "Full"

    multi_metric_table_config {
      compare_to_previous_value = true
      row_group_by              = "COUNTRY"
      limit                     = 10
    }

    multi_metric_columns {
      data_source  = "ALERTS"
      metric_group = "ALERTS"
      metric       = "ALERT_COUNT"

      measure {
        type = "MEAN"
      }
    }

    multi_metric_columns {
      data_source  = "CLOUD_AND_ENTERPRISE_AGENTS"
      metric_group = "HTTP_SERVER"
      metric       = "WEB_FETCH"

      measure {
        type = "MEAN"
      }
    }
  }
}
```

Each `multi_metric_columns` block has its own `data_source`, `metric_group`, `metric`, and `measure`. Some data sources also support column-level `direction`, such as `TO_TARGET` and `FROM_TARGET` for directional metrics.

## Test Table

Use `Test Table` when you want to include or exclude tests using textual criteria such as test name, target, or test ID.

Minimum shape:

- `type`
- `title`
- `test_table_config`

Example:

```terraform
resource "thousandeyes_dashboard" "test_table" {
  title = "Tests"

  widgets {
    type        = "Test Table"
    title       = "API tests"
    visual_mode = "Full"

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
```

`Test Table` does not use the generic widget `filter` block. It uses `test_table_config.filter` and `test_table_config.exclude`, each containing `filters { key, value }` items. Product documentation notes that dashboard-wide saved filters do not currently affect `Test Table` widgets.

## Next step

Use [Dashboard widgets: charts and maps](dashboard-widgets-charts-and-maps.md) for time series, breakdown charts, and map widgets.
