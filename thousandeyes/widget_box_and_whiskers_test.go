package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildBoxAndWhiskersWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic box and whiskers widget",
			input: map[string]interface{}{
				"type":        "Box and Whiskers",
				"title":       "Test Box and Whiskers",
				"visual_mode": "Full",
				"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiBoxAndWhiskersWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Box and Whiskers", w.GetType())
				assert.Equal(t, "Test Box and Whiskers", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.BoxAndWhiskersDatasource("CLOUD_AND_ENTERPRISE_AGENTS"), w.GetDataSource())
			},
		},
		{
			name: "box and whiskers widget with config",
			input: map[string]interface{}{
				"type":         "Box and Whiskers",
				"title":        "Test Box and Whiskers",
				"visual_mode":  "Full",
				"data_source":  "CLOUD_AND_ENTERPRISE_AGENTS",
				"metric_group": "WEB_HTTP_SERVER",
				"metric":       "RESPONSE_TIME",
				"box_and_whiskers_config": []interface{}{
					map[string]interface{}{
						"min_scale": float64(0),
						"max_scale": float64(1000),
						"group_by":  "AGENT",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiBoxAndWhiskersWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Test Box and Whiskers", w.GetTitle())
				assert.Equal(t, dashboards.MetricGroup("WEB_HTTP_SERVER"), w.GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("RESPONSE_TIME"), w.GetMetric())
				assert.Equal(t, float32(0), w.GetMinScale())
				assert.Equal(t, float32(1000), w.GetMaxScale())
				assert.Equal(t, dashboards.ApiAggregateProperty("AGENT"), w.GetGroupBy())
			},
		},
		{
			name: "box and whiskers widget with measure",
			input: map[string]interface{}{
				"type":        "Box and Whiskers",
				"title":       "Test Box and Whiskers",
				"visual_mode": "Full",
				"measure": []interface{}{
					map[string]interface{}{
						"type": "MEAN",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiBoxAndWhiskersWidget
				assert.NotNil(t, w)
				measure := w.GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("MEAN"), measure.GetType())
			},
		},
		{
			name: "box and whiskers widget with unit",
			input: map[string]interface{}{
				"type":        "Box and Whiskers",
				"title":       "Test Box and Whiskers",
				"visual_mode": "Full",
				"box_and_whiskers_config": []interface{}{
					map[string]interface{}{
						"unit": "ms",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiBoxAndWhiskersWidget
				assert.NotNil(t, w)
				assert.Equal(t, dashboards.ApiWidgetFixedYScalePrefix("ms"), w.GetUnit())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildBoxAndWhiskersWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapBoxAndWhiskersWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic box and whiskers widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
				w.SetId("widget-123")
				w.SetTitle("Test Box and Whiskers")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.BoxAndWhiskersDatasource("CLOUD_AND_ENTERPRISE_AGENTS"))
				return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Box and Whiskers", data["type"])
				assert.Equal(t, "widget-123", data["id"])
				assert.Equal(t, "Test Box and Whiskers", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "CLOUD_AND_ENTERPRISE_AGENTS", data["data_source"])
			},
		},
		{
			name: "box and whiskers widget with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
				w.SetTitle("Test Box and Whiskers")
				w.SetMinScale(0)
				w.SetMaxScale(1000)
				w.SetGroupBy(dashboards.ApiAggregateProperty("AGENT"))
				return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Box and Whiskers", data["type"])
				config := data["box_and_whiskers_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, float64(0), config["min_scale"])
				assert.Equal(t, float64(1000), config["max_scale"])
				assert.Equal(t, "AGENT", config["group_by"])
			},
		},
		{
			name: "box and whiskers widget with metric fields",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
				w.SetTitle("Test Box and Whiskers")
				w.SetMetricGroup(dashboards.MetricGroup("WEB_HTTP_SERVER"))
				w.SetMetric(dashboards.DashboardMetric("RESPONSE_TIME"))
				return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "WEB_HTTP_SERVER", data["metric_group"])
				assert.Equal(t, "RESPONSE_TIME", data["metric"])
			},
		},
		{
			name: "box and whiskers widget with measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
				w.SetTitle("Test Box and Whiskers")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("MEAN"))
				w.SetMeasure(*measure)
				return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				measureList := data["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "MEAN", measureMap["type"])
			},
		},
		{
			name: "box and whiskers widget with unit",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiBoxAndWhiskersWidget("Box and Whiskers")
				w.SetTitle("Test Box and Whiskers")
				w.SetUnit(dashboards.ApiWidgetFixedYScalePrefix("ms"))
				return dashboards.ApiBoxAndWhiskersWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				config := data["box_and_whiskers_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, "ms", config["unit"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapBoxAndWhiskersWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapBoxAndWhiskersWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapBoxAndWhiskersWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
