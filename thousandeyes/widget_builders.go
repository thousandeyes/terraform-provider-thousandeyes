package thousandeyes

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// buildGeoMapWidget builds a GeoMap widget from Terraform data
func buildGeoMapWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiGeoMapWidget("Map")
	setCommonBuilderFields(widget, data)

	// Set data_source (GeoMap-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.GeoMapDatasource(dataSource))
	}

	if configList := getListValue(data, "geo_map_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setFloat32FromMapIfPresent(config, "min_scale", widget.SetMinScale)
		setFloat32FromMapIfPresent(config, "max_scale", widget.SetMaxScale)
		if v := getStringValue(config, "unit"); v != "" {
			widget.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v, ok := config["is_geo_map_per_test"].(bool); ok {
			widget.SetIsGeoMapPerTest(v)
		}
	}

	return dashboards.ApiGeoMapWidgetAsApiWidget(widget)
}

// buildListWidget builds a List widget from Terraform data
func buildListWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiListWidget("List")
	setCommonBuilderFields(widget, data)

	// Set data_source (List-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.ListDatasource(dataSource))
	}

	if configList := getListValue(data, "list_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		// Handle active_within
		activeWithinValue := 0
		activeWithinUnit := ""
		if v, ok := config["active_within_value"].(int); ok && v != 0 {
			activeWithinValue = v
		}
		if v := getStringValue(config, "active_within_unit"); v != "" {
			activeWithinUnit = v
		}
		if activeWithinValue != 0 || activeWithinUnit != "" {
			activeWithin := dashboards.NewActiveWithin()
			if activeWithinValue != 0 {
				activeWithin.SetValue(int32(activeWithinValue))
			}
			if activeWithinUnit != "" {
				activeWithin.SetUnit(dashboards.LegacyDurationUnit(activeWithinUnit))
			}
			widget.SetActiveWithin(*activeWithin)
		}
	}

	return dashboards.ApiListWidgetAsApiWidget(widget)
}

// buildAlertListWidget builds an Alert List widget from Terraform data
func buildAlertListWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiAlertListWidget("Alert List")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.AlertListDatasource(dataSource))
	}

	if configList := getListValue(data, "alert_list_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setAlertTypesFromConfig(config, "alert_types", widget.SetAlertTypes)
		if rawLimit, exists := config["limit_to"]; exists {
			if limit, ok := rawLimit.(int); ok {
				widget.SetLimitTo(int32(limit))
			}
		}
		if activeWithin := buildActiveWithinFromConfig(config); activeWithin != nil {
			widget.SetActiveWithin(*activeWithin)
		}
	}

	return dashboards.ApiAlertListWidgetAsApiWidget(widget)
}

// buildBoxAndWhiskersWidget builds a Box and Whiskers widget from Terraform data
func buildBoxAndWhiskersWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
	setCommonBuilderFields(widget, data)

	// Set data_source (BoxAndWhiskers-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.BoxAndWhiskersDatasource(dataSource))
	}

	if configList := getListValue(data, "box_and_whiskers_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setFloat32FromMapIfPresent(config, "min_scale", widget.SetMinScale)
		setFloat32FromMapIfPresent(config, "max_scale", widget.SetMaxScale)
		if v := getStringValue(config, "unit"); v != "" {
			widget.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
	}

	return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(widget)
}

// buildTableWidget builds a Table widget from Terraform data
func buildTableWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiTableWidget("Table")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.TableDatasource(dataSource))
	}

	if configList := getListValue(data, "table_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v, ok := config["compare_to_previous_value"].(bool); ok {
			widget.SetCompareToPreviousValue(v)
		}
		if v := getStringValue(config, "row_group_by"); v != "" {
			widget.SetRowGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getStringValue(config, "column_group_by"); v != "" {
			widget.SetColumnGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getIntValue(config, "limit"); v != 0 {
			widget.SetLimit(int32(v))
		}
	}

	return dashboards.ApiTableWidgetAsApiWidget(widget)
}

