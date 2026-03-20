package thousandeyes

import (
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

// WidgetBuilder builds an API widget from Terraform resource data (map)
// Used for Create and Update operations
type WidgetBuilder func(data map[string]interface{}) dashboards.ApiWidget

// WidgetMapper maps an API widget back to Terraform resource data (map)
// Used for Read operations
type WidgetMapper func(widget dashboards.ApiWidget) map[string]interface{}

// WidgetTypeRegistry maps widget type strings to their builder and mapper functions
type WidgetTypeRegistry struct {
	Builder WidgetBuilder
	Mapper  WidgetMapper
}

// widgetRegistry holds the mapping functions for each widget type
var widgetRegistry = map[string]WidgetTypeRegistry{
	"Map":                      {Builder: buildGeoMapWidget, Mapper: mapGeoMapWidget},
	"Agent Status":             {Builder: buildAgentStatusWidget, Mapper: mapAgentStatusWidget},
	"Time Series: Line":        {Builder: buildTimeseriesWidget, Mapper: mapTimeseriesWidget},
	"Time Series: Stacked Area": {Builder: buildStackedAreaWidget, Mapper: mapStackedAreaWidget},
}

// BuildWidget builds an API widget from Terraform data using the appropriate builder
func BuildWidget(data map[string]interface{}) dashboards.ApiWidget {
	widgetType, ok := data["type"].(string)
	if !ok {
		return dashboards.ApiWidget{}
	}

	registry, exists := widgetRegistry[widgetType]
	if !exists {
		return dashboards.ApiWidget{}
	}

	return registry.Builder(data)
}

// MapWidget maps an API widget to Terraform data using the appropriate mapper
func MapWidget(widget dashboards.ApiWidget) map[string]interface{} {
	instance := widget.GetActualInstance()
	if instance == nil {
		return nil
	}

	var widgetType string
	switch instance.(type) {
	case *dashboards.ApiGeoMapWidget:
		widgetType = "Map"
	case *dashboards.ApiAgentStatusWidget:
		widgetType = "Agent Status"
	case *dashboards.ApiTimeseriesWidget:
		widgetType = "Time Series: Line"
	case *dashboards.ApiStackedAreaChartWidget:
		widgetType = "Time Series: Stacked Area"
	default:
		return nil
	}

	registry, exists := widgetRegistry[widgetType]
	if !exists {
		return nil
	}

	return registry.Mapper(widget)
}

// BuildWidgets builds a slice of API widgets from Terraform data
func BuildWidgets(widgetsData []interface{}) []dashboards.ApiWidget {
	if len(widgetsData) == 0 {
		return nil
	}

	widgets := make([]dashboards.ApiWidget, 0, len(widgetsData))
	for _, w := range widgetsData {
		if widgetData, ok := w.(map[string]interface{}); ok {
			widget := BuildWidget(widgetData)
			widgets = append(widgets, widget)
		}
	}
	return widgets
}

// MapWidgets maps a slice of API widgets to Terraform data
func MapWidgets(widgets []dashboards.ApiWidget) []interface{} {
	if len(widgets) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(widgets))
	for _, widget := range widgets {
		if mapped := MapWidget(widget); mapped != nil {
			result = append(result, mapped)
		}
	}
	return result
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

func getBoolValue(data map[string]interface{}, key string) bool {
	if v, ok := data[key].(bool); ok {
		return v
	}
	return false
}

func getFloat64Value(data map[string]interface{}, key string) float64 {
	if v, ok := data[key].(float64); ok {
		return v
	}
	return 0
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
