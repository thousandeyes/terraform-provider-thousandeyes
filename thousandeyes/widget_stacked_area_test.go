package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildStackedAreaWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic stacked area widget",
			input: map[string]interface{}{
				"type":        "Time Series: Stacked Area",
				"title":       "Test Stacked Area",
				"visual_mode": "Full",
				"data_source": "ALERTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiStackedAreaChartWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Time Series: Stacked Area", w.GetType())
				assert.Equal(t, "Test Stacked Area", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.StackedAreaChartDatasource("ALERTS"), w.GetDataSource())
			},
		},
		{
			name: "stacked area widget with config",
			input: map[string]interface{}{
				"type":         "Time Series: Stacked Area",
				"title":        "Test Stacked Area",
				"visual_mode":  "Full",
				"data_source":  "ALERTS",
				"metric_group": "ALERTS",
				"metric":       "ALERT_COUNT_AGENT",
				"stacked_area_config": []interface{}{
					map[string]interface{}{
						"min_scale": float64(0),
						"max_scale": float64(100),
						"group_by":  "AGENT",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiStackedAreaChartWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Test Stacked Area", w.GetTitle())
				assert.Equal(t, dashboards.MetricGroup("ALERTS"), w.GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("ALERT_COUNT_AGENT"), w.GetMetric())
				assert.Equal(t, float32(0), w.GetMinScale())
				assert.Equal(t, float32(100), w.GetMaxScale())
				assert.Equal(t, dashboards.ApiAggregateProperty("AGENT"), w.GetGroupBy())
			},
		},
		{
			name: "stacked area widget with measure",
			input: map[string]interface{}{
				"type":        "Time Series: Stacked Area",
				"title":       "Test Stacked Area",
				"visual_mode": "Full",
				"measure": []interface{}{
					map[string]interface{}{
						"type": "TOTAL",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiStackedAreaChartWidget
				assert.NotNil(t, w)
				measure := w.GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("TOTAL"), measure.GetType())
			},
		},
		{
			name: "stacked area widget with unit",
			input: map[string]interface{}{
				"type":        "Time Series: Stacked Area",
				"title":       "Test Stacked Area",
				"visual_mode": "Full",
				"stacked_area_config": []interface{}{
					map[string]interface{}{
						"unit": "ms",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiStackedAreaChartWidget
				assert.NotNil(t, w)
				assert.Equal(t, dashboards.ApiWidgetFixedYScalePrefix("ms"), w.GetUnit())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildStackedAreaWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapStackedAreaWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic stacked area widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
				w.SetId("widget-123")
				w.SetTitle("Test Stacked Area")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.StackedAreaChartDatasource("ALERTS"))
				return dashboards.ApiStackedAreaChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Time Series: Stacked Area", data["type"])
				assert.Equal(t, "widget-123", data["id"])
				assert.Equal(t, "Test Stacked Area", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "ALERTS", data["data_source"])
			},
		},
		{
			name: "stacked area widget with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
				w.SetTitle("Test Stacked Area")
				w.SetMinScale(0)
				w.SetMaxScale(100)
				w.SetGroupBy(dashboards.ApiAggregateProperty("AGENT"))
				return dashboards.ApiStackedAreaChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Time Series: Stacked Area", data["type"])
				config := data["stacked_area_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, float64(0), config["min_scale"])
				assert.Equal(t, float64(100), config["max_scale"])
				assert.Equal(t, "AGENT", config["group_by"])
			},
		},
		{
			name: "stacked area widget with metric fields",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
				w.SetTitle("Test Stacked Area")
				w.SetMetricGroup(dashboards.MetricGroup("ALERTS"))
				w.SetMetric(dashboards.DashboardMetric("ALERT_COUNT_AGENT"))
				return dashboards.ApiStackedAreaChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "ALERTS", data["metric_group"])
				assert.Equal(t, "ALERT_COUNT_AGENT", data["metric"])
			},
		},
		{
			name: "stacked area widget with measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
				w.SetTitle("Test Stacked Area")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("TOTAL"))
				w.SetMeasure(*measure)
				return dashboards.ApiStackedAreaChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				measureList := data["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "TOTAL", measureMap["type"])
			},
		},
		{
			name: "stacked area widget with unit",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiStackedAreaChartWidget("Time Series: Stacked Area")
				w.SetTitle("Test Stacked Area")
				w.SetUnit(dashboards.ApiWidgetFixedYScalePrefix("ms"))
				return dashboards.ApiStackedAreaChartWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				config := data["stacked_area_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, "ms", config["unit"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapStackedAreaWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapStackedAreaWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapStackedAreaWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
