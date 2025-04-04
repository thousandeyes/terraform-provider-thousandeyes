package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

// Structs used for mapping
var _ = tests.OAuth{}

var oauth = &schema.Schema{
	Type:        schema.TypeSet,
	Description: "Use this only if you want to use OAuth as the authentication mechanism.",
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"config_id": {
				Type:        schema.TypeString,
				Description: "The ID of the OAuth configuration.",
				Required:    true,
			},
			"test_url": {
				Type:        schema.TypeString,
				Description: "Target for the test.",
				Required:    true,
			},
			"request_method": {
				Type:        schema.TypeString,
				Description: "[get, post, put, delete, patch, options, trace] Request method.",
				Optional:    true,
				ValidateFunc: validation.StringInSlice([]string{
					"get",
					"post",
					"put",
					"delete",
					"patch",
					"options",
					"trace",
				}, false),
			},
			"post_body": {
				Type:        schema.TypeString,
				Description: "Enter the OAuth body for the HTTP POST request in this field when using OAuth as the authentication mechanism. No special escaping is required. If content is provided in the post body, the `requestMethod` is automatically set to POST.",
				Optional:    true,
			},
			"headers": {
				Type:        schema.TypeString,
				Description: "Request headers used for OAuth.",
				Optional:    true,
				Sensitive:   true,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Description: "[none, basic, ntlm, kerberos, oauth] The HTTP authentication type. Defaults to 'none'.",
				Optional:    true,
				Default:     "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"basic",
					"ntlm",
					"kerberos",
					"oauth",
				}, false),
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAuth username",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAuth password",
				Sensitive:   true,
			},
		},
	},
}