// buildTestTableWidget builds a Test Table widget from Terraform data
func buildTestTableWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiTestTableWidget("Test Table")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.TestTableDatasource(dataSource))
	}

	if configList := getListValue(data, "test_table_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if filter := buildTestTableFilterConfig(config, "filter"); filter != nil {
			widget.SetFilter(*filter)
		}
		if exclude := buildTestTableFilterConfig(config, "exclude"); exclude != nil {
			widget.SetExclude(*exclude)
		}
	}

	return dashboards.ApiTestTableWidgetAsApiWidget(widget)
}

// buildPieChartWidget builds a Pie Chart widget from Terraform data
func buildPieChartWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiPieChartWidget("Pie Chart")
	setCommonBuilderFields(widget, data)

	// Set data_source (PieChart-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.PieChartDatasource(dataSource))
	}

	if configList := getListValue(data, "pie_chart_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
	}

	return dashboards.ApiPieChartWidgetAsApiWidget(widget)
}

// buildStackedBarChartWidget builds a Stacked Bar Chart widget from Terraform data
func buildStackedBarChartWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiStackedBarchartWidget("Bar Chart: Stacked")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.StackedBarChartDatasource(dataSource))
	}

	if configList := getListValue(data, "stacked_bar_chart_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v := getStringValue(config, "axis_group_by"); v != "" {
			widget.SetAxisGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getIntValue(config, "limit"); v != 0 {
			widget.SetLimit(int32(v))
		}
		if v, ok := config["show_labels"].(bool); ok {
			widget.SetShowLabels(v)
		}
		if v, ok := config["is_horizontal_bar_chart"].(bool); ok {
			widget.SetIsHorizontalBarChart(v)
		}
	}

	return dashboards.ApiStackedBarchartWidgetAsApiWidget(widget)
}

// buildGroupedBarChartWidget builds a Grouped Bar Chart widget from Terraform data
func buildGroupedBarChartWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiGroupedBarchartWidget("Bar Chart: Grouped")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.GroupedBarChartDatasource(dataSource))
	}

	if configList := getListValue(data, "grouped_bar_chart_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getStringValue(config, "axis_group_by"); v != "" {
			widget.SetAxisGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getIntValue(config, "limit"); v != 0 {
			widget.SetLimit(int32(v))
		}
		if v, ok := config["show_labels"].(bool); ok {
			widget.SetShowLabels(v)
		}
		if v, ok := config["is_horizontal_bar_chart"].(bool); ok {
			widget.SetIsHorizontalBarChart(v)
		}
	}

	return dashboards.ApiGroupedBarchartWidgetAsApiWidget(widget)
}

// buildStackedAreaWidget builds a Stacked Area Chart widget from Terraform data
func buildStackedAreaWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
	setCommonBuilderFields(widget, data)

	// Set data_source (StackedAreaChart-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.StackedAreaChartDatasource(dataSource))
	}

	if configList := getListValue(data, "stacked_area_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setFloat32FromMapIfPresent(config, "min_scale", widget.SetMinScale)
		setFloat32FromMapIfPresent(config, "max_scale", widget.SetMaxScale)
		if v := getStringValue(config, "unit"); v != "" {
			widget.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
	}

	return dashboards.ApiStackedAreaChartWidgetAsApiWidget(widget)
}

