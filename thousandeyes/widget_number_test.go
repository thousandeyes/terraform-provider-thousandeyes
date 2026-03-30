package thousandeyes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestBuildNumberWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		validate func(t *testing.T, widget dashboards.ApiWidget)
	}{
		{
			name: "basic number widget",
			input: map[string]interface{}{
				"type":        "Number",
				"title":       "Test Number",
				"visual_mode": "Full",
				"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiNumbersCardWidget
				assert.NotNil(t, w)
				assert.Equal(t, "Number", w.GetType())
				assert.Equal(t, "Test Number", w.GetTitle())
				assert.Equal(t, dashboards.VisualMode("Full"), w.GetVisualMode())
				assert.Equal(t, dashboards.NumbersCardDatasource("CLOUD_AND_ENTERPRISE_AGENTS"), w.GetDataSource())
			},
		},
		{
			name: "number widget with cards",
			input: map[string]interface{}{
				"type":        "Number",
				"title":       "Multi-Card Number",
				"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
				"number_cards": []interface{}{
					map[string]interface{}{
						"description":  "Card 1",
						"data_source":  "CLOUD_AND_ENTERPRISE_AGENTS",
						"metric_group": "WEB_HTTP_SERVER",
						"metric":       "RESPONSE_TIME",
						"direction":    "TO_TARGET",
						"measure": []interface{}{
							map[string]interface{}{
								"type": "MEAN",
							},
						},
					},
					map[string]interface{}{
						"description": "Card 2",
						"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
						"metric":      "AVAILABILITY",
						"measure": []interface{}{
							map[string]interface{}{
								"type": "MEAN",
							},
						},
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiNumbersCardWidget
				assert.NotNil(t, w)
				cards := w.GetNumberCards()
				assert.Len(t, cards, 2)

				assert.Equal(t, "Card 1", cards[0].GetDescription())
				assert.Equal(t, dashboards.NumbersCardDatasource("CLOUD_AND_ENTERPRISE_AGENTS"), cards[0].GetDataSource())
				assert.Equal(t, dashboards.MetricGroup("WEB_HTTP_SERVER"), cards[0].GetMetricGroup())
				assert.Equal(t, dashboards.DashboardMetric("RESPONSE_TIME"), cards[0].GetMetric())
				assert.Equal(t, dashboards.DashboardMetricDirection("TO_TARGET"), cards[0].GetDirection())
				measure := cards[0].GetMeasure()
				assert.Equal(t, dashboards.WidgetMeasureType("MEAN"), measure.GetType())

				assert.Equal(t, "Card 2", cards[1].GetDescription())
				assert.Equal(t, dashboards.DashboardMetric("AVAILABILITY"), cards[1].GetMetric())
			},
		},
		{
			name: "number card with scale, timespan, and filters",
			input: map[string]interface{}{
				"type":  "Number",
				"title": "Detailed Card",
				"number_cards": []interface{}{
					map[string]interface{}{
						"min_scale":                  0.0,
						"max_scale":                  100.0,
						"unit":                       "Kilo",
						"compare_to_previous_value":  true,
						"should_exclude_alert_suppression_windows": true,
						"fixed_timespan": []interface{}{
							map[string]interface{}{
								"value": 24,
								"unit":  "Hours",
							},
						},
						"filter": []interface{}{
							map[string]interface{}{
								"property": "TEST",
								"values":   []interface{}{"12345"},
							},
						},
					},
				},
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiNumbersCardWidget
				assert.NotNil(t, w)
				cards := w.GetNumberCards()
				assert.Len(t, cards, 1)

				card := cards[0]
				assert.Equal(t, float32(0.0), card.GetMinScale())
				assert.Equal(t, float32(100.0), card.GetMaxScale())
				assert.Equal(t, dashboards.ApiWidgetFixedYScalePrefix("Kilo"), card.GetUnit())
				assert.True(t, card.GetCompareToPreviousValue())
				assert.True(t, card.GetShouldExcludeAlertSuppressionWindows())

				fts := card.GetFixedTimespan()
				assert.Equal(t, int32(24), fts.GetValue())
				assert.Equal(t, dashboards.LegacyDurationUnit("Hours"), fts.GetUnit())

				filters := card.GetFilters()
				assert.Contains(t, filters, "TEST")
				assert.Equal(t, []interface{}{"12345"}, filters["TEST"])
			},
		},
		{
			name: "removing all cards sends empty array",
			input: map[string]interface{}{
				"type":        "Number",
				"title":       "Empty Cards",
				"data_source": "CLOUD_AND_ENTERPRISE_AGENTS",
			},
			validate: func(t *testing.T, widget dashboards.ApiWidget) {
				w := widget.ApiNumbersCardWidget
				assert.NotNil(t, w)
				cards, ok := w.GetNumberCardsOk()
				assert.True(t, ok, "NumberCards should be explicitly set (not nil)")
				assert.Empty(t, cards, "NumberCards should be an empty slice")
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := buildNumberWidget(tc.input)
			tc.validate(t, widget)
		})
	}
}

func TestMapNumberWidget(t *testing.T) {
	tests := []struct {
		name     string
		input    func() dashboards.ApiWidget
		validate func(t *testing.T, data map[string]interface{})
	}{
		{
			name: "basic number widget",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiNumbersCardWidget("Number")
				w.SetId("widget-num-1")
				w.SetTitle("Test Number")
				w.SetVisualMode(dashboards.VisualMode("Full"))
				w.SetDataSource(dashboards.NumbersCardDatasource("CLOUD_AND_ENTERPRISE_AGENTS"))
				return dashboards.ApiNumbersCardWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Number", data["type"])
				assert.Equal(t, "widget-num-1", data["id"])
				assert.Equal(t, "Test Number", data["title"])
				assert.Equal(t, "Full", data["visual_mode"])
				assert.Equal(t, "CLOUD_AND_ENTERPRISE_AGENTS", data["data_source"])
			},
		},
		{
			name: "number widget with cards round-trip",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiNumbersCardWidget("Number")
				w.SetTitle("Multi-Card")
				w.SetDataSource(dashboards.NumbersCardDatasource("CLOUD_AND_ENTERPRISE_AGENTS"))

				card1 := dashboards.NewApiNumbersCard()
				card1.SetId("card-1")
				card1.SetDescription("Response Time")
				card1.SetDataSource(dashboards.NumbersCardDatasource("CLOUD_AND_ENTERPRISE_AGENTS"))
				card1.SetMetricGroup(dashboards.MetricGroup("WEB_HTTP_SERVER"))
				card1.SetMetric(dashboards.DashboardMetric("RESPONSE_TIME"))
				card1.SetDirection(dashboards.DashboardMetricDirection("TO_TARGET"))
				measure := dashboards.NewApiWidgetMeasure()
				measure.SetType(dashboards.WidgetMeasureType("MEAN"))
				card1.SetMeasure(*measure)

				card2 := dashboards.NewApiNumbersCard()
				card2.SetId("card-2")
				card2.SetDescription("Availability")
				card2.SetMetric(dashboards.DashboardMetric("AVAILABILITY"))

				w.SetNumberCards([]dashboards.ApiNumbersCard{*card1, *card2})
				return dashboards.ApiNumbersCardWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				assert.Equal(t, "Number", data["type"])
				assert.Equal(t, "Multi-Card", data["title"])

				cards := data["number_cards"].([]interface{})
				assert.Len(t, cards, 2)

				c1 := cards[0].(map[string]interface{})
				assert.Equal(t, "card-1", c1["id"])
				assert.Equal(t, "Response Time", c1["description"])
				assert.Equal(t, "CLOUD_AND_ENTERPRISE_AGENTS", c1["data_source"])
				assert.Equal(t, "WEB_HTTP_SERVER", c1["metric_group"])
				assert.Equal(t, "RESPONSE_TIME", c1["metric"])
				assert.Equal(t, "TO_TARGET", c1["direction"])
				measureList := c1["measure"].([]interface{})
				assert.Len(t, measureList, 1)
				measureMap := measureList[0].(map[string]interface{})
				assert.Equal(t, "MEAN", measureMap["type"])

				c2 := cards[1].(map[string]interface{})
				assert.Equal(t, "card-2", c2["id"])
				assert.Equal(t, "Availability", c2["description"])
				assert.Equal(t, "AVAILABILITY", c2["metric"])
			},
		},
		{
			name: "number card with scale and timespan",
			input: func() dashboards.ApiWidget {
				w := dashboards.NewApiNumbersCardWidget("Number")
				w.SetTitle("Scaled")

				card := dashboards.NewApiNumbersCard()
				card.SetMinScale(0)
				card.SetMaxScale(100)
				card.SetUnit(dashboards.ApiWidgetFixedYScalePrefix("Kilo"))
				card.SetCompareToPreviousValue(true)
				card.SetShouldExcludeAlertSuppressionWindows(true)

				fts := dashboards.NewApiDuration()
				fts.SetValue(24)
				fts.SetUnit(dashboards.LegacyDurationUnit("Hours"))
				card.SetFixedTimespan(*fts)

				card.SetFilters(map[string][]interface{}{
					"TEST": {"12345"},
				})

				w.SetNumberCards([]dashboards.ApiNumbersCard{*card})
				return dashboards.ApiNumbersCardWidgetAsApiWidget(w)
			},
			validate: func(t *testing.T, data map[string]interface{}) {
				cards := data["number_cards"].([]interface{})
				assert.Len(t, cards, 1)

				c := cards[0].(map[string]interface{})
				assert.Equal(t, float64(0), c["min_scale"])
				assert.Equal(t, float64(100), c["max_scale"])
				assert.Equal(t, "Kilo", c["unit"])
				assert.Equal(t, true, c["compare_to_previous_value"])
				assert.Equal(t, true, c["should_exclude_alert_suppression_windows"])

				fts := c["fixed_timespan"].([]interface{})[0].(map[string]interface{})
				assert.Equal(t, 24, fts["value"])
				assert.Equal(t, "Hours", fts["unit"])

				filters := c["filter"].([]interface{})
				assert.Len(t, filters, 1)
				f := filters[0].(map[string]interface{})
				assert.Equal(t, "TEST", f["property"])
				assert.Equal(t, []interface{}{"12345"}, f["values"])
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			widget := tc.input()
			data, err := mapNumberWidget(widget)
			assert.NoError(t, err)
			assert.NotNil(t, data)
			tc.validate(t, data)
		})
	}
}

func TestMapNumberWidgetNil(t *testing.T) {
	widget := dashboards.ApiWidget{}
	data, err := mapNumberWidget(widget)
	assert.NoError(t, err)
	assert.Nil(t, data)
}

func TestNumberWidgetRegistration(t *testing.T) {
	input := map[string]interface{}{
		"type":  "Number",
		"title": "Registered Number",
	}
	widget, err := BuildWidget(input)
	assert.NoError(t, err)
	assert.NotNil(t, widget.ApiNumbersCardWidget)

	data, err := MapWidget(widget)
	assert.NoError(t, err)
	assert.Equal(t, "Number", data["type"])
	assert.Equal(t, "Registered Number", data["title"])
}
