package thousandeyes

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeConfiguredWidgets_FieldConfiguration(t *testing.T) {
	assert.Contains(t, dashboardWidgetPresenceSensitiveTopLevelFields, "should_exclude_alert_suppression_windows")
	assert.Contains(t, dashboardWidgetPresenceSensitiveConfigBlocks, "number_cards")
	assert.ElementsMatch(t,
		[]string{"min_scale", "max_scale", "compare_to_previous_value", "should_exclude_alert_suppression_windows"},
		dashboardWidgetPresenceSensitiveConfigBlocks["number_cards"],
	)
}

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

func TestNormalizeConfiguredWidgets_WidgetConfigUnknownScaleIsConfigured(t *testing.T) {
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
						"min_scale": cty.UnknownVal(cty.Number),
					}),
				}),
			}),
		}),
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)
	block := normalized[0].(map[string]interface{})["color_grid_config"].([]interface{})[0].(map[string]interface{})

	assert.Contains(t, block, "min_scale")
	assert.NotContains(t, block, "max_scale")
}

func TestMarkConfiguredWidgetScalePresence_NumberCardScales(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type": "Number",
			"number_cards": []interface{}{
				map[string]interface{}{
					"description": "omitted scales",
					"min_scale":   0.0,
					"max_scale":   0.0,
				},
				map[string]interface{}{
					"description": "explicit zero scales",
					"min_scale":   0.0,
					"max_scale":   0.0,
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
						"description": cty.StringVal("omitted scales"),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"description": cty.StringVal("explicit zero scales"),
						"min_scale":   cty.NumberIntVal(0),
						"max_scale":   cty.NumberIntVal(0),
					}),
				}),
			}),
		}),
	})

	marked := markConfiguredWidgetScalePresence(widgetList, rawConfig)
	cards := marked[0].(map[string]interface{})["number_cards"].([]interface{})
	omitted := cards[0].(map[string]interface{})
	explicitZero := cards[1].(map[string]interface{})

	assert.Equal(t, false, omitted["min_scale_configured"])
	assert.Equal(t, false, omitted["max_scale_configured"])
	assert.Equal(t, true, explicitZero["min_scale_configured"])
	assert.Equal(t, true, explicitZero["max_scale_configured"])
}

func TestMarkConfiguredWidgetScalePresence_WidgetConfigScales(t *testing.T) {
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
						map[string]interface{}{
							"min_scale": 0.0,
							"max_scale": 0.0,
						},
					},
				},
			}

			firstRawBlock := map[string]cty.Value{}
			for key, value := range tc.rawExtra {
				firstRawBlock[key] = value
			}
			secondRawBlock := map[string]cty.Value{
				"min_scale": cty.NumberIntVal(0),
				"max_scale": cty.NumberIntVal(0),
			}

			rawConfig := cty.ObjectVal(map[string]cty.Value{
				"widgets": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"type": cty.StringVal(tc.widgetType),
						tc.blockName: cty.TupleVal([]cty.Value{
							cty.ObjectVal(firstRawBlock),
							cty.ObjectVal(secondRawBlock),
						}),
					}),
				}),
			})

			marked := markConfiguredWidgetScalePresence(widgetList, rawConfig)
			blocks := marked[0].(map[string]interface{})[tc.blockName].([]interface{})
			omitted := blocks[0].(map[string]interface{})
			explicitZero := blocks[1].(map[string]interface{})

			assert.Equal(t, false, omitted["min_scale_configured"])
			assert.Equal(t, false, omitted["max_scale_configured"])
			assert.Equal(t, true, explicitZero["min_scale_configured"])
			assert.Equal(t, true, explicitZero["max_scale_configured"])
		})
	}
}