// buildColorGridWidget builds a Color Grid widget from Terraform data
func buildColorGridWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiColorGridWidget("Color Grid")
	setCommonBuilderFields(widget, data)

	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.ColorGridDatasource(dataSource))
	}

	if configList := getListValue(data, "color_grid_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setFloat32FromMapIfPresent(config, "min_scale", widget.SetMinScale)
		setFloat32FromMapIfPresent(config, "max_scale", widget.SetMaxScale)
		if v := getStringValue(config, "unit"); v != "" {
			widget.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v := getStringValue(config, "cards"); v != "" {
			widget.SetCards(dashboards.ApiAggregateProperty(v))
		}
		if v := getStringValue(config, "group_cards_by"); v != "" {
			widget.SetGroupCardsBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getIntValue(config, "columns"); v != 0 {
			widget.SetColumns(int32(v))
		}
		if v := getIntValue(config, "limit"); v != 0 {
			widget.SetLimit(int32(v))
		}
	}

	return dashboards.ApiColorGridWidgetAsApiWidget(widget)
}

// buildTimeseriesWidget builds a Timeseries widget from Terraform data
func buildTimeseriesWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiTimeseriesWidget("Time Series: Line")
	setCommonBuilderFields(widget, data)

	// Set data_source (Timeseries-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.TimeseriesDatasource(dataSource))
	}

	if configList := getListValue(data, "timeseries_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		setFloat32FromMapIfPresent(config, "min_scale", widget.SetMinScale)
		setFloat32FromMapIfPresent(config, "max_scale", widget.SetMaxScale)
		if v := getStringValue(config, "unit"); v != "" {
			widget.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v := getStringValue(config, "group_by"); v != "" {
			widget.SetGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v, ok := config["show_timeseries_overall_baseline"].(bool); ok {
			widget.SetShowTimeseriesOverallBaseline(v)
		}
		if v, ok := config["is_timeseries_one_chart_per_line"].(bool); ok {
			widget.SetIsTimeseriesOneChartPerLine(v)
		}
	}

	return dashboards.ApiTimeseriesWidgetAsApiWidget(widget)
}

// buildAgentStatusWidget builds an Agent Status widget from Terraform data
func buildAgentStatusWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiAgentStatusWidget("Agent Status")
	setCommonBuilderFields(widget, data)

	// Set data_source (AgentStatus-specific type)
	if dataSource := getStringValue(data, "data_source"); dataSource != "" {
		widget.SetDataSource(dashboards.AgentStatusDatasource(dataSource))
	}

	if configList := getListValue(data, "agent_status_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v := getStringValue(config, "show"); v != "" {
			widget.SetShow(dashboards.LegacyAgentWidgetShow(v))
		}
		if v := getStringValue(config, "agent_type"); v != "" {
			widget.SetAgents(dashboards.LegacyAgentWidgetType(v))
		}
	}

	return dashboards.ApiAgentStatusWidgetAsApiWidget(widget)
}

// buildNumberWidget builds a Number widget from Terraform data
func buildNumberWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiNumbersCardWidget("Number")
	setCommonBuilderFields(widget, data)

	cardsList := getListValue(data, "number_cards")
	if len(cardsList) > 0 {
		widget.SetNumberCards(buildNumberCards(cardsList))
	} else {
		widget.SetNumberCards([]dashboards.ApiNumbersCard{})
	}

	return dashboards.ApiNumbersCardWidgetAsApiWidget(widget)
}

