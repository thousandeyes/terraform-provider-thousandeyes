package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

var _ = connectors.WebhookOperation{}
var _ = connectors.Header{}

var WebhookOperationSchema = map[string]*schema.Schema{
	// id
	"id": {
		Type:        schema.TypeString,
		Description: "The operation ID.",
		Computed:    true,
	},
	// name
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the webhook operation.",
		Required:    true,
	},
	// enabled
	"enabled": {
		Type:        schema.TypeBool,
		Description: "Whether the webhook operation is enabled.",
		Optional:    true,
	},
	// category
	"category": {
		Type:        schema.TypeString,
		Description: "[alerts, recommendations, traffic-monitoring] The category of the webhook operation.",
		Required:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"alerts",
			"recommendations",
			"traffic-monitoring",
		}, false),
	},
	// status
	"status": {
		Type:        schema.TypeString,
		Description: "[pending, connected, failing, unverified] The status of the webhook operation.",
		Required:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"pending",
			"connected",
			"failing",
			"unverified",
		}, false),
	},
	// path
	"path": {
		Type:        schema.TypeString,
		Description: "The custom path for the webhook endpoint.",
		Optional:    true,
	},
	// payload
	"payload": {
		Type:        schema.TypeString,
		Description: "Handlebars template for the payload. Must be a valid JSON string with Handlebars variables.",
		Optional:    true,
	},
	// headers
	"headers": {
		Type:        schema.TypeList,
		Description: "Custom headers to include in the webhook request.",
		Optional:    true,
		Sensitive:   true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "Header name.",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "Header value. Note that this value is obfuscated in the response.",
					Required:    true,
					Sensitive:   true,
				},
			},
		},
	},
	// query_params
	"query_params": {
		Type:        schema.TypeString,
		Description: "Handlebars template for the query params. Must compile into a proper JSON object where each object property will define the query param name and the object property value define the corresponding query param value.",
		Optional:    true,
	},
	// type
	"type": {
		Type:        schema.TypeString,
		Description: "The type of operation. Always 'webhook' for webhook operations.",
		Optional:    true,
		Default:     "webhook",
		ValidateFunc: validation.StringInSlice([]string{
			"webhook",
		}, false),
	},
	// link (_links.self.href)
	"link": {
		Type:        schema.TypeString,
		Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
		Computed:    true,
	},
}
