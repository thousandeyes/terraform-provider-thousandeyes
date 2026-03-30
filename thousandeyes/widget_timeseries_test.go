package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildTimeseriesWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic timeseries widget",
			input: map[string]interface{}{
				"type":        "Time Series: Line",
				"title":       "Test Timeseries",
				"visual_mode": "Full",
				"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiTimeseriesWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Time Series: Line", w.GetType())
				assert.Equal(t, "Test Timeseries", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.TimeseriesDatasource("CLOUD_AND_ENTERPRISE_AGENTS"), w.GetDataSource())
			},
		},
		{
			name: "timeseries widget with config",
			input: map[string]interface{}{
				"type":         "Time Series: Line",
				"title":        "Test Timeseries",
				"visual_mode":  "Full",
				"data_source":  "CLOUD_AND_ENTERPRISE_AGENTS",
				"metric_group": "WEB_HTTP_SERVER",
				"metric":       "RESPONSE_TIME",
				"timeseries_config": []interface{}{
					map[string]interface{}{
						"min_scale":                        float64(0),
						"max_scale":                        float64(1000),
						"group_by":                         "AGENT",
						"show_timeseries_overall_baseline": true,
						"is_timeseries_one_chart_per_line": false,
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiTimeseriesWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Test Timeseries", w.GetTitle())
				assert.Equal(t, dashboards.MetricGroup("WEB_HTTP_SERVER"), w.GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("RESPONSE_TIME"), w.GetMetric())
				assert.Equal(t, float32(1000), w.GetMaxScale())
				assert.Equal(t, dashboards.ApiAggregateProperty("AGENT"), w.GetGroupBy())
				assert.True(t, w.GetShowTimeseriesOverallBaseline())
				assert.False(t, w.GetIsTimeseriesOneChartPerLine())
			},
		},
		{
			name: "timeseries widget with measure",
			input: map[string]interface{}{
				"type":        "Time Series: Line",
				"title":       "Test Timeseries",
				"visual_mode": "Full",
				"measure": []interface{}{
					map[string]interface{}{
						"type": "MEAN",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiTimeseriesWidget
				assert.NotNil(t, w)
				measure := w.GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("MEAN"), measure.GetType())
			},
		},
		{
			name: "timeseries widget with unit",
			input: map[string]interface{}{
				"type":        "Time Series: Line",
				"title":       "Test Timeseries",
				"visual_mode": "Full",
				"timeseries_config": []interface{}{
					map[string]interface{}{
						"unit": "ms",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiTimeseriesWidget
				assert.NotNil(t, w)
				assert.Equal(t, dashboards.ApiWidgetFixedYScalePrefix("ms"), w.GetUnit())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildTimeseriesWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapTimeseriesWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic timeseries widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
				w.SetId("widget-123")
				w.SetTitle("Test Timeseries")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.TimeseriesDatasource("CLOUD_AND_ENTERPRISE_AGENTS"))
				return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Time Series: Line", data["type"])
				assert.Equal(t, "widget-123", data["id"])
				assert.Equal(t, "Test Timeseries", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "CLOUD_AND_ENTERPRISE_AGENTS", data["data_source"])
			},
		},
		{
			name: "timeseries widget with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
				w.SetTitle("Test Timeseries")
				w.SetMinScale(0)
				w.SetMaxScale(1000)
				w.SetGroupBy(dashboards.ApiAggregateProperty("AGENT"))
				w.SetShowTimeseriesOverallBaseline(true)
				w.SetIsTimeseriesOneChartPerLine(false)
				return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Time Series: Line", data["type"])
				config := data["timeseries_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, float64(0), config["min_scale"])
				assert.Equal(t, float64(1000), config["max_scale"])
				assert.Equal(t, "AGENT", config["group_by"])
				assert.True(t, config["show_timeseries_overall_baseline"].(bool))
				assert.False(t, config["is_timeseries_one_chart_per_line"].(bool))
			},
		},
		{
			name: "timeseries widget with metric fields",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
				w.SetTitle("Test Timeseries")
				w.SetMetricGroup(dashboards.MetricGroup("WEB_HTTP_SERVER"))
				w.SetMetric(dashboards.DashboardMetric("RESPONSE_TIME"))
				return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "WEB_HTTP_SERVER", data["metric_group"])
				assert.Equal(t, "RESPONSE_TIME", data["metric"])
			},
		},
		{
			name: "timeseries widget with measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
				w.SetTitle("Test Timeseries")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("MEAN"))
				w.SetMeasure(*measure)
				return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				measureList := data["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "MEAN", measureMap["type"])
			},
		},
		{
			name: "timeseries widget with unit",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiTimeseriesWidget("Time Series: Line")
				w.SetTitle("Test Timeseries")
				w.SetUnit(dashboards.ApiWidgetFixedYScalePrefix("ms"))
				return dashboards.ApiTimeseriesWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				config := data["timeseries_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, "ms", config["unit"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapTimeseriesWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapTimeseriesWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapTimeseriesWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