func buildNumberCards(cardsList []interface{}) []dashboards.ApiNumbersCard {
	cards := make([]dashboards.ApiNumbersCard, 0, len(cardsList))
	for _, c := range cardsList {
		cardData, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		card := dashboards.NewApiNumbersCard()

		if v := getStringValue(cardData, "id"); v != "" {
			card.SetId(v)
		}
		if v := getStringValue(cardData, "description"); v != "" {
			card.SetDescription(v)
		}
		setFloat32FromMapIfPresent(cardData, "min_scale", card.SetMinScale)
		setFloat32FromMapIfPresent(cardData, "max_scale", card.SetMaxScale)
		if v := getStringValue(cardData, "unit"); v != "" {
			card.SetUnit(dashboards.ApiWidgetFixedYScalePrefix(v))
		}
		if v, ok := boolFromMapIfPresent(cardData, "compare_to_previous_value"); ok {
			card.SetCompareToPreviousValue(v)
		}
		if v, ok := boolFromMapIfPresent(cardData, "should_exclude_alert_suppression_windows"); ok {
			card.SetShouldExcludeAlertSuppressionWindows(v)
		}
		if dataSource := getStringValue(cardData, "data_source"); dataSource != "" {
			card.SetDataSource(dashboards.NumbersCardDatasource(dataSource))
		}
		if v := getStringValue(cardData, "metric_group"); v != "" {
			card.SetMetricGroup(dashboards.MetricGroup(v))
		}
		if v := getStringValue(cardData, "direction"); v != "" {
			card.SetDirection(dashboards.DashboardMetricDirection(v))
		}
		if v := getStringValue(cardData, "metric"); v != "" {
			card.SetMetric(dashboards.DashboardMetric(v))
		}

		if measureList := getListValue(cardData, "measure"); len(measureList) > 0 {
			measureData := measureList[0].(map[string]interface{})
			m := dashboards.NewApiWidgetMeasure()
			if v := getStringValue(measureData, "type"); v != "" {
				m.SetType(dashboards.WidgetMeasureType(v))
			}
			if v := getFloat64Value(measureData, "percentile_value"); v != 0 {
				m.SetPercentileValue(float32(v))
			}
			card.SetMeasure(*m)
		}

		if filterList := getSetValue(cardData, "filter"); len(filterList) > 0 {
			apiFilters := make(map[string][]interface{})
			for _, f := range filterList {
				filterData := f.(map[string]interface{})
				property := getStringValue(filterData, "property")
				if property == "" {
					continue
				}
				var values []interface{}
				switch v := filterData["values"].(type) {
				case *schema.Set:
					strs := make([]string, 0, v.Len())
					for _, item := range v.List() {
						strs = append(strs, item.(string))
					}
					sort.Strings(strs)
					values = make([]interface{}, len(strs))
					for i, s := range strs {
						values[i] = s
					}
				case []interface{}:
					values = v
				}
				if len(values) > 0 {
					apiFilters[property] = values
				}
			}
			if len(apiFilters) > 0 {
				card.SetFilters(apiFilters)
			}
		}

		if fixedTimespanList := getListValue(cardData, "fixed_timespan"); len(fixedTimespanList) > 0 {
			fixedTimespan := fixedTimespanList[0].(map[string]interface{})
			duration := dashboards.NewApiDuration()
			if v := getIntValue(fixedTimespan, "value"); v != 0 {
				duration.SetValue(int32(v))
			}
			if v := getStringValue(fixedTimespan, "unit"); v != "" {
				duration.SetUnit(dashboards.LegacyDurationUnit(v))
			}
			card.SetFixedTimespan(*duration)
		}

		cards = append(cards, *card)
	}
	return cards
}

// buildMultiMetricTableWidget builds a Multi Metric Table widget from Terraform data
func buildMultiMetricTableWidget(data map[string]interface{}) dashboards.ApiWidget {
	widget := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
	setCommonBuilderFields(widget, data)

	if configList := getListValue(data, "multi_metric_table_config"); len(configList) > 0 {
		config := configList[0].(map[string]interface{})
		if v, ok := boolFromMapIfPresent(config, "compare_to_previous_value"); ok {
			widget.SetCompareToPreviousValue(v)
		}
		if v := getStringValue(config, "row_group_by"); v != "" {
			widget.SetRowGroupBy(dashboards.ApiAggregateProperty(v))
		}
		if v := getIntValue(config, "limit"); v != 0 {
			widget.SetLimit(int32(v))
		}
	}

	columnsList := getListValue(data, "multi_metric_columns")
	if len(columnsList) > 0 {
		widget.SetMultiMetricColumns(buildMultiMetricColumns(columnsList))
	} else {
		widget.SetMultiMetricColumns([]dashboards.ApiMultiMetricColumn{})
	}

	return dashboards.ApiMultiMetricTableWidgetAsApiWidget(widget)
}