func TestNormalizeConfiguredWidgets_PrunesOtherPresenceSensitiveFields(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type": "Alert List",
			"should_exclude_alert_suppression_windows": false,
			"alert_list_config": []interface{}{
				map[string]interface{}{
					"limit_to": 0,
				},
			},
		},
		map[string]interface{}{
			"type": "Number",
			"number_cards": []interface{}{
				map[string]interface{}{
					"description":                              "card",
					"compare_to_previous_value":                false,
					"should_exclude_alert_suppression_windows": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Table",
			"table_config": []interface{}{
				map[string]interface{}{
					"compare_to_previous_value": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Bar Chart: Stacked",
			"stacked_bar_chart_config": []interface{}{
				map[string]interface{}{
					"show_labels":             false,
					"is_horizontal_bar_chart": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Bar Chart: Grouped",
			"grouped_bar_chart_config": []interface{}{
				map[string]interface{}{
					"show_labels":             false,
					"is_horizontal_bar_chart": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Multi Metric Table",
			"multi_metric_table_config": []interface{}{
				map[string]interface{}{
					"compare_to_previous_value": false,
				},
			},
		},
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.TupleVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Alert List"),
				"alert_list_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Number"),
				"number_cards": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"description": cty.StringVal("card"),
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Table"),
				"table_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Bar Chart: Stacked"),
				"stacked_bar_chart_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Bar Chart: Grouped"),
				"grouped_bar_chart_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Multi Metric Table"),
				"multi_metric_table_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{}),
				}),
			}),
		}),
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)

	assert.NotContains(t, normalized[0].(map[string]interface{}), "should_exclude_alert_suppression_windows")

	alertConfig := normalized[0].(map[string]interface{})["alert_list_config"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, alertConfig, "limit_to")

	card := normalized[1].(map[string]interface{})["number_cards"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, card, "compare_to_previous_value")
	assert.NotContains(t, card, "should_exclude_alert_suppression_windows")

	tableConfig := normalized[2].(map[string]interface{})["table_config"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, tableConfig, "compare_to_previous_value")

	stackedBarConfig := normalized[3].(map[string]interface{})["stacked_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, stackedBarConfig, "show_labels")
	assert.NotContains(t, stackedBarConfig, "is_horizontal_bar_chart")

	groupedBarConfig := normalized[4].(map[string]interface{})["grouped_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, groupedBarConfig, "show_labels")
	assert.NotContains(t, groupedBarConfig, "is_horizontal_bar_chart")

	multiMetricTableConfig := normalized[5].(map[string]interface{})["multi_metric_table_config"].([]interface{})[0].(map[string]interface{})
	assert.NotContains(t, multiMetricTableConfig, "compare_to_previous_value")

	widgets, err := BuildWidgets(normalized)
	require.NoError(t, err)
	require.Len(t, widgets, 6)

	_, limitSet := widgets[0].ApiAlertListWidget.GetLimitToOk()
	assert.False(t, limitSet)

	numberCards := widgets[1].ApiNumbersCardWidget.GetNumberCards()
	require.Len(t, numberCards, 1)
	_, compareSet := numberCards[0].GetCompareToPreviousValueOk()
	_, cardShouldExcludeSet := numberCards[0].GetShouldExcludeAlertSuppressionWindowsOk()
	assert.False(t, compareSet)
	assert.False(t, cardShouldExcludeSet)
}

func TestNormalizeConfiguredWidgets_PreservesExplicitFalseAndLimitZero(t *testing.T) {
	widgetList := []interface{}{
		map[string]interface{}{
			"type": "Alert List",
			"should_exclude_alert_suppression_windows": false,
			"alert_list_config": []interface{}{
				map[string]interface{}{
					"limit_to": 0,
				},
			},
		},
		map[string]interface{}{
			"type": "Number",
			"number_cards": []interface{}{
				map[string]interface{}{
					"description":                              "card",
					"compare_to_previous_value":                false,
					"should_exclude_alert_suppression_windows": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Table",
			"table_config": []interface{}{
				map[string]interface{}{
					"compare_to_previous_value": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Bar Chart: Stacked",
			"stacked_bar_chart_config": []interface{}{
				map[string]interface{}{
					"show_labels":             false,
					"is_horizontal_bar_chart": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Bar Chart: Grouped",
			"grouped_bar_chart_config": []interface{}{
				map[string]interface{}{
					"show_labels":             false,
					"is_horizontal_bar_chart": false,
				},
			},
		},
		map[string]interface{}{
			"type": "Multi Metric Table",
			"multi_metric_table_config": []interface{}{
				map[string]interface{}{
					"compare_to_previous_value": false,
				},
			},
		},
	}
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.TupleVal([]cty.Value{
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Alert List"),
				"should_exclude_alert_suppression_windows": cty.False,
				"alert_list_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"limit_to": cty.NumberIntVal(0),
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Number"),
				"number_cards": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"description":                              cty.StringVal("card"),
						"compare_to_previous_value":                cty.False,
						"should_exclude_alert_suppression_windows": cty.False,
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Table"),
				"table_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"compare_to_previous_value": cty.False,
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Bar Chart: Stacked"),
				"stacked_bar_chart_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"show_labels":             cty.False,
						"is_horizontal_bar_chart": cty.False,
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Bar Chart: Grouped"),
				"grouped_bar_chart_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"show_labels":             cty.False,
						"is_horizontal_bar_chart": cty.False,
					}),
				}),
			}),
			cty.ObjectVal(map[string]cty.Value{
				"type": cty.StringVal("Multi Metric Table"),
				"multi_metric_table_config": cty.TupleVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"compare_to_previous_value": cty.False,
					}),
				}),
			}),
		}),
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)

	assert.Contains(t, normalized[0].(map[string]interface{}), "should_exclude_alert_suppression_windows")

	alertConfig := normalized[0].(map[string]interface{})["alert_list_config"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, alertConfig, "limit_to")

	card := normalized[1].(map[string]interface{})["number_cards"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, card, "compare_to_previous_value")
	assert.Contains(t, card, "should_exclude_alert_suppression_windows")

	tableConfig := normalized[2].(map[string]interface{})["table_config"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, tableConfig, "compare_to_previous_value")

	stackedBarConfig := normalized[3].(map[string]interface{})["stacked_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, stackedBarConfig, "show_labels")
	assert.Contains(t, stackedBarConfig, "is_horizontal_bar_chart")

	groupedBarConfig := normalized[4].(map[string]interface{})["grouped_bar_chart_config"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, groupedBarConfig, "show_labels")
	assert.Contains(t, groupedBarConfig, "is_horizontal_bar_chart")

	multiMetricTableConfig := normalized[5].(map[string]interface{})["multi_metric_table_config"].([]interface{})[0].(map[string]interface{})
	assert.Contains(t, multiMetricTableConfig, "compare_to_previous_value")

	widgets, err := BuildWidgets(normalized)
	require.NoError(t, err)
	require.Len(t, widgets, 6)

	limit, limitSet := widgets[0].ApiAlertListWidget.GetLimitToOk()
	require.True(t, limitSet)
	require.NotNil(t, limit)
	assert.Equal(t, int32(0), *limit)

	numberCards := widgets[1].ApiNumbersCardWidget.GetNumberCards()
	require.Len(t, numberCards, 1)
	compare, compareSet := numberCards[0].GetCompareToPreviousValueOk()
	cardShouldExclude, cardShouldExcludeSet := numberCards[0].GetShouldExcludeAlertSuppressionWindowsOk()
	require.True(t, compareSet)
	require.True(t, cardShouldExcludeSet)
	require.NotNil(t, compare)
	require.NotNil(t, cardShouldExclude)
	assert.False(t, *compare)
	assert.False(t, *cardShouldExclude)
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

func TestNormalizeConfiguredWidgets_FallbackKeepsOriginalListWhenRawWidgetsEmpty(t *testing.T) {
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
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"widgets": cty.EmptyTupleVal,
	})

	normalized := normalizeConfiguredWidgets(widgetList, rawConfig)
	normalized[0].(map[string]interface{})["sentinel"] = true

	assert.Contains(t, widgetList[0].(map[string]interface{}), "sentinel")
}
