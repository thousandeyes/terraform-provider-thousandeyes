package thousandeyes

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func TestAccThousandEyesDashboard(t *testing.T) {
	var resourceName = "thousandeyes_dashboard.test_dashboard"
	var resourceNameTimeRange = "thousandeyes_dashboard.test_dashboard_time_range"
	var resourceNameMapDefaults = "thousandeyes_dashboard.test_dashboard_map_defaults"
	var resourceNameMapWidget = "thousandeyes_dashboard.test_dashboard_map_widget"
	var resourceNameAgentStatusDefaults = "thousandeyes_dashboard.test_dashboard_agent_status_defaults"
	var resourceNameAgentStatusWidget = "thousandeyes_dashboard.test_dashboard_agent_status_widget"
	var resourceNameTimeseriesDefaults = "thousandeyes_dashboard.test_dashboard_timeseries_defaults"
	var resourceNameTimeseriesWidget = "thousandeyes_dashboard.test_dashboard_timeseries_widget"
	var resourceNameStackedAreaDefaults = "thousandeyes_dashboard.test_dashboard_stacked_area_defaults"
	var resourceNameStackedAreaWidget = "thousandeyes_dashboard.test_dashboard_stacked_area_widget"
	var resourceNamePieChartDefaults = "thousandeyes_dashboard.test_dashboard_pie_chart_defaults"
	var resourceNamePieChartWidget = "thousandeyes_dashboard.test_dashboard_pie_chart_widget"
	var resourceNameBoxAndWhiskersDefaults = "thousandeyes_dashboard.test_dashboard_box_and_whiskers_defaults"
	var resourceNameBoxAndWhiskersWidget = "thousandeyes_dashboard.test_dashboard_box_and_whiskers_widget"
	var resourceNameFilterWidget = "thousandeyes_dashboard.test_dashboard_filter_widget"
	//var resourceNameListWidget = "thousandeyes_dashboard.test_dashboard_list_widget"
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
			name:                 "create_dashboard_map_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_map_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_map_defaults.tf",
			resourceName:         resourceNameMapDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "title", "Test Dashboard Map Defaults"),
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "widgets.0.type", "Map"),
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "widgets.0.title", "Map With Defaults"),
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "widgets.0.data_source", "ALERTS"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "title", "Test Dashboard Map Defaults"),
				resource.TestCheckResourceAttr(resourceNameMapDefaults, "widgets.0.type", "Map"),
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
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.min_scale", "0"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.max_scale", "100"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.group_by", "COUNTRY"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.is_geo_map_per_test", "false"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameMapWidget, "title", "Test Dashboard Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "description", "Test Dashboard with Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.type", "Map"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.title", "Test Map Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.min_scale", "10"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.max_scale", "200"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.group_by", "CONTINENT"),
				resource.TestCheckResourceAttr(resourceNameMapWidget, "widgets.0.geo_map_config.0.is_geo_map_per_test", "true"),
			},
		},
		{
			name:                 "create_dashboard_agent_status_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_agent_status_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_agent_status_defaults.tf",
			resourceName:         resourceNameAgentStatusDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "title", "Test Dashboard Agent Status Defaults"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "widgets.0.type", "Agent Status"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "widgets.0.title", "Agent Status With Defaults"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "widgets.0.data_source", "CLOUD_AND_ENTERPRISE_AGENTS"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "title", "Test Dashboard Agent Status Defaults"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusDefaults, "widgets.0.type", "Agent Status"),
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
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.data_source", "CLOUD_AND_ENTERPRISE_AGENTS"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.type", "Agent Status"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.title", "Test Agent Status Widget"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.agent_type", "Enterprise Agents"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "title", "Test Dashboard Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "description", "Test Dashboard with Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.data_source", "ENDPOINT_AGENTS"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.type", "Agent Status"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.title", "Agent Status Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.show", "Owned Agents"),
				resource.TestCheckResourceAttr(resourceNameAgentStatusWidget, "widgets.0.agent_status_config.0.agent_type", "Endpoint Agents"),
			},
		},
		{
			name:                 "create_dashboard_timeseries_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_timeseries_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_timeseries_defaults.tf",
			resourceName:         resourceNameTimeseriesDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "title", "Test Dashboard Timeseries Defaults"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "widgets.0.type", "Time Series: Line"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "widgets.0.title", "Timeseries With Defaults"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "widgets.0.data_source", "ALERTS"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "title", "Test Dashboard Timeseries Defaults"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesDefaults, "widgets.0.type", "Time Series: Line"),
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
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.group_by", "TEST"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.show_timeseries_overall_baseline", "true"),
				resource.TestCheckResourceAttr(resourceNameTimeseriesWidget, "widgets.0.timeseries_config.0.is_timeseries_one_chart_per_line", "true"),
			},
		},
		{
			name:                 "create_dashboard_stacked_area_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_stacked_area_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_stacked_area_defaults.tf",
			resourceName:         resourceNameStackedAreaDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "title", "Test Dashboard Stacked Area Defaults"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.type", "Time Series: Stacked Area"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.title", "Stacked Area With Defaults"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.stacked_area_config.0.group_by", "CLOUD_NATIVE_MONITORING-REGION"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "title", "Test Dashboard Stacked Area Defaults"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.type", "Time Series: Stacked Area"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaDefaults, "widgets.0.stacked_area_config.0.group_by", "CLOUD_NATIVE_MONITORING-REGION"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_stacked_area_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_stacked_area_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_stacked_area_update.tf",
			resourceName:         resourceNameStackedAreaWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "title", "Test Dashboard Stacked Area Widget"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "description", "Test Dashboard with Stacked Area Widget"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.type", "Time Series: Stacked Area"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.title", "Test Stacked Area Widget"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.metric_group", "CLOUD_NATIVE_MONITORING-EVENTS"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.metric", "CLOUD_NATIVE_MONITORING-ALL_EVENTS"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.measure.0.type", "CLOUD_NATIVE_MONITORING-SUM"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.stacked_area_config.0.group_by", "CLOUD_NATIVE_MONITORING-REGION"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "title", "Test Dashboard Stacked Area Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "description", "Test Dashboard with Stacked Area Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.type", "Time Series: Stacked Area"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.title", "Test Stacked Area Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.metric_group", "CLOUD_NATIVE_MONITORING-EVENTS"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.metric", "CLOUD_NATIVE_MONITORING-ALL_EVENTS"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.measure.0.type", "CLOUD_NATIVE_MONITORING-SUM"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameStackedAreaWidget, "widgets.0.stacked_area_config.0.group_by", "CLOUD_NATIVE_MONITORING-ACCOUNT"),
			},
		},
		{
			name:                 "create_dashboard_pie_chart_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_pie_chart_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_pie_chart_defaults.tf",
			resourceName:         resourceNamePieChartDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "title", "Test Dashboard Pie Chart Defaults"),
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "widgets.0.type", "Pie Chart"),
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "widgets.0.title", "Pie Chart With Defaults"),
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "title", "Test Dashboard Pie Chart Defaults"),
				resource.TestCheckResourceAttr(resourceNamePieChartDefaults, "widgets.0.type", "Pie Chart"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_pie_chart_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_pie_chart_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_pie_chart_update.tf",
			resourceName:         resourceNamePieChartWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "title", "Test Dashboard Pie Chart Widget"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "description", "Test Dashboard with Pie Chart Widget"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.type", "Pie Chart"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.title", "Test Pie Chart Widget"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.metric_group", "CLOUD_NATIVE_MONITORING-EVENTS"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.metric", "CLOUD_NATIVE_MONITORING-ALL_EVENTS"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.measure.0.type", "CLOUD_NATIVE_MONITORING-SUM"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.pie_chart_config.0.group_by", "CLOUD_NATIVE_MONITORING-REGION"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "title", "Test Dashboard Pie Chart Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "description", "Test Dashboard with Pie Chart Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.type", "Pie Chart"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.title", "Test Pie Chart Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.data_source", "CLOUD_NATIVE_MONITORING"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.metric_group", "CLOUD_NATIVE_MONITORING-EVENTS"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.metric", "CLOUD_NATIVE_MONITORING-ALL_EVENTS"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.measure.0.type", "CLOUD_NATIVE_MONITORING-SUM"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNamePieChartWidget, "widgets.0.pie_chart_config.0.group_by", "CLOUD_NATIVE_MONITORING-ACCOUNT"),
			},
		},
		{
			name:                 "create_dashboard_box_and_whiskers_defaults_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_box_and_whiskers_defaults.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_box_and_whiskers_defaults.tf",
			resourceName:         resourceNameBoxAndWhiskersDefaults,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "title", "Test Dashboard Box and Whiskers Defaults"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "widgets.0.type", "Box and Whiskers"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "widgets.0.title", "Box and Whiskers With Defaults"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "widgets.0.data_source", "ALERTS"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "title", "Test Dashboard Box and Whiskers Defaults"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersDefaults, "widgets.0.type", "Box and Whiskers"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_box_and_whiskers_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_box_and_whiskers_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_box_and_whiskers_update.tf",
			resourceName:         resourceNameBoxAndWhiskersWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "title", "Test Dashboard Box and Whiskers Widget"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "description", "Test Dashboard with Box and Whiskers Widget"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.type", "Box and Whiskers"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.title", "Test Box and Whiskers Widget"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.data_source", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.metric_group", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.metric", "ALERT_COUNT_AGENT"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.measure.0.type", "MEAN"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.box_and_whiskers_config.0.group_by", "COUNTRY"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "title", "Test Dashboard Box and Whiskers Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "description", "Test Dashboard with Box and Whiskers Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.type", "Box and Whiskers"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.title", "Test Box and Whiskers Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.visual_mode", "Full"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.data_source", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.metric_group", "ALERTS"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.metric", "ALERT_COUNT_AGENT"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.measure.0.type", "MEAN"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.fixed_timespan.0.value", "1"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
				resource.TestCheckResourceAttr(resourceNameBoxAndWhiskersWidget, "widgets.0.box_and_whiskers_config.0.group_by", "CONTINENT"),
			},
		},
		{
			name:                 "create_update_delete_dashboard_filter_widget_test",
			createResourceFile:   "acceptance_resources/dashboard/widget_filter_basic.tf",
			updateResourceFile:   "acceptance_resources/dashboard/widget_filter_update.tf",
			resourceName:         resourceNameFilterWidget,
			checkDestroyFunction: testAccCheckDashboardResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "title", "Test Dashboard Filter Widget"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "description", "Test Dashboard with Widget Filter"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.type", "Time Series: Line"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.#", "1"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.0.property", "INSIGHTS_NETWORK"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.#", "3"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "32133"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "262287"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "46606"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "title", "Test Dashboard Filter Widget (Updated)"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "description", "Test Dashboard with Widget Filter (Updated)"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.type", "Time Series: Line"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.#", "1"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.0.property", "INSIGHTS_NETWORK"),
				resource.TestCheckResourceAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.#", "5"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "32133"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "4230"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "8075"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "262287"),
				resource.TestCheckTypeSetElemAttr(resourceNameFilterWidget, "widgets.0.filter.0.values.*", "46606"),
			},
		},
		// This API is return invalid sortDirection
		//  "sortDirection": "ASC",
		// The expected values according to the docs
		// sortDirection: LegacyWidgetSortDirection (Deprecated) Specifies the order in which cards are sorted.
		// enum = ["Ascending", "Descending"]
		//{
		//	name:                 "create_update_delete_dashboard_list_widget_test",
		//	createResourceFile:   "acceptance_resources/dashboard/widget_list_basic.tf",
		//	updateResourceFile:   "acceptance_resources/dashboard/widget_list_update.tf",
		//	resourceName:         resourceNameListWidget,
		//	checkDestroyFunction: testAccCheckDashboardResourceDestroy,
		//	checkCreateFunc: []resource.TestCheckFunc{
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "title", "Test Dashboard List Widget"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "description", "Test Dashboard with List Widget"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.type", "List"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.title", "Test List Widget"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.visual_mode", "Full"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.data_source", "EVENT_DETECTION"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.measure.0.type", "MEAN"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.fixed_timespan.0.value", "1"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.list_config.0.active_within_value", "7"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.list_config.0.active_within_unit", "Days"),
		//	},
		//	checkUpdateFunc: []resource.TestCheckFunc{
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "title", "Test Dashboard List Widget (Updated)"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "description", "Test Dashboard with List Widget (Updated)"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.type", "List"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.title", "Test List Widget (Updated)"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.visual_mode", "Full"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.data_source", "ALERTS"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.measure.0.type", "MEAN"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.fixed_timespan.0.value", "1"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.fixed_timespan.0.unit", "Days"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.list_config.0.active_within_value", "14"),
		//		resource.TestCheckResourceAttr(resourceNameListWidget, "widgets.0.list_config.0.active_within_unit", "Days"),
		//	},
		//},
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

