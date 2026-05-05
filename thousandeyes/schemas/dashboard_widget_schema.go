package schemas

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type DashboardWidgetSchemaType map[string]*schema.Schema

var deprecatedDashboardFilterProperties = map[string]struct{}{
	"TEST_LABEL":          {},
	"AGENT_LABEL":         {},
	"ENDPOINT_TEST_LABEL": {},
}

func validateDashboardFilterProperty(v interface{}, path string) ([]string, []error) {
	value, ok := v.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected %s to be a string", path)}
	}

	if _, isDeprecated := deprecatedDashboardFilterProperties[value]; isDeprecated {
		return nil, []error{fmt.Errorf("%s must not use deprecated label filter property %q", path, value)}
	}

	return nil, nil
}

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
			"Table",
			"Test Table",
			"Bar Chart: Stacked",
			"Bar Chart: Grouped",
			"Color Grid",
			"Alert List",
			"List",
			"Number",
			"Multi Metric Table",
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
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"value": {
					Type:        schema.TypeInt,
					Description: "Time value.",
					Optional:    true,
					Computed:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Time unit.",
					Optional:    true,
					Computed:    true,
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
		Type:        schema.TypeSet,
		Description: "Filters applied to the widget. Each filter specifies a property and list of values. Order is not significant.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"property": {
					Type:         schema.TypeString,
					Description:  "Filter property (e.g., 'TEST', 'AGENT', 'ENDPOINT_MACHINE_ID', 'MONITOR').",
					Required:     true,
					ValidateFunc: validateDashboardFilterProperty,
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
				"min_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether min_scale was configured.",
					Computed:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value.",
					Optional:    true,
				},
				"max_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether max_scale was configured.",
					Computed:    true,
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
					Computed:    true,
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
				"min_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether min_scale was configured.",
					Computed:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether max_scale was configured.",
					Computed:    true,
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
					Computed:    true,
				},
				"is_timeseries_one_chart_per_line": {
					Type:        schema.TypeBool,
					Description: "Displays a separate chart for each line if set to true.",
					Optional:    true,
					Computed:    true,
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
				"min_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether min_scale was configured.",
					Computed:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether max_scale was configured.",
					Computed:    true,
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
				"min_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether min_scale was configured.",
					Computed:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value for the Y-axis.",
					Optional:    true,
				},
				"max_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether max_scale was configured.",
					Computed:    true,
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

	"table_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Table widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"compare_to_previous_value": {
					Type:        schema.TypeBool,
					Description: "Enables comparison of the current metric value with the previous value.",
					Optional:    true,
					Computed:    true,
				},
				"row_group_by": {
					Type:        schema.TypeString,
					Description: "Group rows by property.",
					Optional:    true,
				},
				"column_group_by": {
					Type:        schema.TypeString,
					Description: "Group columns by property.",
					Optional:    true,
				},
				"limit": {
					Type:        schema.TypeInt,
					Description: "Maximum number of rows to display.",
					Optional:    true,
				},
			},
		},
	},

	"test_table_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Test Table widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"filter": {
					Type:        schema.TypeList,
					Description: "Include filter configuration.",
					Optional:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:        schema.TypeString,
								Description: "How to combine the filters.",
								Optional:    true,
								ValidateFunc: validation.StringInSlice([]string{
									"all",
									"any",
								}, false),
							},
							"filters": {
								Type:        schema.TypeList,
								Description: "Filter terms.",
								Optional:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"key": {
											Type:        schema.TypeString,
											Description: "Filter key.",
											Required:    true,
											ValidateFunc: validation.StringInSlice([]string{
												"Anything",
												"Test Name",
												"Target",
												"Test ID",
												"Test type",
												"Tag ID",
											}, false),
										},
										"value": {
											Type:        schema.TypeString,
											Description: "Filter value.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
				"exclude": {
					Type:        schema.TypeList,
					Description: "Exclude filter configuration.",
					Optional:    true,
					MaxItems:    1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"type": {
								Type:        schema.TypeString,
								Description: "How to combine the filters.",
								Optional:    true,
								ValidateFunc: validation.StringInSlice([]string{
									"all",
									"any",
								}, false),
							},
							"filters": {
								Type:        schema.TypeList,
								Description: "Filter terms.",
								Optional:    true,
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"key": {
											Type:        schema.TypeString,
											Description: "Filter key.",
											Required:    true,
											ValidateFunc: validation.StringInSlice([]string{
												"Anything",
												"Test Name",
												"Target",
												"Test ID",
												"Test type",
												"Tag ID",
											}, false),
										},
										"value": {
											Type:        schema.TypeString,
											Description: "Filter value.",
											Required:    true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},

	"stacked_bar_chart_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Bar Chart: Stacked widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"axis_group_by": {
					Type:        schema.TypeString,
					Description: "Axis grouping property.",
					Optional:    true,
				},
				"limit": {
					Type:        schema.TypeInt,
					Description: "Maximum number of bars to display.",
					Optional:    true,
				},
				"show_labels": {
					Type:        schema.TypeBool,
					Description: "Displays labels on each bar.",
					Optional:    true,
					Computed:    true,
				},
				"is_horizontal_bar_chart": {
					Type:        schema.TypeBool,
					Description: "Displays bars horizontally when set to true.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	},

	"grouped_bar_chart_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Bar Chart: Grouped widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group_by": {
					Type:        schema.TypeString,
					Description: "Group bars by property.",
					Optional:    true,
				},
				"axis_group_by": {
					Type:        schema.TypeString,
					Description: "Axis grouping property.",
					Optional:    true,
				},
				"limit": {
					Type:        schema.TypeInt,
					Description: "Maximum number of bars to display.",
					Optional:    true,
				},
				"show_labels": {
					Type:        schema.TypeBool,
					Description: "Displays labels on each bar.",
					Optional:    true,
					Computed:    true,
				},
				"is_horizontal_bar_chart": {
					Type:        schema.TypeBool,
					Description: "Displays bars horizontally when set to true.",
					Optional:    true,
					Computed:    true,
				},
			},
		},
	},

	"color_grid_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Color Grid widgets.",
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
				"min_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether min_scale was configured.",
					Computed:    true,
				},
				"max_scale": {
					Type:        schema.TypeFloat,
					Description: "Maximum scale value.",
					Optional:    true,
				},
				"max_scale_configured": {
					Type:        schema.TypeBool,
					Description: "Internal marker indicating whether max_scale was configured.",
					Computed:    true,
				},
				"unit": {
					Type:        schema.TypeString,
					Description: "Unit for the scale.",
					Optional:    true,
				},
				"cards": {
					Type:        schema.TypeString,
					Description: "Aggregate property used for cards.",
					Optional:    true,
				},
				"group_cards_by": {
					Type:        schema.TypeString,
					Description: "Aggregate property used to group cards.",
					Optional:    true,
				},
				"columns": {
					Type:        schema.TypeInt,
					Description: "Number of columns.",
					Optional:    true,
					Computed:    true,
				},
				"limit": {
					Type:        schema.TypeInt,
					Description: "Maximum number of cards to display.",
					Optional:    true,
				},
			},
		},
	},

	"alert_list_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Alert List widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"alert_types": {
					Type:        schema.TypeSet,
					Description: "Alert types to include. Empty means all alert types.",
					Optional:    true,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
				},
				"limit_to": {
					Type:        schema.TypeInt,
					Description: "Maximum number of alerts to display.",
					Optional:    true,
					Computed:    true,
				},
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

	// Type-specific: Number Cards (for "Number" type)
	"number_cards": {
		Type:        schema.TypeList,
		Description: "List of number cards within a Number widget. Each card can have its own data source, metric, and measure.",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: NumberCardSchema,
		},
	},

	// Type-specific: Multi Metric Table configuration
	"multi_metric_table_config": {
		Type:        schema.TypeList,
		Description: "Configuration for Multi Metric Table widgets.",
		Optional:    true,
		Computed:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"compare_to_previous_value": {
					Type:        schema.TypeBool,
					Description: "Enables comparison of the current metric value with the previous value.",
					Optional:    true,
					Computed:    true,
				},
				"row_group_by": {
					Type:        schema.TypeString,
					Description: "Property to group rows by.",
					Optional:    true,
					Computed:    true,
				},
				"limit": {
					Type:        schema.TypeInt,
					Description: "Maximum number of rows displayed.",
					Optional:    true,
				},
			},
		},
	},

	// Type-specific: Multi Metric Columns (for "Multi Metric Table" type)
	"multi_metric_columns": {
		Type:        schema.TypeList,
		Description: "List of columns within a Multi Metric Table widget. Each column has its own data source, metric, and measure.",
		Optional:    true,
		Computed:    true,
		Elem: &schema.Resource{
			Schema: MultiMetricColumnSchema,
		},
	},
}

var NumberCardSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "Identifier of the number card.",
		Computed:    true,
	},
	"description": {
		Type:        schema.TypeString,
		Description: "Description of the number card.",
		Optional:    true,
	},
	"min_scale": {
		Type:        schema.TypeFloat,
		Description: "Minimum scale configured for the card.",
		Optional:    true,
	},
	"min_scale_configured": {
		Type:        schema.TypeBool,
		Description: "Internal marker indicating whether min_scale was configured.",
		Computed:    true,
	},
	"max_scale": {
		Type:        schema.TypeFloat,
		Description: "Maximum scale configured for the card.",
		Optional:    true,
	},
	"max_scale_configured": {
		Type:        schema.TypeBool,
		Description: "Internal marker indicating whether max_scale was configured.",
		Computed:    true,
	},
	"unit": {
		Type:        schema.TypeString,
		Description: "Unit for the scale.",
		Optional:    true,
	},
	"compare_to_previous_value": {
		Type:        schema.TypeBool,
		Description: "Enables comparison with the previous metric value.",
		Optional:    true,
		Computed:    true,
	},
	"fixed_timespan": {
		Type:        schema.TypeList,
		Description: "Fixed timespan for the card.",
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
		Description: "Data source for the card.",
		Optional:    true,
		Computed:    true,
	},
	"metric_group": {
		Type:        schema.TypeString,
		Description: "Metric group for the card.",
		Optional:    true,
	},
	"direction": {
		Type:        schema.TypeString,
		Description: "Direction for the metric.",
		Optional:    true,
		Computed:    true,
	},
	"metric": {
		Type:        schema.TypeString,
		Description: "Metric for the card.",
		Optional:    true,
	},
	"measure": {
		Type:        schema.TypeList,
		Description: "Measure configuration for the card.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Description: "Measure type.",
					Optional:    true,
				},
				"percentile_value": {
					Type:        schema.TypeFloat,
					Description: "Percentile value when type is NTH_PERCENTILE.",
					Optional:    true,
				},
			},
		},
	},
	"filter": {
		Type:        schema.TypeSet,
		Description: "Filters applied to the card. Order is not significant.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"property": {
					Type:         schema.TypeString,
					Description:  "Filter property.",
					Required:     true,
					ValidateFunc: validateDashboardFilterProperty,
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
}
var MultiMetricColumnSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "Identifier of the column.",
		Computed:    true,
	},
	"data_source": {
		Type:        schema.TypeString,
		Description: "Data source for the column.",
		Optional:    true,
		Computed:    true,
	},
	"metric_group": {
		Type:        schema.TypeString,
		Description: "Metric group for the column.",
		Optional:    true,
	},
	"direction": {
		Type:        schema.TypeString,
		Description: "Direction for the metric (e.g., TO_TARGET, FROM_TARGET). Only applicable to certain data sources.",
		Optional:    true,
		Computed:    true,
	},
	"metric": {
		Type:        schema.TypeString,
		Description: "Metric for the column.",
		Optional:    true,
	},
	"measure": {
		Type:        schema.TypeList,
		Description: "Measure configuration for the column.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Description: "Measure type.",
					Optional:    true,
				},
				"percentile_value": {
					Type:        schema.TypeFloat,
					Description: "Percentile value when type is NTH_PERCENTILE.",
					Optional:    true,
				},
			},
		},
	},
}
