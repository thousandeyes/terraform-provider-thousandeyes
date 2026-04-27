package thousandeyes

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccThousandEyesDashboard_presenceSensitiveDefaultsStablePlan(t *testing.T) {
	cases := []struct {
		name         string
		resourceName string
		config       string
		attrs        []string
	}{
		{
			name:         "map",
			resourceName: "thousandeyes_dashboard.test_dashboard_map_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_map_basic.tf",
				"      min_scale = 0\n",
				"      max_scale = 100\n",
				"      is_geo_map_per_test = false\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.geo_map_config.0.min_scale",
				"widgets.0.geo_map_config.0.max_scale",
				"widgets.0.geo_map_config.0.is_geo_map_per_test",
			},
		},
		{
			name:         "timeseries",
			resourceName: "thousandeyes_dashboard.test_dashboard_timeseries_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_timeseries_basic.tf",
				"      show_timeseries_overall_baseline = false\n",
				"      is_timeseries_one_chart_per_line = false\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.timeseries_config.0.min_scale",
				"widgets.0.timeseries_config.0.max_scale",
				"widgets.0.timeseries_config.0.show_timeseries_overall_baseline",
				"widgets.0.timeseries_config.0.is_timeseries_one_chart_per_line",
			},
		},
		{
			name:         "stacked_area",
			resourceName: "thousandeyes_dashboard.test_dashboard_stacked_area_defaults",
			config:       testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_stacked_area_defaults.tf"),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.stacked_area_config.0.min_scale",
				"widgets.0.stacked_area_config.0.max_scale",
			},
		},
		{
			name:         "box_and_whiskers",
			resourceName: "thousandeyes_dashboard.test_dashboard_box_and_whiskers_widget",
			config:       testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_box_and_whiskers_basic.tf"),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.box_and_whiskers_config.0.min_scale",
				"widgets.0.box_and_whiskers_config.0.max_scale",
			},
		},
		{
			name:         "color_grid",
			resourceName: "thousandeyes_dashboard.test_dashboard_color_grid_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_color_grid_basic.tf",
				"      min_scale      = 0\n",
				"      max_scale      = 100\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.color_grid_config.0.min_scale",
				"widgets.0.color_grid_config.0.max_scale",
			},
		},
		{
			name:         "number",
			resourceName: "thousandeyes_dashboard.test_dashboard_number_defaults",
			config:       testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_number_defaults.tf"),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.number_cards.0.min_scale",
				"widgets.0.number_cards.0.max_scale",
				"widgets.0.number_cards.0.compare_to_previous_value",
				"widgets.0.number_cards.0.should_exclude_alert_suppression_windows",
			},
		},
		{
			name:         "alert_list",
			resourceName: "thousandeyes_dashboard.test_dashboard_alert_list_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_alert_list_basic.tf",
				"      limit_to    = 15\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.alert_list_config.0.limit_to",
			},
		},
		{
			name:         "table",
			resourceName: "thousandeyes_dashboard.test_dashboard_table_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_table_basic.tf",
				"      compare_to_previous_value = true\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.table_config.0.compare_to_previous_value",
			},
		},
		{
			name:         "stacked_bar_chart",
			resourceName: "thousandeyes_dashboard.test_dashboard_stacked_bar_chart_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_stacked_bar_chart_basic.tf",
				"      show_labels             = true\n",
				"      is_horizontal_bar_chart = true\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.stacked_bar_chart_config.0.show_labels",
				"widgets.0.stacked_bar_chart_config.0.is_horizontal_bar_chart",
			},
		},
		{
			name:         "grouped_bar_chart",
			resourceName: "thousandeyes_dashboard.test_dashboard_grouped_bar_chart_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_grouped_bar_chart_basic.tf",
				"      show_labels             = true\n",
				"      is_horizontal_bar_chart = false\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.grouped_bar_chart_config.0.show_labels",
				"widgets.0.grouped_bar_chart_config.0.is_horizontal_bar_chart",
			},
		},
		{
			name:         "multi_metric_table",
			resourceName: "thousandeyes_dashboard.test_dashboard_multi_metric_table_widget",
			config: withoutConfigLines("acceptance_resources/dashboard/widget_multi_metric_table_basic.tf",
				"      compare_to_previous_value = true\n",
			),
			attrs: []string{
				"widgets.0.should_exclude_alert_suppression_windows",
				"widgets.0.multi_metric_table_config.0.compare_to_previous_value",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckDashboardResourceDestroy,
				Steps: []resource.TestStep{
					{
						Config: tc.config,
						Check:  resource.ComposeTestCheckFunc(testCheckResourceAttrsSet(tc.resourceName, tc.attrs)...),
					},
					{
						Config:   tc.config,
						PlanOnly: true,
					},
				},
			})
		})
	}
}

func withoutConfigLines(path string, removals ...string) string {
	cfg := testAccThousandEyesDashboardConfig(path)
	for _, removal := range removals {
		cfg = strings.ReplaceAll(cfg, removal, "")
	}
	return cfg
}

func testCheckResourceAttrsSet(resourceName string, attrs []string) []resource.TestCheckFunc {
	checks := make([]resource.TestCheckFunc, 0, len(attrs))
	for _, attr := range attrs {
		checks = append(checks, resource.TestCheckResourceAttrSet(resourceName, attr))
	}
	return checks
}
