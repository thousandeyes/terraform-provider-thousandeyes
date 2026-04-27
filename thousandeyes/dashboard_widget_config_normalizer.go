package thousandeyes

import "github.com/hashicorp/go-cty/cty"

var dashboardWidgetPresenceSensitiveConfigBlocks = map[string][]string{
	"geo_map_config":            {"min_scale", "max_scale", "is_geo_map_per_test"},
	"timeseries_config":         {"min_scale", "max_scale", "show_timeseries_overall_baseline", "is_timeseries_one_chart_per_line"},
	"stacked_area_config":       {"min_scale", "max_scale"},
	"box_and_whiskers_config":   {"min_scale", "max_scale"},
	"color_grid_config":         {"min_scale", "max_scale"},
	"alert_list_config":         {"limit_to"},
	"table_config":              {"compare_to_previous_value"},
	"stacked_bar_chart_config":  {"show_labels", "is_horizontal_bar_chart"},
	"grouped_bar_chart_config":  {"show_labels", "is_horizontal_bar_chart"},
	"multi_metric_table_config": {"compare_to_previous_value"},
}

func normalizeConfiguredWidgets(widgetList []interface{}, rawConfig cty.Value) []interface{} {
	rawWidgets, ok := rawWidgetsFromConfig(rawConfig)
	if !ok {
		return widgetList
	}

	normalized := cloneInterfaceSlice(widgetList)
	for i := range normalized {
		if i >= len(rawWidgets) {
			break
		}
		widget, ok := normalized[i].(map[string]interface{})
		if !ok {
			continue
		}
		normalizeConfiguredWidget(widget, rawWidgets[i])
	}
	return normalized
}

func normalizeConfiguredWidget(widget map[string]interface{}, rawWidget cty.Value) {
	pruneConfiguredFields(widget, rawWidget, []string{
		"should_exclude_alert_suppression_windows",
	})
	for blockName, fieldNames := range dashboardWidgetPresenceSensitiveConfigBlocks {
		pruneConfiguredBlockFields(widget, rawWidget, blockName, fieldNames)
	}
	pruneConfiguredBlockFields(widget, rawWidget, "number_cards", []string{
		"min_scale",
		"max_scale",
		"compare_to_previous_value",
		"should_exclude_alert_suppression_windows",
	})
}

func pruneConfiguredFields(parent map[string]interface{}, rawParent cty.Value, fieldNames []string) {
	for _, fieldName := range fieldNames {
		if !rawObjectHasConfiguredAttr(rawParent, fieldName) {
			delete(parent, fieldName)
		}
	}
}

func pruneConfiguredBlockFields(parent map[string]interface{}, rawParent cty.Value, blockName string, fieldNames []string) {
	rawBlocks, ok := rawBlockValues(rawParent, blockName)
	if !ok {
		return
	}

	blocks, ok := parent[blockName].([]interface{})
	if !ok {
		return
	}

	for i := range blocks {
		if i >= len(rawBlocks) {
			break
		}
		block, ok := blocks[i].(map[string]interface{})
		if !ok {
			continue
		}
		for _, fieldName := range fieldNames {
			if !rawObjectHasConfiguredAttr(rawBlocks[i], fieldName) {
				delete(block, fieldName)
			}
		}
	}
}

func rawWidgetsFromConfig(rawConfig cty.Value) ([]cty.Value, bool) {
	if rawConfig == cty.NilVal || rawConfig.IsNull() || !rawConfig.IsKnown() {
		return nil, false
	}
	rawType := rawConfig.Type()
	if !rawType.IsObjectType() || !rawType.HasAttribute("widgets") {
		return nil, false
	}
	widgets := rawConfig.GetAttr("widgets")
	rawWidgets, ok := ctyListValues(widgets)
	if !ok || len(rawWidgets) == 0 {
		return nil, false
	}
	return rawWidgets, true
}

func rawBlockValues(rawParent cty.Value, blockName string) ([]cty.Value, bool) {
	if !rawObjectHasConfiguredAttr(rawParent, blockName) {
		return nil, false
	}
	return ctyListValues(rawParent.GetAttr(blockName))
}

func rawObjectHasConfiguredAttr(rawObject cty.Value, attrName string) bool {
	if rawObject == cty.NilVal || rawObject.IsNull() || !rawObject.IsKnown() {
		return false
	}
	rawType := rawObject.Type()
	if !rawType.IsObjectType() || !rawType.HasAttribute(attrName) {
		return false
	}
	attr := rawObject.GetAttr(attrName)
	return !attr.IsNull()
}

func ctyListValues(value cty.Value) ([]cty.Value, bool) {
	if value == cty.NilVal || value.IsNull() || !value.IsKnown() {
		return nil, false
	}
	valueType := value.Type()
	if !valueType.IsListType() && !valueType.IsTupleType() {
		return nil, false
	}
	return value.AsValueSlice(), true
}

func cloneInterfaceSlice(input []interface{}) []interface{} {
	output := make([]interface{}, len(input))
	for i, value := range input {
		output[i] = cloneInterfaceValue(value)
	}
	return output
}

func cloneInterfaceValue(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[string]interface{}:
		copied := make(map[string]interface{}, len(typed))
		for key, nested := range typed {
			copied[key] = cloneInterfaceValue(nested)
		}
		return copied
	case []interface{}:
		return cloneInterfaceSlice(typed)
	default:
		return value
	}
}
