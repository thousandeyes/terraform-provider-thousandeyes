package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

// Structs used for mapping
var _ = tags.TagAssignment{}
var _ = tags.BulkTagAssignment{}
var _ = tags.Assignment{}

var TagAssignmentSchema = map[string]*schema.Schema{
	// tagId
	"tag_id": {
		Type:        schema.TypeString,
		Description: "The ID of the tag to assign.",
		Required:    true,
		ForceNew:    true,
	},
	// assignments
	"assignments": {
		Description: "",
		Required:    true,
		ForceNew:    true,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:        schema.TypeString,
					Description: "Object Id.",
					Required:    true,
				},
				"type": {
					Type:        schema.TypeString,
					Description: "[test, v-agent, endpoint-test, dashboard] Type of assignment.",
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						"test",
						"v-agent",
						"endpoint-test",
						"dashboard",
					}, false),
				},
			},
		},
	},
}
