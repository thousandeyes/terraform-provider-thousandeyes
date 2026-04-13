package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildStackedBarChartWidget(t *testing.T) {
	widget := buildStackedBarChartWidget(map[string]interface{}{
		"type":        "Bar Chart: Stacked",
		"title":       "Stacked Bars",
		"visual_mode": "Full",
		"data_source": "CLOUD_NATIVE_MONITORING",
		"stacked_bar_chart_config": []interface{}{
			map[string]interface{}{
				"axis_group_by":           "CLOUD_NATIVE_MONITORING-REGION",
				"limit":                   8,
				"show_labels":             true,
				"is_horizontal_bar_chart": true,
			},
		},
	})

	w := widget.ApiStackedBarchartWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Bar Chart: Stacked", w.GetType())
	assert.Equal(t, dashboards.StackedBarChartDatasource("CLOUD_NATIVE_MONITORING"), w.GetDataSource())
	assert.Equal(t, dashboards.ApiAggregateProperty("CLOUD_NATIVE_MONITORING-REGION"), w.GetAxisGroupBy())
	assert.Equal(t, int32(8), w.GetLimit())
	assert.True(t, w.GetShowLabels())
	assert.True(t, w.GetIsHorizontalBarChart())
}

func TestMapStackedBarChartWidget(t *testing.T) {
	w := dashboards.NewApiStackedBarchartWidget("Bar Chart: Stacked")
	w.SetId("widget-123")
	w.SetTitle("Stacked Bars")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.StackedBarChartDatasource("CLOUD_NATIVE_MONITORING"))
	w.SetAxisGroupBy(dashboards.ApiAggregateProperty("CLOUD_NATIVE_MONITORING-REGION"))
	w.SetSortBy(dashboards.LegacyWidgetSortProperty("Value"))
	w.SetSortDirection(dashboards.LegacyWidgetSortDirection("Descending"))
	w.SetLimit(8)
	w.SetShowLabels(true)
	w.SetIsHorizontalBarChart(true)

	data, err := mapStackedBarChartWidget(dashboards.ApiStackedBarchartWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Bar Chart: Stacked", data["type"])
	assert.Equal(t, "CLOUD_NATIVE_MONITORING", data["data_source"])

	config := data["stacked_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "CLOUD_NATIVE_MONITORING-REGION", config["axis_group_by"])
	assert.NotContains(t, config, "sort_by")
	assert.NotContains(t, config, "sort_direction")
	assert.Equal(t, 8, config["limit"])
	assert.Equal(t, true, config["show_labels"])
	assert.Equal(t, true, config["is_horizontal_bar_chart"])
}
