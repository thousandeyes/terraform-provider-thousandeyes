package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

// Structs used for mapping
var _ = tags.Tag{}
var _ = tags.TagInfo{}

var TagSchema = map[string]*schema.Schema{
	// accessType
	"access_type": {
		Type:        schema.TypeString,
		Description: "[all, partner, system] The access level of the tag. The access level impacts the visibility of the label in UI and the permissions to modify the label.",
		Optional:    true,
		Required:    false,
		Default:     "use-agent-policy",
		ValidateFunc: validation.StringInSlice([]string{
			"all",
			"partner",
			"system",
		}, false),
	},
	// aid
	"aid": {
		Type:        schema.TypeString, // integer in model
		Description: "The account group ID.",
		Computed:    true,
	},
	// color
	"color": {
		Type:        schema.TypeString,
		Description: "Tag color. When omitted, uses the API-provided value.",
		Optional:    true,
		Computed:    true,
	},
	// createDate
	"create_date": {
		Type:        schema.TypeString,
		Description: "Tag creation date.",
		Computed:    true,
	},
	// modifiedDate
	"modified_date": {
		Type:        schema.TypeString,
		Description: "The date and time the tag was last modified.",
		Computed:    true,
	},
	// builtIn
	"built_in": {
		Type:        schema.TypeBool,
		Description: "Indicates whether this is a built-in (system) tag.",
		Computed:    true,
	},
	// icon
	"icon": {
		Type:        schema.TypeString,
		Description: "Tag icon. When omitted, uses the API-provided value.",
		Optional:    true,
		Computed:    true,
	},
	// description
	"description": {
		Type:        schema.TypeString,
		Description: "The tag's description.",
		Optional:    true,
	},
	// id
	"id": {
		Type:        schema.TypeString,
		Description: "The tag ID.",
		Computed:    true,
	},
	// key
	"key": {
		Type:        schema.TypeString,
		Description: "The tags's key.",
		Optional:    true,
	},
	// legacyId
	"legacy_id": {
		Type:        schema.TypeInt, // float in model
		Description: "Legacy Id.",
		Computed:    true,
	},
	// objectType
	"object_type": {
		Type:        schema.TypeString,
		Description: "[test, dashboard, endpoint-test, v-agent, endpoint-agent] The object type associated with the tag. Endpoint Agent tags require type = dynamic.",
		Optional:    true,
		Required:    false,
		ValidateFunc: validation.StringInSlice([]string{
			"test",
			"dashboard",
			"endpoint-test",
			"v-agent",
			"endpoint-agent",
		}, false),
	},
	// type
	"type": {
		Type:        schema.TypeString,
		Description: "The nature of the tag: static or dynamic. Only endpoint-agent tags support dynamic assignment. Existing static tags can omit this field.",
		Optional:    true,
		Computed:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"static",
			"dynamic",
		}, false),
	},
	"match_type": {
		Type:         schema.TypeString,
		Description:  "How dynamic Endpoint Agent tag filters are combined: and requires all filters to match; or requires any filter to match. Required for dynamic tags.",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"and", "or"}, false),
	},
	"filters": {
		Type:        schema.TypeSet,
		Description: "Unordered filter criteria for dynamic Endpoint Agent tags. At least one filter is required for dynamic tags.",
		Optional:    true,
		MinItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:         schema.TypeString,
					Description:  "The built-in agent field or custom check-in metadata key to filter on. The API validates supported keys.",
					Required:     true,
					ValidateFunc: validation.StringIsNotEmpty,
				},
				"values": {
					Type:        schema.TypeSet,
					Description: "Unordered values to match, passed unchanged to the API. Their format and matching semantics depend on the filter key.",
					Required:    true,
					MinItems:    1,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
				"mode": {
					Type:         schema.TypeString,
					Description:  "Whether the agent's value must match a value (in) or must not match any value (not-in).",
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"in", "not-in"}, false),
				},
				"scope": {
					Type:         schema.TypeString,
					Description:  "Where the filter key comes from: default for built-in ThousandEyes agent fields, or custom for metadata supplied when an agent checks in.",
					Optional:     true,
					Default:      "default",
					ValidateFunc: validation.StringInSlice([]string{"default", "custom"}, false),
				},
			},
		},
	},
	// value
	"value": {
		Type:        schema.TypeString,
		Description: "The tag's value.",
		Optional:    true,
	},
	// link (_links.self.href)
	"link": {
		Type:        schema.TypeString,
		Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
		Computed:    true,
	},
}
