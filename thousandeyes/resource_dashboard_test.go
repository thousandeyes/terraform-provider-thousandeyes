package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestAccThousandEyesDashboard(t *testing.T) {
	var resourceName = "thousandeyes_dashboard.test_dashboard"
	var resourceNameTimeRange = "thousandeyes_dashboard.test_dashboard_time_range"
	var resourceNameMapWidget = "thousandeyes_dashboard.test_dashboard_map_widget"
	var resourceNameAgentStatusWidget = "thousandeyes_dashboard.test_dashboard_agent_status_widget"
	var resourceNameTimeseriesWidget = "thousandeyes_dashboard.test_dashboard_timeseries_widget"
	var testCases = []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:                 "create_update_delete_dashboard_duration_test",
			createResourceFile:   "acceptance_resources/dashboard/basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard"),
				resource.TestCheckResourceAttr(resourceName, "description", "Test Dashboard Description"),
				resource.TestCheckResourceAttr(resourceName, "is_private", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_global_override", "false"),
				resource.TestCheckResourceAttr(resourceName, "default_timespan.0.duration", "3600"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "description", "Updated Test Dashboard Description"),
				resource.TestCheckResourceAttr(resourceName, "is_private", "true"),
				resource.TestCheckResourceAttr(resourceName, "is_global_override", "true"),
				resource.TestCheckResourceAttr(resourceName, "default_timespan.0.duration", "7200"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_time_range_test",
			createResourceFile:   "acceptance_resources/dashboard/time_range_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/time_range_update.tf",
			resourceName:         resourceNameTimeRange,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeRange, "title", "Test Dashboard Time Range"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "description", "Test Dashboard with Time Range"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "is_private", "false"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "default_timespan.0.start", "2026-01-01T00:00:00Z"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "default_timespan.0.end", "2026-02-01T00:00:00Z"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeRange, "title", "Test Dashboard Time Range (Updated)"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "description", "Updated Test Dashboard with Time Range"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "is_private", "true"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "default_timespan.0.start", "2026-02-01T00:00:00Z"),
				resource.TestCheckResourceAttr(resourceNameTimeRange, "default_timespan.0.end", "2026-03-01T00:00:00Z"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_map_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_map_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_map_update.tf",
			resourceName:         resourceNameMapWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameMapWidget, "title", "Test Dashboard Map Widget"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "description", "Test Dashboard with Map Widget"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.type", "Map"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.title", "Test Map Widget"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.group_by", "COUNTRY"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.is_geo_map_per_test", "false"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameMapWidget, "title", "Test Dashboard Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "description", "Test Dashboard with Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.type", "Map"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.title", "Test Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.group_by", "COUNTRY"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.is_geo_map_per_test", "true"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_agent_status_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_agent_status_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_agent_status_update.tf",
			resourceName:         resourceNameAgentStatusWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "title", "Test Dashboard Agent Status Widget"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "description", "Test Dashboard with Agent Status Widget"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.type", "Agent Status"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.title", "Test Agent Status Widget"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.visual_mode", "Full"),
				//resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.agent_type", "Endpoint Agents"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "title", "Test Dashboard Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "description", "Test Dashboard with Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.type", "Agent Status"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.title", "Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.visual_mode", "Full"),
				//resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.show", "All Assigned Agents"),
				//resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.agent_type", "Endpoint Agents"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_timeseries_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_timeseries_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_timeseries_update.tf",
			resourceName:         resourceNameTimeseriesWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "title", "Test Dashboard Timeseries Widget"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "description", "Test Dashboard with Timeseries Widget"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.type", "Time Series: Line"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.title", "Test Timeseries Widget"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.data_source", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.metric_group", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.metric", "ALERT_COUNT_AGENT"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.measure.0.type", "TOTAL"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.group_by", "AGENT"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.show_timeseries_overall_baseline", "false"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.is_timeseries_one_chart_per_line", "false"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "title", "Test Dashboard Timeseries Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "description", "Test Dashboard with Timeseries Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.type", "Time Series: Line"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.title", "Test Timeseries Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.data_source", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.metric_group", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.metric", "ALERT_COUNT_AGENT"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.measure.0.type", "TOTAL"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.group_by", "AGENT"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.show_timeseries_overall_baseline", "true"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.is_timeseries_one_chart_per_line", "true"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      tc.checkDestroyFunction,
				Steps: []resource.TestStep{
					{
						Config: testAccThousandEyesDashboardConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						ResourceName:      tc.resourceName,
						ImportState:       true,
						ImportStateVerify: true,
					},
					{
						Config: testAccThousandEyesDashboardConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDashboardResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_dashboard",
			GetResource: func(id string) (interface{}, error) {
				return getDashboard(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesDashboardConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getDashboard(id string) (interface{}, error) {
	api := (*dashboards.DashboardsAPIService)(&testClient.Common)
	req := api.GetDashboard(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
