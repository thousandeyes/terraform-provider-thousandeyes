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
		Description: "Tag color.",
		Optional:    true,
	},
	// createDate
	"create_date": {
		Type:        schema.TypeString,
		Description: "Tag creation date.",
		Computed:    true,
	},
	// icon
	"icon": {
		Type:        schema.TypeString,
		Description: "Tag icon.",
		Optional:    true,
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
		Description: "[test, dashboard, endpoint-test, v-agent] The object type associated with the tag.",
		Optional:    true,
		Required:    false,
		ValidateFunc: validation.StringInSlice([]string{
			"test",
			"dashboard",
			"endpoint-test",
			"v-agent",
		}, false),
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
