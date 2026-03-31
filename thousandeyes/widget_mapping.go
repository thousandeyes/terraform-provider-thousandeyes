package thousandeyes

import (
	"fmt"
	"log"

	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// Widget type constants
const (
	WidgetTypeMap            = "Map"
	WidgetTypeAgentStatus    = "Agent Status"
	WidgetTypeTimeseriesLine = "Time Series: Line"
	WidgetTypeStackedArea    = "Time Series: Stacked Area"
	WidgetTypePieChart       = "Pie Chart"
	WidgetTypeBoxAndWhiskers = "Box and Whiskers"
	WidgetTypeList           = "List"
)

// WidgetBuilder builds an API widget from Terraform resource data (map)
// Used for Create and Update operations
type WidgetBuilder func(data map[string]interface{}) dashboards.ApiWidget

// WidgetMapper maps an API widget back to Terraform resource data (map)
// Used for Read operations
type WidgetMapper func(widget dashboards.ApiWidget) (map[string]interface{}, error)

// WidgetTypeRegistry maps widget type strings to their builder and mapper functions
type WidgetTypeRegistry struct {
	Builder WidgetBuilder
	Mapper  WidgetMapper
}

// widgetRegistry holds the mapping functions for each widget type
var widgetRegistry = map[string]WidgetTypeRegistry{
	WidgetTypeMap:            {Builder: buildGeoMapWidget, Mapper: mapGeoMapWidget},
	WidgetTypeAgentStatus:    {Builder: buildAgentStatusWidget, Mapper: mapAgentStatusWidget},
	WidgetTypeTimeseriesLine: {Builder: buildTimeseriesWidget, Mapper: mapTimeseriesWidget},
	WidgetTypeStackedArea:    {Builder: buildStackedAreaWidget, Mapper: mapStackedAreaWidget},
	WidgetTypePieChart:       {Builder: buildPieChartWidget, Mapper: mapPieChartWidget},
	WidgetTypeBoxAndWhiskers: {Builder: buildBoxAndWhiskersWidget, Mapper: mapBoxAndWhiskersWidget},
	WidgetTypeList:           {Builder: buildListWidget, Mapper: mapListWidget},
}

// BuildWidget builds an API widget from Terraform data using the appropriate builder
func BuildWidget(data map[string]interface{}) (dashboards.ApiWidget, error) {
	widgetType, ok := data["type"].(string)
	if !ok || widgetType == "" {
		return dashboards.ApiWidget{}, fmt.Errorf("type is required and must be a non-empty string")
	}

	registry, exists := widgetRegistry[widgetType]
	if !exists {
		return dashboards.ApiWidget{}, fmt.Errorf("unsupported widget type %q", widgetType)
	}

	return registry.Builder(data), nil
}

// widgetTypeFromInstance maps a concrete API widget instance to the Terraform "type" string key.
func widgetTypeFromInstance(instance interface{}) (string, error) {
	switch instance.(type) {
	case *dashboards.ApiGeoMapWidget:
		return WidgetTypeMap, nil
	case *dashboards.ApiAgentStatusWidget:
		return WidgetTypeAgentStatus, nil
	case *dashboards.ApiTimeseriesWidget:
		return WidgetTypeTimeseriesLine, nil
	case *dashboards.ApiStackedAreaChartWidget:
		return WidgetTypeStackedArea, nil
	case *dashboards.ApiPieChartWidget:
		return WidgetTypePieChart, nil
	case *dashboards.ApiBoxAndWhiskersWidget:
		return WidgetTypeBoxAndWhiskers, nil
	case *dashboards.ApiListWidget:
		return WidgetTypeList, nil
	default:
		return "", nil
	}
}

// mapWidgetWithInstance maps using instance as returned by GetActualInstance (tests may pass a synthetic type).
func mapWidgetWithInstance(widget dashboards.ApiWidget, instance interface{}) (map[string]interface{}, error) {
	if instance == nil {
		return nil, fmt.Errorf("widget instance is nil")
	}

	widgetType, err := widgetTypeFromInstance(instance)
	if err != nil {
		return nil, err
	}
	if widgetType == "" {
		log.Printf("[WARN] Skipping unmanaged widget type: %T", instance)
		return nil, nil
	}

	registry, exists := widgetRegistry[widgetType]
	if !exists {
		return nil, fmt.Errorf("no mapper registered for widget type: %s", widgetType)
	}

	return registry.Mapper(widget)
}

// MapWidget maps an API widget to Terraform data using the appropriate mapper
func MapWidget(widget dashboards.ApiWidget) (map[string]interface{}, error) {
	return mapWidgetWithInstance(widget, widget.GetActualInstance())
}

// BuildWidgets builds a slice of API widgets from Terraform data
func BuildWidgets(widgetsData []interface{}) ([]dashboards.ApiWidget, error) {
	if len(widgetsData) == 0 {
		return nil, nil
	}

	widgets := make([]dashboards.ApiWidget, 0, len(widgetsData))
	for i, w := range widgetsData {
		if w == nil {
			return nil, fmt.Errorf("widgets.%d: must not be null", i)
		}
		widgetData, ok := w.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("widgets.%d: expected object, got %T", i, w)
		}
		widget, err := BuildWidget(widgetData)
		if err != nil {
			return nil, fmt.Errorf("widgets.%d: %w", i, err)
		}
		widgets = append(widgets, widget)
	}
	return widgets, nil
}

