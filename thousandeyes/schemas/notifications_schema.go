package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
)

// Structs used for mapping
var _ = alerts.AlertNotification{}

var notifications = &schema.Schema{
	Type:        schema.TypeSet,
	Description: "The list of notifications for the alert rule.",
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeSet,
				Description: "The email notification.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"message": {
							Type:        schema.TypeString,
							Description: "The contents of the email, as a string.",
							Optional:    true,
						},
						"recipient": {
							Type:        schema.TypeSet,
							Description: "The email addresses to send the notification to.",
							Optional:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"third_party": {
				Type:        schema.TypeSet,
				Description: "Third party notification.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_id": {
							Type:        schema.TypeString,
							Description: "The integration ID, as a string.",
							Required:    true,
						},
						"integration_type": {
							Type:        schema.TypeString,
							Description: "The integration type, as a string.",
							Required:    true,
						},
					},
				},
			},
			"webhook": {
				Type:        schema.TypeSet,
				Description: "Webhook notification.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_id": {
							Type:        schema.TypeString,
							Description: "The integration ID, as a string.",
							Required:    true,
						},
						"integration_type": {
							Type:        schema.TypeString,
							Description: "The integration type, as a string.",
							Required:    true,
						},
						"integration_name": {
							Type:        schema.TypeString,
							Description: "Name of the integration, configured by the user.",
							Required:    true,
						},
						"target": {
							Type:        schema.TypeString,
							Description: "Webhook target URL.",
							Required:    true,
						},
					},
				},
			},
			"custom_webhook": {
				Type:        schema.TypeSet,
				Description: "Webhook notification.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"integration_id": {
							Type:        schema.TypeString,
							Description: "The integration ID, as a string.",
							Required:    true,
						},
						"integration_type": {
							Type:        schema.TypeString,
							Description: "The integration type, as a string.",
							Required:    true,
						},
						"integration_name": {
							Type:        schema.TypeString,
							Description: "Name of the integration, configured by the user.",
							Required:    true,
						},
						"target": {
							Type:        schema.TypeString,
							Description: "Webhook target URL.",
							Required:    true,
						},
					},
				},
			},
		},
	},
}
