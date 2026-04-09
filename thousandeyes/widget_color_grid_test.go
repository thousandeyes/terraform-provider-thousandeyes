package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildColorGridWidget(t *testing.T) {
	widget := buildColorGridWidget(map[string]interface{}{
		"type":        "Color Grid",
		"title":       "Color Grid",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"color_grid_config": []interface{}{
			map[string]interface{}{
				"min_scale":            float64(0),
				"max_scale":            float64(100),
				"unit":                 "Mbps",
				"cards":                "COUNTRY",
				"group_cards_by":       "TEST",
				"columns":              2,
				"limit":                6,
				"sort_by":              "Value",
				"sort_direction":       "Descending",
				"sort_group_by":        "Alphabetical",
				"sort_group_direction": "Ascending",
			},
		},
	})

	w := widget.ApiColorGridWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Color Grid", w.GetType())
	assert.Equal(t, dashboards.ColorGridDatasource("ALERTS"), w.GetDataSource())
	assert.Equal(t, float32(0), w.GetMinScale())
	assert.Equal(t, float32(100), w.GetMaxScale())
	assert.Equal(t, dashboards.ApiWidgetFixedYScalePrefix("Mbps"), w.GetUnit())
	assert.Equal(t, dashboards.ApiAggregateProperty("COUNTRY"), w.GetCards())
	assert.Equal(t, dashboards.ApiAggregateProperty("TEST"), w.GetGroupCardsBy())
	assert.Equal(t, int32(2), w.GetColumns())
	assert.Equal(t, int32(6), w.GetLimit())
	assert.Equal(t, dashboards.LegacyWidgetSortProperty("Value"), w.GetSortBy())
	assert.Equal(t, dashboards.LegacyWidgetSortDirection("Descending"), w.GetSortDirection())
	assert.Equal(t, dashboards.LegacyWidgetSortProperty("Alphabetical"), w.GetSortGroupBy())
	assert.Equal(t, dashboards.LegacyWidgetSortDirection("Ascending"), w.GetSortGroupDirection())
}

func TestMapColorGridWidget(t *testing.T) {
	w := dashboards.NewApiColorGridWidget("Color Grid")
	w.SetId("widget-123")
	w.SetTitle("Color Grid")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.ColorGridDatasource("ALERTS"))
	w.SetMinScale(0)
	w.SetMaxScale(100)
	w.SetUnit(dashboards.ApiWidgetFixedYScalePrefix("Mbps"))
	w.SetCards(dashboards.ApiAggregateProperty("COUNTRY"))
	w.SetGroupCardsBy(dashboards.ApiAggregateProperty("TEST"))
	w.SetColumns(2)
	w.SetLimit(6)
	w.SetSortBy(dashboards.LegacyWidgetSortProperty("Value"))
	w.SetSortDirection(dashboards.LegacyWidgetSortDirection("Descending"))
	w.SetSortGroupBy(dashboards.LegacyWidgetSortProperty("Alphabetical"))
	w.SetSortGroupDirection(dashboards.LegacyWidgetSortDirection("Ascending"))

	data, err := mapColorGridWidget(dashboards.ApiColorGridWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Color Grid", data["type"])
	assert.Equal(t, "ALERTS", data["data_source"])

	config := data["color_grid_config"].([]interface{})[0].(map[string]interface{})
	assert.Equal(t, float64(0), config["min_scale"])
	assert.Equal(t, float64(100), config["max_scale"])
	assert.Equal(t, "Mbps", config["unit"])
	assert.Equal(t, "COUNTRY", config["cards"])
	assert.Equal(t, "TEST", config["group_cards_by"])
	assert.Equal(t, 2, config["columns"])
	assert.Equal(t, 6, config["limit"])
	assert.Equal(t, "Value", config["sort_by"])
	assert.Equal(t, "Descending", config["sort_direction"])
	assert.Equal(t, "Alphabetical", config["sort_group_by"])
	assert.Equal(t, "Ascending", config["sort_group_direction"])
}
