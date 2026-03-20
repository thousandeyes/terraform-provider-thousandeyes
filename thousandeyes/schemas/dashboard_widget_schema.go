package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DashboardWidgetSchemaType map[string]*schema.Schema

var DashboardWidgetSchema = DashboardWidgetSchemaType{
	// Widget type - currently only supports Map and Agent Status
	"type": {
		Type:        schema.TypeString,
		Description: "Widget type.",
		Required:    true,
		ForceNew:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"Map",
			"Agent Status",
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
		Description: "Visual mode of the widget (e.g., 'Full', 'Half').",
		Optional:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"Full",
			"Half",
		}, false),
	},
	"embed_url": {
		Type:        schema.TypeString,
		Description: "When isEmbedded is set to true, an embedUrl is provided.",
		Computed:    true,
	},
	"is_embedded": {
		Type:        schema.TypeBool,
		Description: "Set to true if widget is marked as embedded.",
		Optional:    true,
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
	},
	"data_source": {
		Type:        schema.TypeString,
		Description: "Data source for the widget.",
		Optional:    true,
	},
	"filters": {
		Type:        schema.TypeMap,
		Description: "Filters applied to the widget.",
		Optional:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
	},

	// Type-specific: GeoMap configuration (for "Map" type)
	"geo_map_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Map widgets.",
		Optional:    true,
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
}
