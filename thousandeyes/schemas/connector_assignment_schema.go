package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ConnectorAssignmentSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The ID of this resource. Matches `connector_id`.",
		Computed:    true,
	},
	"connector_id": {
		Type:        schema.TypeString,
		Description: "The ID of the connector whose webhook operation assignments are managed by this resource.",
		Required:    true,
		ForceNew:    true,
	},
	"operation_ids": {
		Type:        schema.TypeSet,
		Description: "The webhook operation IDs assigned to the connector. This list is authoritative and replaces all existing assignments on apply.",
		Required:    true,
		Elem: &schema.Schema{
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
	},
}