// TestAccThousandEyesDashboard_removeAllWidgets is a dedicated test for the bug where removing
// all widget blocks from config produced no diff and left widgets unchanged on the API.
func TestAccThousandEyesDashboard_removeAllWidgets(t *testing.T) {
	resourceName := "thousandeyes_dashboard.test_dashboard_remove_all_widgets"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDashboardResourceDestroy,
		Steps: []resource.TestStep{
			{
				// Create dashboard with two widgets.
				Config: testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_remove_all_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard Remove All Widgets"),
					resource.TestCheckResourceAttr(resourceName, "widgets.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "widgets.0.type", "Time Series: Line"),
					resource.TestCheckResourceAttr(resourceName, "widgets.1.type", "Box and Whiskers"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// Remove all widget blocks — previously produced no diff (the bug).
				Config: testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_remove_all_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "title", "Test Dashboard Remove All Widgets"),
					resource.TestCheckResourceAttr(resourceName, "widgets.#", "0"),
				),
			},
		},
	})
}

// TestAccThousandEyesDashboard_preserveUnmanagedWidgets verifies that when a dashboard
// contains widget types not supported by the provider, those widgets are preserved
// across Terraform updates rather than being silently dropped.
func TestAccThousandEyesDashboard_preserveUnmanagedWidgets(t *testing.T) {
	resourceName := "thousandeyes_dashboard.test_dashboard_preserve_unmanaged"
	var dashboardID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckDashboardResourceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_unmanaged_basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "widgets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "widgets.0.type", "Time Series: Line"),
					resource.TestCheckResourceAttr(resourceName, "widgets.0.title", "Managed Timeseries Widget"),
					func(s *terraform.State) error {
						rs, ok := s.RootModule().Resources[resourceName]
						if !ok {
							return fmt.Errorf("resource not found: %s", resourceName)
						}
						dashboardID = rs.Primary.ID
						return nil
					},
				),
			},
			{
				PreConfig: func() {
					addUnmanagedWidgetViaDashboardAPI(t, dashboardID)
				},
				Config: testAccThousandEyesDashboardConfig("acceptance_resources/dashboard/widget_unmanaged_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "widgets.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "widgets.0.type", "Time Series: Line"),
					resource.TestCheckResourceAttr(resourceName, "widgets.0.title", "Managed Timeseries Widget (Updated)"),
					func(s *terraform.State) error {
						return checkDashboardTotalWidgetCount(dashboardID, 2)
					},
				),
			},
		},
	})
}

