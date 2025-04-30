package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
)

var _ = administrative.UserRequest{}
var _ = administrative.UserDetail{}

var UserSchema = map[string]*schema.Schema{
	// uid
	"uid": {
		Type:        schema.TypeString,
		Description: "The unique ID of the user.",
		Computed:    true,
	},
	// name
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the user.",
		Required:    true,
	},
	// email
	"email": {
		Type:        schema.TypeString,
		Description: "The email of the user.",
		Optional:    true,
	},
	// loginAccountGroupId
	"login_account_group_id": {
		Type:        schema.TypeString,
		Description: "Unique ID of the login account group.",
		Optional:    true,
	},
	// dateRegistered
	"date_registered": {
		Type:        schema.TypeString,
		Description: "UTC date the user registered their account (ISO date-time format).",
		Computed:    true,
	},
	// lastLogin
	"last_login": {
		Type:        schema.TypeString,
		Description: "UTC date the user registered their account (ISO date-time format).",
		Computed:    true,
	},
	// accountGroupRoles
	"account_group_roles": {
		Type:        schema.TypeSet,
		Description: "Unique IDs representing the roles.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account_group_id": {
					Type:        schema.TypeString,
					Description: "Unique ID of the account group.",
					Required:    true,
				},
				"role_ids": {
					Type:        schema.TypeSet,
					Description: "Unique role IDs.",
					Optional:    true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
		},
	},
	// allAccountGroupRoleIds
	"all_account_group_role_ids": {
		Type:        schema.TypeSet,
		Description: "Unique IDs representing the roles.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// link (_links.self.href)
	"link": {
		Type:        schema.TypeString,
		Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
		Computed:    true,
	},
}
