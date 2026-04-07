package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DashboardWidgetSchemaType map[string]*schema.Schema

var DashboardWidgetSchema = DashboardWidgetSchemaType{
	"type": {
		Type:        schema.TypeString,
		Description: "Widget type.",
		Required:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"Map",
			"Agent Status",
			"Time Series: Line",
			"Time Series: Stacked Area",
			"Pie Chart",
			"Box and Whiskers",
			// "List" is temporarily disabled until the API returns valid
			// sortDirection values (CP-2702). The SDK cannot deserialize the
			// ASC/DESC values the API currently sends, breaking refresh.
		}, false),
	},

	// Common fields
	"id": {
		Type:        schema.TypeString,
		Description: "Identifier of the widget.",
		Computed:    true,
	},
	"title": {
		Type:        schema.TypeString,
		Description: "Title of the widget.",
		Optional:    true,
	},
	"visual_mode": {
		Type:        schema.TypeString,
		Description: "Visual mode of the widget (e.g., 'Full', 'Half screen').",
		Optional:    true,
		Computed:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"Full",
			"Half screen",
		}, false),
	},
	"embed_url": {
		Type:        schema.TypeString,
		Description: "When isEmbedded is set to true, an embedUrl is provided.",
		Computed:    true,
	},
	"is_embedded": {
		Type:        schema.TypeBool,
		Description: "Indicates whether the widget is marked as embedded.",
		Computed:    true,
	},
	"metric_group": {
		Type:        schema.TypeString,
		Description: "Metric group for the widget.",
		Optional:    true,
	},
	"direction": {
		Type:        schema.TypeString,
		Description: "Direction for the metric (e.g., 'To Target', 'From Target', 'Bidirectional').",
		Optional:    true,
		Computed:    true,
	},
	"metric": {
		Type:        schema.TypeString,
		Description: "Metric for the widget.",
		Optional:    true,
	},
	"measure": {
		Type:        schema.TypeList,
		Description: "Measure configuration for the widget.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Description: "Measure type (e.g., 'MEAN', 'MEDIAN', 'MAXIMUM', 'MINIMUM', 'NTH_PERCENTILE').",
					Optional:    true,
				},
				"percentile_value": {
					Type:        schema.TypeFloat,
					Description: "The percentile value to use when type is NTH_PERCENTILE.",
					Optional:    true,
				},
			},
		},
	},
	"fixed_timespan": {
		Type:        schema.TypeList,
		Description: "Fixed timespan for the widget.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:        schema.TypeInt,
					Description: "Time value.",
					Optional:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Time unit.",
					Optional:    true,
				},
			},
		},
	},
	"should_exclude_alert_suppression_windows": {
		Type:        schema.TypeBool,
		Description: "Excludes alert suppression window data if set to true.",
		Optional:    true,
		Computed:    true,
	},
	"data_source": {
		Type:        schema.TypeString,
		Description: "Data source for the widget.",
		Optional:    true,
		Computed:    true,
	},
	"filter": {
		Type:        schema.TypeList,
		Description: "Filters applied to the widget. Each filter specifies a property and list of values.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"property": {
					Type:        schema.TypeString,
					Description: "Filter property (e.g., 'TEST', 'AGENT', 'ENDPOINT_MACHINE_ID', 'MONITOR').",
					Required:    true,
				},
				"values": {
					Type:        schema.TypeSet,
					Description: "Set of filter values (IDs). Order is not significant.",
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	},

	// Type-specific: GeoMap configuration (for "Map" type)
	"geo_map_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Map widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_scale": {
					Type:        schema.TypeFloat,
					Description: "Minimum scale value.",
					Optional:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value.",
					Optional:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Unit for the scale.",
					Optional:    true,
				},
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group by property.",
					Optional:    true,
					Computed:    true,
				},
				"is_geo_map_per_test": {
					Type:        schema.TypeBool,
					Description: "Show one map per test.",
					Optional:    true,
				},
			},
		},
	},

	// Type-specific: Agent Status configuration (for "Agent Status" type)
	"agent_status_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Agent Status widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"show": {
					Type:        schema.TypeString,
					Description: "What to show (e.g., 'All Agents', 'Online Agents', 'Offline Agents').",
					Optional:    true,
				},
				"agent_type": {
					Type:        schema.TypeString,
					Description: "Type of agent (e.g., 'Enterprise', 'Cloud').",
					Optional:    true,
				},
			},
		},
	},

	// Type-specific: Timeseries configuration (for "Time Series: Line" type)
	"timeseries_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Time Series: Line widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_scale": {
					Type:        schema.TypeFloat,
					Description: "Minimum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Unit for the Y-axis scale.",
					Optional:    true,
				},
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group by property.",
					Optional:    true,
					Computed:    true,
				},
				"show_timeseries_overall_baseline": {
					Type:        schema.TypeBool,
					Description: "Displays the overall baseline if set to true.",
					Optional:    true,
				},
				"is_timeseries_one_chart_per_line": {
					Type:        schema.TypeBool,
					Description: "Displays a separate chart for each line if set to true.",
					Optional:    true,
				},
			},
		},
	},

	// Type-specific: Stacked Area configuration (for "Time Series: Stacked Area" type)
	"stacked_area_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Time Series: Stacked Area widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_scale": {
					Type:        schema.TypeFloat,
					Description: "Minimum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Unit for the Y-axis scale.",
					Optional:    true,
				},
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group by property.",
					Required:    true,
				},
			},
		},
	},

	// Type-specific: Pie Chart configuration (for "Pie Chart" type)
	"pie_chart_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Pie Chart widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group by property.",
					Required:    true,
				},
			},
		},
	},

	// Type-specific: Box and Whiskers configuration (for "Box and Whiskers" type)
	"box_and_whiskers_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Box and Whiskers widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"min_scale": {
					Type:        schema.TypeFloat,
					Description: "Minimum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Unit for the Y-axis scale.",
					Optional:    true,
				},
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group by property.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	},

	// Type-specific: List configuration (for "List" type)
	"list_config": {
		Type:        schema.TypeList,
		Description: "Configuration for List widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"active_within_value": {
					Type:        schema.TypeInt,
					Description: "Timespan value for active within filter.",
					Optional:    true,
				},
				"active_within_unit": {
					Type:        schema.TypeString,
					Description: "Timespan unit for active within filter.",
					Optional:    true,
				},
			},
		},
	},
}
