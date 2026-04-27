package thousandeyes

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeConfiguredWidgets_NumberCardScales(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type":  "Number",
			"title": "Number",
			"number_cards": []interface{}{
				map[string]interface{}{
					"description":  "omitted scales",
					"data_source":  "ALERTS",
					"metric_group": "ALERTS",
					"metric":       "ALERT_COUNT_AGENT",
					"min_scale":    0.0,
					"max_scale":    0.0,
				},
				map[string]interface{}{
					"description":  "explicit zero scales",
					"data_source":  "ALERTS",
					"metric_group": "ALERTS",
					"metric":       "ALERT_COUNT_AGENT",
					"min_scale":    0.0,
					"max_scale":    0.0,
				},
			},
		},
	}

	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.TupleVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Number"),
				"number_cards": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"description":  cty.StringVal("omitted scales"),
						"data_source":  cty.StringVal("ALERTS"),
						"metric_group": cty.StringVal("ALERTS"),
						"metric":       cty.StringVal("ALERT_COUNT_AGENT"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"description":  cty.StringVal("explicit zero scales"),
						"data_source":  cty.StringVal("ALERTS"),
						"metric_group": cty.StringVal("ALERTS"),
						"metric":       cty.StringVal("ALERT_COUNT_AGENT"),
						"min_scale":    cty.NumberIntVal(0),
						"max_scale":    cty.NumberIntVal(0),
					}),
				}),
			}),
		}),
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)
	cards := normalized[0].(map[string]interface{})["number_cards"].([]interface{})
	omitted := cards[0].(map[string]interface{})
	explicitZero := cards[1].(map[string]interface{})

	assert.NotContains(t, omitted, "min_scale")
	assert.NotContains(t, omitted, "max_scale")
	assert.Contains(t, explicitZero, "min_scale")
	assert.Contains(t, explicitZero, "max_scale")

	widgets, err := BuildWidgets(normalized)
	require.NoError(t, err)
	require.Len(t, widgets, 1)

	numberCards := widgets[0].ApiNumbersCardWidget.GetNumberCards()
	require.Len(t, numberCards, 2)

	_, minScaleOmitted := numberCards[0].GetMinScaleOk()
	_, maxScaleOmitted := numberCards[0].GetMaxScaleOk()
	assert.False(t, minScaleOmitted)
	assert.False(t, maxScaleOmitted)

	minScale, minScaleExplicit := numberCards[1].GetMinScaleOk()
	maxScale, maxScaleExplicit := numberCards[1].GetMaxScaleOk()
	require.True(t, minScaleExplicit)
	require.True(t, maxScaleExplicit)
	assert.Equal(t, float32(0), *minScale)
	assert.Equal(t, float32(0), *maxScale)
}

func TestNormalizeConfiguredWidgets_WidgetConfigScales(t *testing.T) {
	tests := []struct {
		name       string
		widgetType string
		blockName  string
		rawExtra   map[string]cty.Value
	}{
		{name: "map", widgetType: "Map", blockName: "geo_map_config"},
		{name: "time series line", widgetType: "Time Series: Line", blockName: "timeseries_config"},
		{
			name:       "stacked area",
			widgetType: "Time Series: Stacked Area",
			blockName:  "stacked_area_config",
			rawExtra: map[string]cty.Value{
				"group_by": cty.StringVal("TEST"),
			},
		},
		{name: "box and whiskers", widgetType: "Box and Whiskers", blockName: "box_and_whiskers_config"},
		{name: "color grid", widgetType: "Color Grid", blockName: "color_grid_config"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widgetList := []interface{}{
				map[string]interface{}{
					"type": tc.widgetType,
					tc.blockName: []interface{}{
						map[string]interface{}{
							"min_scale": 0.0,
							"max_scale": 0.0,
						},
					},
				},
			}

			rawBlock := map[string]cty.Value{}
			for key, value := range tc.rawExtra {
				rawBlock[key] = value
			}
			rawConfig := cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"type":       cty.StringVal(tc.widgetType),
						tc.blockName: cty.TupleVal([]cty.Value{cty.ObjectVal(rawBlock)}),
					}),
				}),
			})

			normalized := normalizeConfiguredWidgets(widgetList, rawConfig)
			block := normalized[0].(map[string]interface{})[tc.blockName].([]interface{})[0].(map[string]interface{})

			assert.NotContains(t, block, "min_scale")
			assert.NotContains(t, block, "max_scale")
		})
	}
}

func TestNormalizeConfiguredWidgets_WidgetConfigExplicitZeroScales(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type": "Color Grid",
			"color_grid_config": []interface{}{
				map[string]interface{}{
					"min_scale": 0.0,
					"max_scale": 0.0,
				},
			},
		},
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.TupleVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Color Grid"),
				"color_grid_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"min_scale": cty.NumberIntVal(0),
						"max_scale": cty.NumberIntVal(0),
					}),
				}),
			}),
		}),
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)
	block := normalized[0].(map[string]interface{})["color_grid_config"].([]interface{})[0].(map[string]interface{})

	assert.Contains(t, block, "min_scale")
	assert.Contains(t, block, "max_scale")
}

func TestNormalizeConfiguredWidgets_FallbackKeepsOriginalShape(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type": "Number",
			"number_cards": []interface{}{
				map[string]interface{}{
					"min_scale": 0.0,
					"max_scale": 0.0,
				},
			},
		},
	}

	normalized := normalizeConfiguredWidgets(widgetList, cty.NilVal)

	assert.Equal(t, widgetList, normalized)
}