// addUnmanagedWidgetViaDashboardAPI adds a Color Grid widget (unsupported by the provider)
// directly to the dashboard via the API, simulating a widget created outside Terraform.
func addUnmanagedWidgetViaDashboardAPI(t *testing.T, dashboardID string) {
	t.Helper()
	api := (*dashboards.DashboardsAPIService)(&testClient.Common)
	ctx := testClient.GetConfig().Context

	getReq := api.GetDashboard(dashboardID)
	getReq = SetAidFromContext(ctx, getReq)
	existing, _, err := getReq.Execute()
	if err != nil {
		t.Fatalf("failed to GET dashboard %s: %v", dashboardID, err)
	}

	widgets := existing.GetWidgets()
	colorGrid := dashboards.NewApiColorGridWidget("Color Grid")
	colorGrid.SetTitle("Unmanaged Color Grid Widget")
	widgets = append(widgets, dashboards.ApiColorGridWidgetAsApiWidget(colorGrid))

	update := dashboards.Dashboard{}
	update.SetTitle(existing.GetTitle())
	update.SetDescription(existing.GetDescription())
	update.SetIsPrivate(existing.GetIsPrivate())
	if ts, ok := existing.GetDefaultTimespanOk(); ok && ts != nil {
		update.SetDefaultTimespan(dashboards.DefaultTimespan{Duration: ts.Duration})
	}
	update.SetWidgets(widgets)

	updateReq := api.UpdateDashboard(dashboardID).Dashboard(update)
	updateReq = SetAidFromContext(ctx, updateReq)
	if _, _, err := updateReq.Execute(); err != nil {
		t.Fatalf("failed to add unmanaged widget to dashboard %s: %v", dashboardID, err)
	}
}

// checkDashboardTotalWidgetCount verifies the dashboard has exactly the expected
// number of widgets by querying the API directly (since unmanaged widgets are
// filtered from Terraform state).
func checkDashboardTotalWidgetCount(dashboardID string, expected int) error {
	api := (*dashboards.DashboardsAPIService)(&testClient.Common)
	ctx := testClient.GetConfig().Context

	getReq := api.GetDashboard(dashboardID)
	getReq = SetAidFromContext(ctx, getReq)
	resp, _, err := getReq.Execute()
	if err != nil {
		return fmt.Errorf("failed to GET dashboard %s: %w", dashboardID, err)
	}

	actual := len(resp.GetWidgets())
	if actual != expected {
		return fmt.Errorf("expected %d total widgets on dashboard, got %d", expected, actual)
	}
	return nil
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
