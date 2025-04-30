package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
)

var _ = administrative.RoleRequestBody{}
var _ = administrative.RoleDetail{}

var RoleSchema = map[string]*schema.Schema{
	// roleId
	"role_id": {
		Type:        schema.TypeString,
		Description: "The unique ID of the role.",
		Computed:    true,
	},
	// name
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the role.",
		Optional:    true,
	},
	// permissions
	"permissions": {
		Type:        schema.TypeSet,
		Description: "Contains list of test permission IDs",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// isBuiltin
	"is_builtin": {
		Type:        schema.TypeBool,
		Description: "Flag indicating if the role is built-in (Account Admin, Organization Admin, Regular User).",
		Computed:    true,
	},
	// link (_links.self.href)
	"link": {
		Type:        schema.TypeString,
		Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
		Computed:    true,
	},
}
