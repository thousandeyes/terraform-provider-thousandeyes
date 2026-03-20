package thousandeyes

import (
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// mapGeoMapWidget maps a GeoMap widget to Terraform data
func mapGeoMapWidget(widget dashboards.ApiWidget) map[string]interface{} {
	w := widget.ApiGeoMapWidget
	if w == nil {
		return nil
	}

	data := map[string]interface{}{
		"type": "Map",
	}
	setCommonWidgetFields(data, w.GetId(), w.GetTitle(), w.GetEmbedUrl(), w.GetIsEmbedded(), string(w.GetVisualMode()))
	setCommonMapperFields(data, w)

	// Map data_source (GeoMap-specific type)
	if v := w.GetDataSource(); v != "" {
		data["data_source"] = string(v)
	}

	config := map[string]interface{}{}
	if v, ok := w.GetMinScaleOk(); ok && v != nil {
		config["min_scale"] = float64(*v)
	}
	if v, ok := w.GetMaxScaleOk(); ok && v != nil {
		config["max_scale"] = float64(*v)
	}
	if v := w.GetUnit(); v != "" {
		config["unit"] = string(v)
	}
	if v := w.GetGroupBy(); v != "" {
		config["group_by"] = string(v)
	}
	if v, ok := w.GetIsGeoMapPerTestOk(); ok && v != nil {
		config["is_geo_map_per_test"] = *v
	}
	if len(config) > 0 {
		data["geo_map_config"] = []interface{}{config}
	}

	return data
}

// mapTimeseriesWidget maps a Timeseries widget to Terraform data
func mapTimeseriesWidget(widget dashboards.ApiWidget) map[string]interface{} {
	w := widget.ApiTimeseriesWidget
	if w == nil {
		return nil
	}

	data := map[string]interface{}{
		"type": "Time Series: Line",
	}
	setCommonWidgetFields(data, w.GetId(), w.GetTitle(), w.GetEmbedUrl(), w.GetIsEmbedded(), string(w.GetVisualMode()))
	setCommonMapperFields(data, w)

	// Map data_source (Timeseries-specific type)
	if v := w.GetDataSource(); v != "" {
		data["data_source"] = string(v)
	}

	config := map[string]interface{}{}
	if v, ok := w.GetMinScaleOk(); ok && v != nil {
		config["min_scale"] = float64(*v)
	}
	if v, ok := w.GetMaxScaleOk(); ok && v != nil {
		config["max_scale"] = float64(*v)
	}
	if v := w.GetUnit(); v != "" {
		config["unit"] = string(v)
	}
	if v := w.GetGroupBy(); v != "" {
		config["group_by"] = string(v)
	}
	if v, ok := w.GetShowTimeseriesOverallBaselineOk(); ok && v != nil {
		config["show_timeseries_overall_baseline"] = *v
	}
	if v, ok := w.GetIsTimeseriesOneChartPerLineOk(); ok && v != nil {
		config["is_timeseries_one_chart_per_line"] = *v
	}
	if len(config) > 0 {
		data["timeseries_config"] = []interface{}{config}
	}

	return data
}

// mapAgentStatusWidget maps an Agent Status widget to Terraform data
func mapAgentStatusWidget(widget dashboards.ApiWidget) map[string]interface{} {
	w := widget.ApiAgentStatusWidget
	if w == nil {
		return nil
	}

	data := map[string]interface{}{
		"type": "Agent Status",
	}
	setCommonWidgetFields(data, w.GetId(), w.GetTitle(), w.GetEmbedUrl(), w.GetIsEmbedded(), string(w.GetVisualMode()))
	setCommonMapperFields(data, w)

	// Map data_source (AgentStatus-specific type)
	if v := w.GetDataSource(); v != "" {
		data["data_source"] = string(v)
	}

	config := map[string]interface{}{}
	if v := w.GetShow(); v != "" {
		config["show"] = string(v)
	}
	if v := w.GetAgents(); v != "" {
		config["agent_type"] = string(v)
	}
	if len(config) > 0 {
		data["agent_status_config"] = []interface{}{config}
	}

	return data
}

// setCommonMapperFields sets common fields from any widget that has them
func setCommonMapperFields(data map[string]interface{}, widget interface{}) {
	if w, ok := widget.(interface{ GetMetricGroup() dashboards.MetricGroup }); ok {
		if v := w.GetMetricGroup(); v != "" {
			data["metric_group"] = string(v)
		}
	}
	if w, ok := widget.(interface {
		GetDirection() dashboards.DashboardMetricDirection
	}); ok {
		if v := w.GetDirection(); v != "" {
			data["direction"] = string(v)
		}
	}
	if w, ok := widget.(interface{ GetMetric() dashboards.DashboardMetric }); ok {
		if v := w.GetMetric(); v != "" {
			data["metric"] = string(v)
		}
	}
	// Handle measure - nested block with type and percentile_value
	if w, ok := widget.(interface {
		GetMeasureOk() (*dashboards.ApiWidgetMeasure, bool)
	}); ok {
		if measure, ok := w.GetMeasureOk(); ok && measure != nil {
			measureMap := map[string]interface{}{}
			if measureType := measure.GetType(); measureType != "" {
				measureMap["type"] = string(measureType)
			}
			if percentile, ok := measure.GetPercentileValueOk(); ok && percentile != nil {
				measureMap["percentile_value"] = float64(*percentile)
			}
			if len(measureMap) > 0 {
				data["measure"] = []interface{}{measureMap}
			}
		}
	}
	if w, ok := widget.(interface {
		GetShouldExcludeAlertSuppressionWindowsOk() (*bool, bool)
	}); ok {
		if v, ok := w.GetShouldExcludeAlertSuppressionWindowsOk(); ok && v != nil {
			data["should_exclude_alert_suppression_windows"] = *v
		}
	}

	// Handle fixed_timespan
	if w, ok := widget.(interface {
		GetFixedTimespanOk() (*dashboards.ApiDuration, bool)
	}); ok {
		if fixedTimespan, ok := w.GetFixedTimespanOk(); ok && fixedTimespan != nil {
			fixedTimespanMap := map[string]interface{}{}
			if v, ok := fixedTimespan.GetValueOk(); ok && v != nil {
				fixedTimespanMap["value"] = int(*v)
			}
			if v := fixedTimespan.GetUnit(); v != "" {
				fixedTimespanMap["unit"] = string(v)
			}
			if len(fixedTimespanMap) > 0 {
				data["fixed_timespan"] = []interface{}{fixedTimespanMap}
			}
		}
	}

	// Handle filters - SDK uses map[string][]interface{}
	if w, ok := widget.(interface {
		GetFiltersOk() (*map[string][]interface{}, bool)
	}); ok {
		if filters, ok := w.GetFiltersOk(); ok && filters != nil && len(*filters) > 0 {
			// Convert to map[string]interface{} for Terraform (TypeMap)
			tfFilters := make(map[string]interface{})
			for key, vals := range *filters {
				if len(vals) > 0 {
					// Take first value as string
					if strVal, ok := vals[0].(string); ok {
						tfFilters[key] = strVal
					}
				}
			}
			if len(tfFilters) > 0 {
				data["filters"] = tfFilters
			}
		}
	}
}
