package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildGroupedBarChartWidget(t *testing.T) {
	widget := buildGroupedBarChartWidget(map[string]interface{}{
		"type":        "Bar Chart: Grouped",
		"title":       "Grouped Bars",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"grouped_bar_chart_config": []interface{}{
			map[string]interface{}{
				"group_by":                "COUNTRY",
				"axis_group_by":           "TEST",
				"sort_by":                 "Alphabetical",
				"sort_direction":          "Ascending",
				"limit":                   12,
				"show_labels":             true,
				"is_horizontal_bar_chart": false,
			},
		},
	})

	w := widget.ApiGroupedBarchartWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Bar Chart: Grouped", w.GetType())
	assert.Equal(t, dashboards.GroupedBarChartDatasource("ALERTS"), w.GetDataSource())
	assert.Equal(t, dashboards.ApiAggregateProperty("COUNTRY"), w.GetGroupBy())
	assert.Equal(t, dashboards.ApiAggregateProperty("TEST"), w.GetAxisGroupBy())
	assert.Equal(t, dashboards.LegacyWidgetSortProperty("Alphabetical"), w.GetSortBy())
	assert.Equal(t, dashboards.LegacyWidgetSortDirection("Ascending"), w.GetSortDirection())
	assert.Equal(t, int32(12), w.GetLimit())
	assert.True(t, w.GetShowLabels())
	assert.False(t, w.GetIsHorizontalBarChart())
}

func TestMapGroupedBarChartWidget(t *testing.T) {
	w := dashboards.NewApiGroupedBarchartWidget("Bar Chart: Grouped")
	w.SetId("widget-123")
	w.SetTitle("Grouped Bars")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.GroupedBarChartDatasource("ALERTS"))
	w.SetGroupBy(dashboards.ApiAggregateProperty("COUNTRY"))
	w.SetAxisGroupBy(dashboards.ApiAggregateProperty("TEST"))
	w.SetSortBy(dashboards.LegacyWidgetSortProperty("Alphabetical"))
	w.SetSortDirection(dashboards.LegacyWidgetSortDirection("Ascending"))
	w.SetLimit(12)
	w.SetShowLabels(true)
	w.SetIsHorizontalBarChart(false)

	data, err := mapGroupedBarChartWidget(dashboards.ApiGroupedBarchartWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Bar Chart: Grouped", data["type"])
	assert.Equal(t, "ALERTS", data["data_source"])

	config := data["grouped_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "COUNTRY", config["group_by"])
	assert.Equal(t, "TEST", config["axis_group_by"])
	assert.Equal(t, "Alphabetical", config["sort_by"])
	assert.Equal(t, "Ascending", config["sort_direction"])
	assert.Equal(t, 12, config["limit"])
	assert.Equal(t, true, config["show_labels"])
	assert.Equal(t, false, config["is_horizontal_bar_chart"])
}
