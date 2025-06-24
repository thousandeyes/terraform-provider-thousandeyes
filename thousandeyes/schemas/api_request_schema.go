package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

var _ = tests.ApiRequest{}

var apiRequest = &schema.Schema{
	Type:        schema.TypeList,
	Description: "List of API requests",
	Required:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// assertions
			"assertions": {
				Type:        schema.TypeSet,
				Description: "List of assertion objects.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "[status-code, response-body] Name of assertion.",
							Optional:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"status-code",
								"response-body",
							}, false),
						},
						"operator": {
							Type:        schema.TypeString,
							Description: "[is, is-not, includes, not-includes] Assertion operator",
							Optional:    true,
							ValidateFunc: validation.StringInSlice([]string{
								"is",
								"is-not",
								"includes",
								"not-includes",
							}, false),
						},
						"value": {
							Type:        schema.TypeString,
							Description: "The value of the assertion. If name = `status-code`, the status code to assert. If name = `response-body`, the lookup value to assert.",
							Optional:    true,
						},
					},
				},
			},
			// authType
			"auth_type": {
				Type:        schema.TypeString,
				Description: "[none, basic, bearer-token, oauth2] Auth type will override the Authorization request header.",
				Optional:    true,
				Default:     "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"basic",
					"bearer-token",
					"oauth2",
				}, false),
			},
			// bearerToken
			"bearer_token": {
				Type:        schema.TypeString,
				Description: "Enter the OAuth body for the HTTP POST request in this field when using OAuth as the authentication mechanism. No special escaping is required. If content is provided in the post body, the `requestMethod` is automatically set to POST.",
				Optional:    true,
				Sensitive:   true,
			},
			// body
			"body": {
				Type:        schema.TypeString,
				Description: "POST/PUT request body. Must be in JSON format.",
				Optional:    true,
			},
			// clientAuthentication
			"client_authentication": {
				Type:        schema.TypeString,
				Description: "[basic-auth-header, in-body] The OAuth2 client authentication location type.",
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"basic-auth-header",
					"in-body",
				}, false),
			},
			// clientId
			"client_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The application ID used when `authType` is set to \"oauth2\".",
			},
			// clientSecret
			"client_secret": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The private client secret used when `authType` is set to \"oauth2\".",
			},
			// collectApiResponse
			"collect_api_response": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Set to `true` if API response body should be collected and saved. Set to `false` if API response body should not be saved.",
			},
			// headers
			"headers": {
				Type:        schema.TypeSet,
				Description: "An array of label IDs used to assign specific Endpoint Agents to the test (get `id` from `/endpoint/labels`). This is applicable when `alertGroupType` is `browser-session`.",
				Optional:    true,
				Sensitive:   true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request header key.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Request header value. Supports variables `{{variableName}}`.",
						},
					},
				},
			},
			// method
			"method": {
				Type:        schema.TypeString,
				Description: "[get, post, put, delete, patch] Request method.",
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"get",
					"post",
					"put",
					"delete",
					"patch",
				}, false),
			},
			// name
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "API step name, must be unique.",
			},
			// password
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The password if `authType = basic`.",
				Sensitive:   true,
			},
			// scope
			"scope": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Application-specific scope values for the access token when `authType` is \"oauth2\".",
			},
			// tokenUrl
			"token_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The endpoint used to request the access token when `authType` is \"oauth2\".",
			},
			// url
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Request url. Supports variables in the format `{{variableName}}`.",
			},
			// username
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The username if `authType = basic`.",
			},
			// variables
			"variables": {
				Type:        schema.TypeSet,
				Description: "Array of API post request variable objects.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Variable name.",
						},
						"value": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The JSON path of data within the Response Body to assign to this variable.",
						},
					},
				},
			},
			// verifyCertificate
			"verify_certificate": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Ignore or acknowledge certificate errors. Set to `false` to ignore certificate errors.",
			},
			// waitTimeMs
			"wait_time_ms": {
				Type:        schema.TypeInt,
				Description: "Post request delay before executing the next API requests, in milliseconds.",
				Optional:    true,
			},
		},
	},
}
