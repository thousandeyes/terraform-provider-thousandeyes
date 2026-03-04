package schemas

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

var ConnectorAssignmentSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The ID of this resource. Matches `webhook_operation_id`.",
		Computed:    true,
	},
	"webhook_operation_id": {
		Type:        schema.TypeString,
		Description: "The ID of the webhook operation whose connector assignments are managed by this resource.",
		Required:    true,
		ForceNew:    true,
	},
	"connector_ids": {
		Type:        schema.TypeSet,
		Description: "The connector IDs assigned to the webhook operation. This list is authoritative and replaces all existing assignments on apply. The current API supports one connector per webhook operation.",
		Required:    true,
		MaxItems:    1,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
