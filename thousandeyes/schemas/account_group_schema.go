package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
)

var _ = administrative.AccountGroupRequest{}
var _ = administrative.AccountGroupDetail{}

var AccountGroupSchema = map[string]*schema.Schema{
	// aid
	"aid": {
		Type:        schema.TypeString,
		Description: "The unique ID of the account group.",
		Computed:    true,
	},
	// accountGroupName
	"account_group_name": {
		Type:        schema.TypeString,
		Description: "The name of the account group.",
		Required:    true,
	},
	// isCurrentAccountGroup
	"is_current_account_group": {
		Type:        schema.TypeBool,
		Description: "Indicates whether the requested aid is the context of the current account.",
		Computed:    true,
	},
	// isDefaultAccountGroup
	"is_default_account_group": {
		Type:        schema.TypeBool,
		Description: "Indicates whether the aid is the default one for the requesting user.",
		Computed:    true,
	},
	// organizationName
	"organization_name": {
		Type:        schema.TypeString,
		Description: "The name of the organization associated with the account group.",
		Computed:    true,
	},
	// orgId
	"org_id": {
		Type:        schema.TypeString,
		Description: "The ID for the organization associated with the account group.",
		Computed:    true,
	},
	// accountToken
	"account_token": {
		Type:        schema.TypeString,
		Description: "The account group token is an alphanumeric string used to bind an Enterprise Agent to a specific account group. This token is not a password that must be kept secret.",
		Computed:    true,
	},
	// users
	"users": {
		Type:        schema.TypeSet,
		Description: "",
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// agents
	"agents": {
		Type:        schema.TypeSet,
		Description: "To grant access to enterprise agents, specify the agent list. Note that this is not an additive list - the full list must be specified if changing access to agents.",
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