func buildMultiMetricColumns(columnsList []interface{}) []dashboards.ApiMultiMetricColumn {
	columns := make([]dashboards.ApiMultiMetricColumn, 0, len(columnsList))
	for _, c := range columnsList {
		colData, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		col := dashboards.NewApiMultiMetricColumn()

		if v := getStringValue(colData, "id"); v != "" {
			col.SetId(v)
		}
		if dataSource := getStringValue(colData, "data_source"); dataSource != "" {
			col.SetDataSource(dashboards.MultiMetricsTableDatasource(dataSource))
		}
		if v := getStringValue(colData, "metric_group"); v != "" {
			col.SetMetricGroup(dashboards.MetricGroup(v))
		}
		if v := getStringValue(colData, "direction"); v != "" {
			col.SetDirection(dashboards.DashboardMetricDirection(v))
		}
		if v := getStringValue(colData, "metric"); v != "" {
			col.SetMetric(dashboards.DashboardMetric(v))
		}

		if measureList := getListValue(colData, "measure"); len(measureList) > 0 {
			measureData := measureList[0].(map[string]interface{})
			m := dashboards.NewApiWidgetMeasure()
			if v := getStringValue(measureData, "type"); v != "" {
				m.SetType(dashboards.WidgetMeasureType(v))
			}
			if v := getFloat64Value(measureData, "percentile_value"); v != 0 {
				m.SetPercentileValue(float32(v))
			}
			col.SetMeasure(*m)
		}

		if filterList := getSetValue(colData, "filter"); len(filterList) > 0 {
			apiFilters := make(map[string][]interface{})
			for _, f := range filterList {
				filterData := f.(map[string]interface{})
				property := getStringValue(filterData, "property")
				if property == "" {
					continue
				}
				var values []interface{}
				switch v := filterData["values"].(type) {
				case *schema.Set:
					strs := make([]string, 0, v.Len())
					for _, item := range v.List() {
						strs = append(strs, item.(string))
					}
					sort.Strings(strs)
					values = make([]interface{}, len(strs))
					for i, s := range strs {
						values[i] = s
					}
				case []interface{}:
					values = v
				}
				if len(values) > 0 {
					apiFilters[property] = values
				}
			}
			if len(apiFilters) > 0 {
				col.SetFilters(apiFilters)
			}
		}

		columns = append(columns, *col)
	}
	return columns
}

// setCommonBuilderFields sets common fields on any widget
func setCommonBuilderFields(widget interface{}, data map[string]interface{}) {
	// Set widget ID if present (important for updates)
	if id := getStringValue(data, "id"); id != "" {
		if w, ok := widget.(interface{ SetId(string) }); ok {
			w.SetId(id)
		}
	}
	if title := getStringValue(data, "title"); title != "" {
		if w, ok := widget.(interface{ SetTitle(string) }); ok {
			w.SetTitle(title)
		}
	}
	if visualMode := getStringValue(data, "visual_mode"); visualMode != "" {
		if w, ok := widget.(interface{ SetVisualMode(dashboards.VisualMode) }); ok {
			w.SetVisualMode(dashboards.VisualMode(visualMode))
		}
	}
	if metricGroup := getStringValue(data, "metric_group"); metricGroup != "" {
		if w, ok := widget.(interface{ SetMetricGroup(dashboards.MetricGroup) }); ok {
			w.SetMetricGroup(dashboards.MetricGroup(metricGroup))
		}
	}
	if metric := getStringValue(data, "metric"); metric != "" {
		if w, ok := widget.(interface {
			SetMetric(dashboards.DashboardMetric)
		}); ok {
			w.SetMetric(dashboards.DashboardMetric(metric))
		}
	}
	// Handle measure - nested block with type and percentile_value
	if measureList := getListValue(data, "measure"); len(measureList) > 0 {
		measureData := measureList[0].(map[string]interface{})
		m := dashboards.NewApiWidgetMeasure()
		if v := getStringValue(measureData, "type"); v != "" {
			m.SetType(dashboards.WidgetMeasureType(v))
		}
		if v := getFloat64Value(measureData, "percentile_value"); v != 0 {
			m.SetPercentileValue(float32(v))
		}
		if w, ok := widget.(interface {
			SetMeasure(dashboards.ApiWidgetMeasure)
		}); ok {
			w.SetMeasure(*m)
		}
	}
	if shouldExclude, ok := boolFromMapIfPresent(data, "should_exclude_alert_suppression_windows"); ok {
		if w, ok := widget.(interface{ SetShouldExcludeAlertSuppressionWindows(bool) }); ok {
			w.SetShouldExcludeAlertSuppressionWindows(shouldExclude)
		}
	}

	// Handle filter blocks - SDK uses map[string][]interface{}
	if filterList := getSetValue(data, "filter"); len(filterList) > 0 {
		apiFilters := make(map[string][]interface{})
		for _, f := range filterList {
			filterData := f.(map[string]interface{})
			property := getStringValue(filterData, "property")
			if property == "" {
				continue
			}
			var values []interface{}
			if valuesSet, ok := filterData["values"].(*schema.Set); ok {
				strs := make([]string, 0, valuesSet.Len())
				for _, v := range valuesSet.List() {
					strs = append(strs, v.(string))
				}
				sort.Strings(strs)
				values = make([]interface{}, len(strs))
				for i, s := range strs {
					values[i] = s
				}
			}
			if len(values) > 0 {
				apiFilters[property] = values
			}
		}
		if w, ok := widget.(interface {
			SetFilters(map[string][]interface{})
		}); ok && len(apiFilters) > 0 {
			w.SetFilters(apiFilters)
		}
	}

	// Handle fixed_timespan
	if fixedTimespanList := getListValue(data, "fixed_timespan"); len(fixedTimespanList) > 0 {
		fixedTimespan := fixedTimespanList[0].(map[string]interface{})
		duration := dashboards.NewApiDuration()
		if v := getIntValue(fixedTimespan, "value"); v != 0 {
			duration.SetValue(int32(v))
		}
		if v := getStringValue(fixedTimespan, "unit"); v != "" {
			duration.SetUnit(dashboards.LegacyDurationUnit(v))
		}
		if w, ok := widget.(interface{ SetFixedTimespan(dashboards.ApiDuration) }); ok {
			w.SetFixedTimespan(*duration)
		}
	}
}

