package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

// Structs used for mapping
var _ = tests.TestSipCredentials{}

var targetSipCredentials = &schema.Schema{
	Type:        schema.TypeSet,
	Description: "Target SIP credentials",
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			// authUser
			"auth_user": {
				Type:        schema.TypeString,
				Description: "The username for authentication with the SIP server.",
				Required:    true,
			},
			// password
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "Password for Basic/NTLM authentication.",
			},
			// port
			"port": {
				Type:         schema.TypeInt,
				Description:  "The target port.",
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
			},
			// protocol
			"protocol": {
				Type:         schema.TypeString,
				Description:  "[tcp, tls, or udp] The transport layer for SIP communication. Can be either TCP, TLS (TLS over TCP), or UDP, and defaults to tcp.",
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "tls", "udp"}, false),
			},
			// sipRegistrar
			"sip_registrar": {
				Type:        schema.TypeString,
				Description: "The SIP server to be tested, specified by domain name or IP address.",
				Required:    true,
			},
			// user
			"user": {
				Type:        schema.TypeString,
				Description: "The username for SIP registration. This should be unique within a ThousandEyes account group.",
				Optional:    true,
			},
		},
	},
}