// mapAllWidgets maps each widget with mapOne (MapWidgets passes MapWidget).
func mapAllWidgets(widgets []dashboards.ApiWidget, mapOne func(dashboards.ApiWidget) (map[string]interface{}, error)) ([]interface{}, error) {
	if len(widgets) == 0 {
		return nil, nil
	}

	result := make([]interface{}, 0, len(widgets))
	for _, widget := range widgets {
		mapped, err := mapOne(widget)
		if err != nil {
			return nil, err
		}
		if mapped != nil {
			result = append(result, mapped)
		}
	}
	return result, nil
}

// MapWidgets maps a slice of API widgets to Terraform data
func MapWidgets(widgets []dashboards.ApiWidget) ([]interface{}, error) {
	return mapAllWidgets(widgets, MapWidget)
}

// isUnmanagedWidget reports whether the given API widget is a type not
// supported by this provider's widget registry.
func isUnmanagedWidget(w dashboards.ApiWidget) bool {
	instance := w.GetActualInstance()
	if instance == nil {
		return false
	}
	wType, _ := widgetTypeFromInstance(instance)
	return wType == ""
}

// mergeUnmanagedWidgets combines managed widgets from config with unmanaged
// widgets from the current API state. Unmanaged widgets are kept at their
// current positions; managed slots are filled in config order. Extra config
// widgets (newly added) are appended at the end.
func mergeUnmanagedWidgets(configWidgets, currentAPIWidgets []dashboards.ApiWidget) []dashboards.ApiWidget {
	merged := make([]dashboards.ApiWidget, 0, len(currentAPIWidgets)+len(configWidgets))
	mi := 0

	for _, apiWidget := range currentAPIWidgets {
		if isUnmanagedWidget(apiWidget) {
			merged = append(merged, apiWidget)
		} else if mi < len(configWidgets) {
			merged = append(merged, configWidgets[mi])
			mi++
		}
	}
	for ; mi < len(configWidgets); mi++ {
		merged = append(merged, configWidgets[mi])
	}
	return merged
}

// Helper functions for extracting values from map data

func getStringValue(data map[string]interface{}, key string) string {
	if v, ok := data[key].(string); ok {
		return v
	}
	return ""
}

func getIntValue(data map[string]interface{}, key string) int {
	if v, ok := data[key].(int); ok {
		return v
	}
	return 0
}

func boolFromMapIfPresent(m map[string]interface{}, key string) (value bool, ok bool) {
	v, present := m[key]
	if !present {
		return false, false
	}
	b, typed := v.(bool)
	return b, typed
}

func getFloat64Value(data map[string]interface{}, key string) float64 {
	if v, ok := data[key].(float64); ok {
		return v
	}
	return 0
}

// asTerraformFloat coerces Terraform/SDK numeric values to float64.
func asTerraformFloat(v interface{}) (float64, bool) {
	switch x := v.(type) {
	case float64:
		return x, true
	case float32:
		return float64(x), true
	case int:
		return float64(x), true
	case int64:
		return float64(x), true
	default:
		return 0, false
	}
}

// setFloat32FromMapIfPresent calls set when key exists and the value is numeric (including 0).
func setFloat32FromMapIfPresent(m map[string]interface{}, key string, set func(float32)) {
	raw, ok := m[key]
	if !ok {
		return
	}
	f, ok := asTerraformFloat(raw)
	if !ok {
		return
	}
	set(float32(f))
}

func getListValue(data map[string]interface{}, key string) []interface{} {
	if v, ok := data[key].([]interface{}); ok {
		return v
	}
	return nil
}

// setCommonWidgetFields sets common fields on a widget data map
func setCommonWidgetFields(data map[string]interface{}, id, title, embedUrl string, isEmbedded bool, visualMode string) {
	if id != "" {
		data["id"] = id
	}
	if title != "" {
		data["title"] = title
	}
	if embedUrl != "" {
		data["embed_url"] = embedUrl
	}
	data["is_embedded"] = isEmbedded
	if visualMode != "" {
		data["visual_mode"] = visualMode
	}
}