func buildActiveWithinFromConfig(config map[string]interface{}) *dashboards.ActiveWithin {
	activeWithinValue := 0
	activeWithinUnit := ""
	if v, ok := config["active_within_value"].(int); ok && v != 0 {
		activeWithinValue = v
	}
	if v := getStringValue(config, "active_within_unit"); v != "" {
		activeWithinUnit = v
	}
	if activeWithinValue == 0 && activeWithinUnit == "" {
		return nil
	}

	activeWithin := dashboards.NewActiveWithin()
	if activeWithinValue != 0 {
		activeWithin.SetValue(int32(activeWithinValue))
	}
	if activeWithinUnit != "" {
		activeWithin.SetUnit(dashboards.LegacyDurationUnit(activeWithinUnit))
	}
	return activeWithin
}

func setAlertTypesFromConfig(m map[string]interface{}, key string, set func([]dashboards.LegacyAlertListAlertType)) {
	raw, ok := m[key].(*schema.Set)
	if !ok || raw.Len() == 0 {
		return
	}

	values := make([]string, 0, raw.Len())
	for _, v := range raw.List() {
		values = append(values, v.(string))
	}
	sort.Strings(values)

	alertTypes := make([]dashboards.LegacyAlertListAlertType, len(values))
	for i, v := range values {
		alertTypes[i] = dashboards.LegacyAlertListAlertType(v)
	}
	set(alertTypes)
}

func buildTestTableFilterConfig(config map[string]interface{}, key string) *dashboards.ApiWidgetFilterApiTestTableFilterKey {
	blocks := getListValue(config, key)
	if len(blocks) == 0 {
		return nil
	}

	block, ok := blocks[0].(map[string]interface{})
	if !ok || block == nil {
		return nil
	}

	filter := dashboards.NewApiWidgetFilterApiTestTableFilterKey()
	if v := getStringValue(block, "type"); v != "" {
		filter.SetType(dashboards.TestTableFilterType(v))
	}

	terms := getListValue(block, "filters")
	if len(terms) > 0 {
		items := make([]dashboards.ApiMultiSearchFilterApiTestTableFilterKey, 0, len(terms))
		for _, rawTerm := range terms {
			term, ok := rawTerm.(map[string]interface{})
			if !ok || term == nil {
				continue
			}
			keyValue := getStringValue(term, "key")
			value := getStringValue(term, "value")
			if keyValue == "" || value == "" {
				continue
			}
			item := dashboards.NewApiMultiSearchFilterApiTestTableFilterKey()
			item.SetKey(dashboards.TestTableFilterKey(keyValue))
			item.SetValue(value)
			items = append(items, *item)
		}
		if len(items) > 0 {
			filter.SetFilters(items)
		}
	}

	if !filter.HasType() && !filter.HasFilters() {
		return nil
	}
	return filter
}
