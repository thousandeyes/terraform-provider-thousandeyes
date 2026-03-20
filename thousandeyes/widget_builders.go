package thousandeyes

import (
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
		if v := getFloat64Value(config, "min_scale"); v != 0 {
			widget.SetMinScale(float32(v))
		}
		if v := getFloat64Value(config, "max_scale"); v != 0 {
			widget.SetMaxScale(float32(v))
		}
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
		if v := getFloat64Value(config, "min_scale"); v != 0 {
			widget.SetMinScale(float32(v))
		}
		if v := getFloat64Value(config, "max_scale"); v != 0 {
			widget.SetMaxScale(float32(v))
		}
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
	if isEmbedded := getBoolValue(data, "is_embedded"); isEmbedded {
		if w, ok := widget.(interface{ SetIsEmbedded(bool) }); ok {
			w.SetIsEmbedded(isEmbedded)
		}
	}
	if metricGroup := getStringValue(data, "metric_group"); metricGroup != "" {
		if w, ok := widget.(interface{ SetMetricGroup(dashboards.MetricGroup) }); ok {
			w.SetMetricGroup(dashboards.MetricGroup(metricGroup))
		}
	}
	if direction := getStringValue(data, "direction"); direction != "" {
		if w, ok := widget.(interface {
			SetDirection(dashboards.DashboardMetricDirection)
		}); ok {
			w.SetDirection(dashboards.DashboardMetricDirection(direction))
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
	if shouldExclude := getBoolValue(data, "should_exclude_alert_suppression_windows"); shouldExclude {
		if w, ok := widget.(interface{ SetShouldExcludeAlertSuppressionWindows(bool) }); ok {
			w.SetShouldExcludeAlertSuppressionWindows(shouldExclude)
		}
	}

	// Handle filters - SDK uses map[string][]interface{}
	if filters, ok := data["filters"].(map[string]interface{}); ok && len(filters) > 0 {
		apiFilters := make(map[string][]interface{})
		for key, val := range filters {
			apiFilters[key] = []interface{}{val}
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
