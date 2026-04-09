package thousandeyes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildAlertListWidget(t *testing.T) {
	widget := buildAlertListWidget(map[string]interface{}{
		"type":        "Alert List",
		"title":       "Test Alert List",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"alert_list_config": []interface{}{
			map[string]interface{}{
				"alert_types": schema.NewSet(schema.HashString, []interface{}{
					"API",
					"DNS Server",
				}),
				"limit_to":            15,
				"active_within_value": 7,
				"active_within_unit":  "Days",
			},
		},
	})

	w := widget.ApiAlertListWidget
	assert.NotNil(t, w)
	assert.Equal(t, "Alert List", w.GetType())
	assert.Equal(t, dashboards.AlertListDatasource("ALERTS"), w.GetDataSource())
	assert.ElementsMatch(t, []dashboards.LegacyAlertListAlertType{
		dashboards.LegacyAlertListAlertType("API"),
		dashboards.LegacyAlertListAlertType("DNS Server"),
	}, w.GetAlertTypes())
	assert.Equal(t, int32(15), w.GetLimitTo())
	activeWithin := w.GetActiveWithin()
	assert.Equal(t, int32(7), (&activeWithin).GetValue())
	assert.Equal(t, dashboards.LegacyDurationUnit("Days"), (&activeWithin).GetUnit())
}

func TestBuildAlertListWidgetPreservesZeroLimit(t *testing.T) {
	widget := buildAlertListWidget(map[string]interface{}{
		"type":        "Alert List",
		"title":       "Zero Limit Alert List",
		"visual_mode": "Full",
		"data_source": "ALERTS",
		"alert_list_config": []interface{}{
			map[string]interface{}{
				"limit_to": 0,
			},
		},
	})

	w := widget.ApiAlertListWidget
	assert.NotNil(t, w)

	limit, ok := w.GetLimitToOk()
	assert.True(t, ok)
	assert.NotNil(t, limit)
	assert.Equal(t, int32(0), *limit)
}

func TestMapAlertListWidget(t *testing.T) {
	w := dashboards.NewApiAlertListWidget("Alert List")
	w.SetId("widget-123")
	w.SetTitle("Test Alert List")
	w.SetVisualMode(dashboards.VisualMode("Full"))
	w.SetDataSource(dashboards.AlertListDatasource("ALERTS"))
	w.SetAlertTypes([]dashboards.LegacyAlertListAlertType{
		dashboards.LegacyAlertListAlertType("API"),
		dashboards.LegacyAlertListAlertType("DNS Server"),
	})
	w.SetLimitTo(15)
	activeWithin := dashboards.NewActiveWithin()
	activeWithin.SetValue(7)
	activeWithin.SetUnit(dashboards.LegacyDurationUnit("Days"))
	w.SetActiveWithin(*activeWithin)

	data, err := mapAlertListWidget(dashboards.ApiAlertListWidgetAsApiWidget(w))
	assert.NoError(t, err)
	assert.Equal(t, "Alert List", data["type"])
	assert.Equal(t, "widget-123", data["id"])
	assert.Equal(t, "ALERTS", data["data_source"])

	config := data["alert_list_config"].([]interface{})[0].(map[string]interface{})
	assert.ElementsMatch(t, []interface{}{"API", "DNS Server"}, config["alert_types"].([]interface{}))
	assert.Equal(t, 15, config["limit_to"])
	assert.Equal(t, 7, config["active_within_value"])
	assert.Equal(t, "Days", config["active_within_unit"])
}
