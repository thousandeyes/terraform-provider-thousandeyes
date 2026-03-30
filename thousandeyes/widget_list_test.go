package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildListWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic list widget",
			input: map[string]interface{}{
				"type":        "List",
				"title":       "Test List",
				"visual_mode": "Full",
				"data_source": "ALERTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiListWidget
				assert.NotNil(t, w)
				assert.Equal(t, "List", w.GetType())
				assert.Equal(t, "Test List", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.ListDatasource("ALERTS"), w.GetDataSource())
			},
		},
		{
			name: "list widget with config",
			input: map[string]interface{}{
				"type":         "List",
				"title":        "Test List",
				"visual_mode":  "Full",
				"data_source":  "ALERTS",
				"metric_group": "ALERTS",
				"metric":       "ALERT_COUNT_AGENT",
				"list_config": []interface{}{
					map[string]interface{}{
						"active_within_value": 7,
						"active_within_unit":  "Days",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiListWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Test List", w.GetTitle())
				assert.Equal(t, dashboards.MetricGroup("ALERTS"), w.GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("ALERT_COUNT_AGENT"), w.GetMetric())
				activeWithin := w.GetActiveWithin()
				assert.Equal(t, int32(7), activeWithin.GetValue())
				assert.Equal(t, dashboards.LegacyDurationUnit("Days"), activeWithin.GetUnit())
			},
		},
		{
			name: "list widget with measure",
			input: map[string]interface{}{
				"type":        "List",
				"title":       "Test List",
				"visual_mode": "Full",
				"measure": []interface{}{
					map[string]interface{}{
						"type": "TOTAL",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiListWidget
				assert.NotNil(t, w)
				measure := w.GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("TOTAL"), measure.GetType())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildListWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapListWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic list widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiListWidget("List")
				w.SetId("widget-123")
				w.SetTitle("Test List")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.ListDatasource("ALERTS"))
				return dashboards.ApiListWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "List", data["type"])
				assert.Equal(t, "widget-123", data["id"])
				assert.Equal(t, "Test List", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "ALERTS", data["data_source"])
			},
		},
		{
			name: "list widget with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiListWidget("List")
				w.SetTitle("Test List")
				activeWithin := dashboards.NewActiveWithin()
				activeWithin.SetValue(7)
				activeWithin.SetUnit(dashboards.LegacyDurationUnit("Days"))
				w.SetActiveWithin(*activeWithin)
				return dashboards.ApiListWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "List", data["type"])
				config := data["list_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, 7, config["active_within_value"])
				assert.Equal(t, "Days", config["active_within_unit"])
			},
		},
		{
			name: "list widget with metric fields",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiListWidget("List")
				w.SetTitle("Test List")
				w.SetMetricGroup(dashboards.MetricGroup("ALERTS"))
				w.SetMetric(dashboards.DashboardMetric("ALERT_COUNT_AGENT"))
				return dashboards.ApiListWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "ALERTS", data["metric_group"])
				assert.Equal(t, "ALERT_COUNT_AGENT", data["metric"])
			},
		},
		{
			name: "list widget with measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiListWidget("List")
				w.SetTitle("Test List")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("TOTAL"))
				w.SetMeasure(*measure)
				return dashboards.ApiListWidgetAsApiWidget(w)
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
			data, err := mapListWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapListWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapListWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}
