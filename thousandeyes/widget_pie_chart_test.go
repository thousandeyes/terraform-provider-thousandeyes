package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildPieChartWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic pie chart widget",
			input: map[string]interface{}{
				"type":        "Pie Chart",
				"title":       "Test Pie Chart",
				"visual_mode": "Full",
				"data_source": "ALERTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiPieChartWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Pie Chart", w.GetType())
				assert.Equal(t, "Test Pie Chart", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.PieChartDatasource("ALERTS"), w.GetDataSource())
			},
		},
		{
			name: "pie chart widget with config",
			input: map[string]interface{}{
				"type":         "Pie Chart",
				"title":        "Test Pie Chart",
				"visual_mode":  "Full",
				"data_source":  "ALERTS",
				"metric_group": "ALERTS",
				"metric":       "ALERT_COUNT_AGENT",
				"pie_chart_config": []interface{}{
					map[string]interface{}{
						"group_by": "AGENT",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiPieChartWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Test Pie Chart", w.GetTitle())
				assert.Equal(t, dashboards.MetricGroup("ALERTS"), w.GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("ALERT_COUNT_AGENT"), w.GetMetric())
				assert.Equal(t, dashboards.ApiAggregateProperty("AGENT"), w.GetGroupBy())
			},
		},
		{
			name: "pie chart widget with measure",
			input: map[string]interface{}{
				"type":        "Pie Chart",
				"title":       "Test Pie Chart",
				"visual_mode": "Full",
				"measure": []interface{}{
					map[string]interface{}{
						"type": "TOTAL",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiPieChartWidget
				assert.NotNil(t, w)
				measure := w.GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("TOTAL"), measure.GetType())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildPieChartWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapPieChartWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic pie chart widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiPieChartWidget("Pie Chart")
				w.SetId("widget-123")
				w.SetTitle("Test Pie Chart")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.PieChartDatasource("ALERTS"))
				return dashboards.ApiPieChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Pie Chart", data["type"])
				assert.Equal(t, "widget-123", data["id"])
				assert.Equal(t, "Test Pie Chart", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "ALERTS", data["data_source"])
			},
		},
		{
			name: "pie chart widget with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiPieChartWidget("Pie Chart")
				w.SetTitle("Test Pie Chart")
				w.SetGroupBy(dashboards.ApiAggregateProperty("AGENT"))
				return dashboards.ApiPieChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Pie Chart", data["type"])
				config := data["pie_chart_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, "AGENT", config["group_by"])
			},
		},
		{
			name: "pie chart widget with metric fields",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiPieChartWidget("Pie Chart")
				w.SetTitle("Test Pie Chart")
				w.SetMetricGroup(dashboards.MetricGroup("ALERTS"))
				w.SetMetric(dashboards.DashboardMetric("ALERT_COUNT_AGENT"))
				return dashboards.ApiPieChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "ALERTS", data["metric_group"])
				assert.Equal(t, "ALERT_COUNT_AGENT", data["metric"])
			},
		},
		{
			name: "pie chart widget with measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiPieChartWidget("Pie Chart")
				w.SetTitle("Test Pie Chart")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("TOTAL"))
				w.SetMeasure(*measure)
				return dashboards.ApiPieChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				measureList := data["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "TOTAL", measureMap["type"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapPieChartWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapPieChartWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapPieChartWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
