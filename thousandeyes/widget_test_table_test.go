package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildTestTableWidget(t *testing.T) {
	widget := buildTestTableWidget(map[string]interface{}{
		"type":        "Test Table",
		"title":       "Test Table Widget",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"test_table_config": []interface{}{
			map[string]interface{}{
				"filter": []interface{}{
					map[string]interface{}{
						"type": "all",
						"filters": []interface{}{
							map[string]interface{}{"key": "Test Name", "value": "API"},
						},
					},
				},
				"exclude": []interface{}{
					map[string]interface{}{
						"type": "any",
						"filters": []interface{}{
							map[string]interface{}{"key": "Label ID", "value": "123"},
						},
					},
				},
			},
		},
	})

	w := widget.ApiTestTableWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Test Table", w.GetType())
	assert.Equal(t, dashboards.TestTableDatasource("ALERTS"), w.GetDataSource())
	filter := w.GetFilter()
	assert.Equal(t, dashboards.TestTableFilterType("all"), (&filter).GetType())
	assert.Len(t, (&filter).GetFilters(), 1)
	assert.Equal(t, dashboards.TestTableFilterKey("Test Name"), (&filter).GetFilters()[0].GetKey())
	assert.Equal(t, "API", (&filter).GetFilters()[0].GetValue())
	exclude := w.GetExclude()
	assert.Equal(t, dashboards.TestTableFilterType("any"), (&exclude).GetType())
	assert.Len(t, (&exclude).GetFilters(), 1)
	assert.Equal(t, dashboards.TestTableFilterKey("Label ID"), (&exclude).GetFilters()[0].GetKey())
	assert.Equal(t, "123", (&exclude).GetFilters()[0].GetValue())
}

func TestMapTestTableWidget(t *testing.T) {
	w := dashboards.NewApiTestTableWidget("Test Table")
	w.SetId("widget-123")
	w.SetTitle("Test Table Widget")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.TestTableDatasource("ALERTS"))

	filter := dashboards.NewApiWidgetFilterApiTestTableFilterKey()
	filter.SetType(dashboards.TestTableFilterType("all"))
	filter.SetFilters([]dashboards.ApiMultiSearchFilterApiTestTableFilterKey{
		func() dashboards.ApiMultiSearchFilterApiTestTableFilterKey {
			item := dashboards.NewApiMultiSearchFilterApiTestTableFilterKey()
			item.SetKey(dashboards.TestTableFilterKey("Test Name"))
			item.SetValue("API")
			return *item
		}(),
	})
	w.SetFilter(*filter)

	exclude := dashboards.NewApiWidgetFilterApiTestTableFilterKey()
	exclude.SetType(dashboards.TestTableFilterType("any"))
	exclude.SetFilters([]dashboards.ApiMultiSearchFilterApiTestTableFilterKey{
		func() dashboards.ApiMultiSearchFilterApiTestTableFilterKey {
			item := dashboards.NewApiMultiSearchFilterApiTestTableFilterKey()
			item.SetKey(dashboards.TestTableFilterKey("Label ID"))
			item.SetValue("123")
			return *item
		}(),
	})
	w.SetExclude(*exclude)

	data, err := mapTestTableWidget(dashboards.ApiTestTableWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Test Table", data["type"])
	assert.Equal(t, "widget-123", data["id"])
	assert.Equal(t, "ALERTS", data["data_source"])

	config := data["test_table_config"].([]interface{})[0].(map[string]interface{})
	filterBlock := config["filter"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "all", filterBlock["type"])
	assert.Equal(t, []interface{}{map[string]interface{}{"key": "Test Name", "value": "API"}}, filterBlock["filters"])
	excludeBlock := config["exclude"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, "any", excludeBlock["type"])
	assert.Equal(t, []interface{}{map[string]interface{}{"key": "Label ID", "value": "123"}}, excludeBlock["filters"])
}
