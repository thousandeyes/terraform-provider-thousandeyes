package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var ConnectorSchema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "The unique identifier of the connector.",
		Computed:    true,
	},
	"name": {
		Type:        schema.TypeString,
		Description: "The name of the connector.",
		Required:    true,
	},
	"target": {
		Type:        schema.TypeString,
		Description: "The target URL where webhook notifications are sent.",
		Required:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The connector type. Currently only 'generic' is supported.",
		Optional:    true,
		Default:     "generic",
		ValidateFunc: validation.StringInSlice([]string{
			"generic",
		}, false),
	},
	"headers": {
		Type:        schema.TypeList,
		Description: "Custom headers to include in webhook requests.",
		Optional:    true,
		Sensitive:   true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The header name.",
					Required:    true,
				},
				"value": {
					Type:        schema.TypeString,
					Description: "The header value. Note: values are obfuscated in responses.",
					Required:    true,
					Sensitive:   true,
				},
			},
		},
	},
	"authentication": {
		Type:        schema.TypeList,
		Description: "Authentication configuration for the connector.",
		Optional:    true,
		MaxItems:    1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:        schema.TypeString,
					Description: "The authentication type: basic, bearer-token, oauth-client-credentials, oauth-auth-code, or other-token.",
					Required:    true,
					ValidateFunc: validation.StringInSlice([]string{
						"basic",
						"bearer-token",
						"oauth-client-credentials",
						"oauth-auth-code",
						"other-token",
					}, false),
				},
				"username": {
					Type:        schema.TypeString,
					Description: "Username for basic authentication.",
					Optional:    true,
				},
				"password": {
					Type:        schema.TypeString,
					Description: "Password for basic authentication.",
					Optional:    true,
					Sensitive:   true,
				},
				"token": {
					Type:        schema.TypeString,
					Description: "Token for bearer-token, other-token, or OAuth authentication.",
					Optional:    true,
					Sensitive:   true,
				},
				"oauth_client_id": {
					Type:        schema.TypeString,
					Description: "OAuth client ID for OAuth authentication.",
					Optional:    true,
				},
				"oauth_client_secret": {
					Type:        schema.TypeString,
					Description: "OAuth client secret for OAuth authentication.",
					Optional:    true,
					Sensitive:   true,
				},
				"oauth_token_url": {
					Type:        schema.TypeString,
					Description: "OAuth token URL for OAuth authentication.",
					Optional:    true,
				},
				"oauth_auth_url": {
					Type:        schema.TypeString,
					Description: "OAuth authorization URL for OAuth auth-code authentication.",
					Optional:    true,
				},
				"code": {
					Type:        schema.TypeString,
					Description: "OAuth authorization code for OAuth auth-code authentication.",
					Optional:    true,
					Sensitive:   true,
				},
				"redirect_uri": {
					Type:        schema.TypeString,
					Description: "OAuth redirect URI for OAuth auth-code authentication.",
					Optional:    true,
				},
				"refresh_token": {
					Type:        schema.TypeString,
					Description: "OAuth refresh token for OAuth auth-code authentication.",
					Optional:    true,
					Sensitive:   true,
				},
			},
		},
	},
	"last_modified_date": {
		Type:        schema.TypeString,
		Description: "The date when the connector was last modified.",
		Computed:    true,
	},
}
