---
page_title: "Dashboard Widgets: Charts and Maps"
subcategory: ""
---

# Dashboard Widgets: Charts and Maps

This guide covers the provider-supported widgets used for trends, breakdowns, and geographic views:

- `Time Series: Line`
- `Time Series: Stacked Area`
- `Pie Chart`
- `Box and Whiskers`
- `Bar Chart: Grouped`
- `Bar Chart: Stacked`
- `Color Grid`
- `Map`

For the resource schema, see [`thousandeyes_dashboard`](../resources/dashboard.md). For dashboard-wide saved filters and timespans, see [Dashboard filters and timespans](dashboard-filters-and-timespans.md).

## Shared chart fields

Most widgets in this guide use the same core fields:

- `data_source`
- `metric_group`
- `metric`
- `measure`
- optional `fixed_timespan`

Several widgets also require a widget-specific config block:

- `timeseries_config`
- `stacked_area_config`
- `pie_chart_config`
- `box_and_whiskers_config`
- `grouped_bar_chart_config`
- `stacked_bar_chart_config`
- `color_grid_config`
- `geo_map_config`

## Time Series: Line

Use `Time Series: Line` for single-metric trends over time.

Example:

```terraform
resource "thousandeyes_dashboard" "timeseries" {
  title = "Time Series"

  widgets {
    type         = "Time Series: Line"
    title        = "Agent alerts over time"
    visual_mode  = "Full"
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

    timeseries_config {
      group_by                         = "AGENT"
      show_timeseries_overall_baseline = false
      is_timeseries_one_chart_per_line = false
    }
  }
}
```

Provider notes:

- Use `timeseries_config.group_by` to define the aggregation dimension.
- `show_timeseries_overall_baseline` and `is_timeseries_one_chart_per_line` are optional display controls.

## Time Series: Stacked Area

Use `Time Series: Stacked Area` when you want grouped values to accumulate over time.

Example:

```terraform
resource "thousandeyes_dashboard" "stacked_area" {
  title = "Cloud Events"

  widgets {
    type         = "Time Series: Stacked Area"
    title        = "Events by region"
    visual_mode  = "Full"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    stacked_area_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }
}
```

Provider notes:

- `stacked_area_config.group_by` is required for usable configurations.
- This widget is a good fit when the same metric needs to be split by region, account, or other aggregate property.

## Pie Chart

Use `Pie Chart` for proportional breakdowns of a single metric.

Example:

```terraform
resource "thousandeyes_dashboard" "pie" {
  title = "Pie Chart"

  widgets {
    type         = "Pie Chart"
    title        = "Events by region"
    visual_mode  = "Full"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    pie_chart_config {
      group_by = "CLOUD_NATIVE_MONITORING-REGION"
    }
  }
}
```

Provider notes:

- `pie_chart_config.group_by` is required for a useful pie chart configuration.
- Keep the chosen metric and grouping aligned so the chart represents a meaningful share-of-total view.

## Box and Whiskers

Use `Box and Whiskers` when distribution matters more than just the average.

Example:

```terraform
resource "thousandeyes_dashboard" "box_and_whiskers" {
  title = "Distribution"

  widgets {
    type         = "Box and Whiskers"
    title        = "Alerts by country"
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

    box_and_whiskers_config {
      group_by = "COUNTRY"
    }
  }
}
```

Provider notes:

- Use `box_and_whiskers_config.group_by` to choose the comparison dimension.
- Scale settings are optional and can be added when a fixed axis improves readability.

## Bar Chart: Grouped

Use `Bar Chart: Grouped` when you need side-by-side comparisons by two dimensions.

Example:

```terraform
resource "thousandeyes_dashboard" "grouped_bar" {
  title = "Grouped Bar Chart"

  widgets {
    type         = "Bar Chart: Grouped"
    title        = "Alerts by country and test"
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

    grouped_bar_chart_config {
      group_by                = "COUNTRY"
      axis_group_by           = "TEST"
      limit                   = 12
      show_labels             = true
      is_horizontal_bar_chart = false
    }
  }
}
```

Provider notes:

- `group_by` defines the grouped series and `axis_group_by` defines the axis split.
- Use `limit` to keep grouped charts readable.

## Bar Chart: Stacked

Use `Bar Chart: Stacked` when one axis category should show cumulative segments.

Example:

```terraform
resource "thousandeyes_dashboard" "stacked_bar" {
  title = "Stacked Bar Chart"

  widgets {
    type         = "Bar Chart: Stacked"
    title        = "Events by region"
    visual_mode  = "Full"
    data_source  = "CLOUD_NATIVE_MONITORING"
    metric_group = "CLOUD_NATIVE_MONITORING-EVENTS"
    metric       = "CLOUD_NATIVE_MONITORING-ALL_EVENTS"

    measure {
      type = "CLOUD_NATIVE_MONITORING-SUM"
    }

    fixed_timespan {
      value = 1
      unit  = "Days"
    }

    stacked_bar_chart_config {
      axis_group_by           = "CLOUD_NATIVE_MONITORING-REGION"
      limit                   = 8
      show_labels             = true
      is_horizontal_bar_chart = true
    }
  }
}
```

Provider notes:

- Use `axis_group_by` to control the bar segmentation dimension.
- `is_horizontal_bar_chart` switches between horizontal bars and vertical columns.

## Color Grid

Use `Color Grid` when you want threshold-like card coloring across many entities.

Example:

```terraform
resource "thousandeyes_dashboard" "color_grid" {
  title = "Color Grid"

  widgets {
    type         = "Color Grid"
    title        = "Alert heat map"
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

    color_grid_config {
      min_scale      = 0
      max_scale      = 100
      cards          = "COUNTRY"
      group_cards_by = "TEST"
      columns        = 2
      limit          = 6
    }
  }
}
```

Provider notes:

- `cards` defines what each card represents.
- `group_cards_by` controls the secondary grouping within the grid.

## Map

Use `Map` when the geographic dimension is the primary way you want to interpret the data.

Example:

```terraform
resource "thousandeyes_dashboard" "map" {
  title = "Map"

  widgets {
    type         = "Map"
    title        = "BGP alerts by country"
    visual_mode  = "Full"
    data_source  = "ALERTS"
    metric_group = "ALERTS"
    metric       = "ALERT_COUNT_BGP"

    measure {
      type = "MEAN"
    }

    geo_map_config {
      min_scale           = 0
      max_scale           = 100
      group_by            = "COUNTRY"
      is_geo_map_per_test = false
    }
  }
}
```

Provider notes:

- `geo_map_config.group_by` controls whether the map aggregates by country, continent, agent, or another supported geography-aware property.
- `is_geo_map_per_test` can be used to split the map by test when the data source supports that view.

## Choosing between chart widgets

- Use `Time Series: Line` for one metric over time.
- Use `Time Series: Stacked Area` when the grouped components should add up over time.
- Use `Pie Chart` or `Bar Chart` widgets for current-period composition rather than time progression.
- Use `Map` when geography matters more than raw ranking.

Use [Dashboard widgets: status and summary](dashboard-widgets-status-and-summary.md) for list, number, and table widgets.
