package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildMultiMetricTableWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic multi metric table widget",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Test Table",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Multi Metric Table", w.GetType())
				assert.Equal(t, "Test Table", w.GetTitle())
				cols, ok := w.GetMultiMetricColumnsOk()
				assert.True(t, ok)
				assert.Empty(t, cols)
			},
		},
		{
			name: "multi metric table with config",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Table With Config",
				"multi_metric_table_config": []interface{}{
					map[string]interface{}{
						"compare_to_previous_value": true,
						"row_group_by":              "TESTS",
						"limit":                     10,
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				assert.True(t, w.GetCompareToPreviousValue())
				assert.Equal(t, dashboards.ApiAggregateProperty("TESTS"), w.GetRowGroupBy())
				assert.Equal(t, int32(10), w.GetLimit())
			},
		},
		{
			name: "multi metric table with columns",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Table With Columns",
				"multi_metric_columns": []interface{}{
					map[string]interface{}{
						"data_source":  "ALERTS",
						"metric_group": "ALERTS",
						"metric":       "ALERT_COUNT_AGENT",
						"measure": []interface{}{
							map[string]interface{}{
								"type": "MEAN",
							},
						},
					},
					map[string]interface{}{
						"data_source":  "ALERTS",
						"metric_group": "ALERTS",
						"metric":       "ACTIVE_ALERT_COUNT",
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				cols := w.GetMultiMetricColumns()
				assert.Len(t, cols, 2)

				assert.Equal(t, dashboards.MultiMetricsTableDatasource("ALERTS"), cols[0].GetDataSource())
				assert.Equal(t, dashboards.MetricGroup("ALERTS"), cols[0].GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("ALERT_COUNT_AGENT"), cols[0].GetMetric())
				measure := cols[0].GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("MEAN"), measure.GetType())

				assert.Equal(t, dashboards.DashboardMetric("ACTIVE_ALERT_COUNT"), cols[1].GetMetric())
			},
		},
		{
			name: "column with filters",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Filtered Table",
				"multi_metric_columns": []interface{}{
					map[string]interface{}{
						"data_source":  "ALERTS",
						"metric_group": "ALERTS",
						"metric":       "ALERT_COUNT_AGENT",
						"filter": testFilterSet(map[string]interface{}{
							"property": "TEST",
							"values":   []interface{}{"12345"},
						}),
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				cols := w.GetMultiMetricColumns()
				assert.Len(t, cols, 1)
				filters := cols[0].GetFilters()
				assert.Contains(t, filters, "TEST")
				assert.Equal(t, []interface{}{"12345"}, filters["TEST"])
			},
		},
		{
			name: "column with percentile measure",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Percentile Table",
				"multi_metric_columns": []interface{}{
					map[string]interface{}{
						"data_source":  "ALERTS",
						"metric_group": "ALERTS",
						"metric":       "ALERT_COUNT_AGENT",
						"measure": []interface{}{
							map[string]interface{}{
								"type":             "NTH_PERCENTILE",
								"percentile_value": 95.0,
							},
						},
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				cols := w.GetMultiMetricColumns()
				assert.Len(t, cols, 1)
				measure := cols[0].GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("NTH_PERCENTILE"), measure.GetType())
				assert.Equal(t, float32(95.0), measure.GetPercentileValue())
			},
		},
		{
			name: "removing all columns sends empty array",
			input: map[string]interface{}{
				"type":  "Multi Metric Table",
				"title": "Empty Columns",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiMultiMetricTableWidget
				assert.NotNil(t, w)
				cols, ok := w.GetMultiMetricColumnsOk()
				assert.True(t, ok, "MultiMetricColumns should be explicitly set (not nil)")
				assert.Empty(t, cols, "MultiMetricColumns should be an empty slice")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildMultiMetricTableWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapMultiMetricTableWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic multi metric table widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
				w.SetId("widget-mmt-1")
				w.SetTitle("Test Table")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				return dashboards.ApiMultiMetricTableWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Multi Metric Table", data["type"])
				assert.Equal(t, "widget-mmt-1", data["id"])
				assert.Equal(t, "Test Table", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
			},
		},
		{
			name: "multi metric table with config",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
				w.SetTitle("Config Table")
				w.SetCompareToPreviousValue(true)
				w.SetRowGroupBy(dashboards.ApiAggregateProperty("TESTS"))
				w.SetLimit(20)
				return dashboards.ApiMultiMetricTableWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				config := data["multi_metric_table_config"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, true, config["compare_to_previous_value"])
				assert.Equal(t, "TESTS", config["row_group_by"])
				assert.Equal(t, 20, config["limit"])
			},
		},
		{
			name: "multi metric table with columns",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
				w.SetTitle("Columns Table")
				col1 := dashboards.NewApiMultiMetricColumn()
				col1.SetId("col-1")
				col1.SetDataSource(dashboards.MultiMetricsTableDatasource("ALERTS"))
				col1.SetMetricGroup(dashboards.MetricGroup("ALERTS"))
				col1.SetMetric(dashboards.DashboardMetric("ALERT_COUNT_AGENT"))
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("MEAN"))
				col1.SetMeasure(*measure)

				col2 := dashboards.NewApiMultiMetricColumn()
				col2.SetId("col-2")
				col2.SetMetric(dashboards.DashboardMetric("ACTIVE_ALERT_COUNT"))

				w.SetMultiMetricColumns([]dashboards.ApiMultiMetricColumn{*col1, *col2})
				return dashboards.ApiMultiMetricTableWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				cols := data["multi_metric_columns"].([]interface{})
				assert.Len(t, cols, 2)

				col1 := cols[0].(map[string]interface{})
				assert.Equal(t, "col-1", col1["id"])
				assert.Equal(t, "ALERTS", col1["data_source"])
				assert.Equal(t, "ALERTS", col1["metric_group"])
				assert.Equal(t, "ALERT_COUNT_AGENT", col1["metric"])
				measureList := col1["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "MEAN", measureMap["type"])

				col2 := cols[1].(map[string]interface{})
				assert.Equal(t, "col-2", col2["id"])
				assert.Equal(t, "ACTIVE_ALERT_COUNT", col2["metric"])
			},
		},
		{
			name: "column with filters",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
				w.SetTitle("Filtered Table")
				col := dashboards.NewApiMultiMetricColumn()
				col.SetId("col-f")
				col.SetFilters(map[string][]interface{}{
					"TEST": {"12345"},
				})
				w.SetMultiMetricColumns([]dashboards.ApiMultiMetricColumn{*col})
				return dashboards.ApiMultiMetricTableWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				cols := data["multi_metric_columns"].([]interface{})
				assert.Len(t, cols, 1)
				col := cols[0].(map[string]interface{})
				filters := col["filter"].([]interface{})
				assert.Len(t, filters, 1)
				f := filters[0].(map[string]interface{})
				assert.Equal(t, "TEST", f["property"])
				assert.Equal(t, []interface{}{"12345"}, f["values"])
			},
		},
		{
			name: "column with percentile measure",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiMultiMetricTableWidget("Multi Metric Table")
				w.SetTitle("Percentile Table")
				col := dashboards.NewApiMultiMetricColumn()
				col.SetId("col-p")
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("NTH_PERCENTILE"))
				measure.SetPercentileValue(99.0)
				col.SetMeasure(*measure)
				w.SetMultiMetricColumns([]dashboards.ApiMultiMetricColumn{*col})
				return dashboards.ApiMultiMetricTableWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				cols := data["multi_metric_columns"].([]interface{})
				col := cols[0].(map[string]interface{})
				measureList := col["measure"].([]interface{})
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "NTH_PERCENTILE", measureMap["type"])
				assert.Equal(t, float64(99.0), measureMap["percentile_value"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapMultiMetricTableWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapMultiMetricTableWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapMultiMetricTableWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestMultiMetricTableWidgetRegistration(t *testing.T) {
	reg, ok := widgetRegistry[WidgetTypeMultiMetricTable]
	assert.True(t, ok, "Multi Metric Table should be registered")
	assert.NotNil(t, reg.Builder)
	assert.NotNil(t, reg.Mapper)
}
