package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildTableWidget(t *testing.T) {
	widget := buildTableWidget(map[string]interface{}{
		"type":        "Table",
		"title":       "Test Table",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"table_config": []interface{}{
			map[string]interface{}{
				"compare_to_previous_value": true,
				"row_group_by":              "AGENT",
				"column_group_by":           "TEST",
				"limit":                     10,
			},
		},
	})

	w := widget.ApiTableWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Table", w.GetType())
	assert.Equal(t, dashboards.TableDatasource("ALERTS"), w.GetDataSource())
	assert.True(t, w.GetCompareToPreviousValue())
	assert.Equal(t, dashboards.ApiAggregateProperty("AGENT"), w.GetRowGroupBy())
	assert.Equal(t, dashboards.ApiAggregateProperty("TEST"), w.GetColumnGroupBy())
	assert.Equal(t, int32(10), w.GetLimit())
}

func TestMapTableWidget(t *testing.T) {
	w := dashboards.NewApiTableWidget("Table")
	w.SetId("widget-123")
	w.SetTitle("Test Table")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.TableDatasource("ALERTS"))
	w.SetCompareToPreviousValue(true)
	w.SetRowGroupBy(dashboards.ApiAggregateProperty("AGENT"))
	w.SetColumnGroupBy(dashboards.ApiAggregateProperty("TEST"))
	w.SetSortBy(dashboards.LegacyWidgetSortProperty("Alphabetical"))
	w.SetSortDirection(dashboards.LegacyWidgetSortDirection("Ascending"))
	w.SetLimit(10)

	data, err := mapTableWidget(dashboards.ApiTableWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Table", data["type"])
	assert.Equal(t, "widget-123", data["id"])
	assert.Equal(t, "ALERTS", data["data_source"])

	config := data["table_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, true, config["compare_to_previous_value"])
	assert.Equal(t, "AGENT", config["row_group_by"])
	assert.Equal(t, "TEST", config["column_group_by"])
	assert.NotContains(t, config, "sort_by")
	assert.NotContains(t, config, "sort_direction")
	assert.Equal(t, 10, config["limit"])
}
