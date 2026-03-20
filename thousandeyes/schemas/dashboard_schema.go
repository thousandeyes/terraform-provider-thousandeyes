package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

var _ = dashboards.Dashboard{}
var _ = dashboards.ApiDashboard{}

var DashboardSchema = map[string]*schema.Schema{
	// id
	"id": {
		Type:        schema.TypeString,
		Description: "The tag ID.",
		Computed:    true,
	},
	// aid
	"aid": {
		Type:        schema.TypeString,
		Description: "Identifier for the account group associated with a dashboard.",
		Computed:    true,
	},
	// createdBy
	"dashboard_created_by": {
		Type:        schema.TypeString,
		Description: "Identifier for the user that created a dashboard.",
		Computed:    true,
	},
	// modifiedBy
	"dashboard_modified_by": {
		Type:        schema.TypeString,
		Description: "Identifier for the user that last modified a dashboard.",
		Computed:    true,
	},
	// modifiedDate
	"dashboard_modified_date": {
		Type:        schema.TypeString,
		Description: "UTC date/time when a dashboard was last modified (ISO date-time format).",
		Computed:    true,
	},
	// description
	"description": {
		Type:        schema.TypeString,
		Description: "A text description of the dashboard's purpose and functionality. This information assists users in understanding the dashboard but isn't displayed when viewing a dashboard.",
		Optional:    true,
	},
	// globalFilterId
	"global_filter_id": {
		Type:        schema.TypeString,
		Description: "Default global dashboard filter ID.",
		Optional:    true,
	},
	// title
	"title": {
		Type:        schema.TypeString,
		Description: "Title of a dashboard.",
		Required:    true,
	},
	// isGlobalOverride
	"is_global_override": {
		Type:        schema.TypeBool,
		Description: "When set to true, the defaultTimespan is used and overrides the widget's timespan. If set to false, the widget's timespan is used.",
		Optional:    true,
	},
	// isMigratedReport
	"is_migrated_report": {
		Type:        schema.TypeBool,
		Description: "True if this dashboard was previously a report.",
		Computed:    true,
	},
	// isPrivate
	"is_private": {
		Type:        schema.TypeBool,
		Description: "A dashboard can be viewed by other users in the account. If true, only the creator of the dashboard may view it. If false, all users in the same account may view it.",
		Optional:    true,
		Default:     false,
	},
	// defaultTimespan
	"default_timespan": {
		Type:        schema.TypeList,
		Description: "Defines the default time range displayed by the dashboard.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				// Duration
				"duration": {
					Type:          schema.TypeInt,
					Description:   "Duration of the timespan in seconds. Mutually exclusive with start/end.",
					Optional:      true,
					ConflictsWith: []string{"default_timespan.0.start", "default_timespan.0.end"},
				},
				// start
				"start": {
					Type:          schema.TypeString,
					Description:   "UTC start date of the timespan range (ISO date-time format). Mutually exclusive with duration.",
					Optional:      true,
					ConflictsWith: []string{"default_timespan.0.duration"},
				},
				// end
				"end": {
					Type:          schema.TypeString,
					Description:   "UTC end date of the timespan range (ISO date-time format). Mutually exclusive with duration.",
					Optional:      true,
					ConflictsWith: []string{"default_timespan.0.duration"},
				},
			},
		},
	},
}
