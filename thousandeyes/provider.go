package thousandeyes

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

// Global variable for account group ID, as we must be aware of it in
// functions that will not have access to it otherwise.
var accountGroupId int64

func New() func() *schema.Provider {
	return func() *schema.Provider {
		return Provider()
	}
}

// Provider for module
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_TOKEN", nil),
				Description: "The ThousandEyes organization's authentication token.",
			},
			"account_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_AID", nil),
				Description: "The ThousandEyes account group's unique ID.",
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_TIMEOUT", 0),
				Description: "The timeout value.",
			},
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_API_ENDPOINT", "https://api.thousandeyes.com/v6"),
				Description: "The ThousandEyes API Endpoint's URL. E.g. https://api.thousandeyes.com/v6",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"thousandeyes_alert_rule":      resourceAlertRule(),
			"thousandeyes_http_server":     resourceHTTPServer(),
			"thousandeyes_page_load":       resourcePageLoad(),
			"thousandeyes_agent_to_server": resourceAgentToServer(),
			"thousandeyes_web_transaction": resourceWebTransaction(),
			"thousandeyes_agent_to_agent":  resourceAgentToAgent(),
			"thousandeyes_dns_server":      resourceDNSServer(),
			"thousandeyes_bgp":             resourceBGP(),
			"thousandeyes_dnssec":          resourceDNSSec(),
			"thousandeyes_dns_trace":       resourceDNSTrace(),
			"thousandeyes_ftp_server":      resourceFTPServer(),
			"thousandeyes_sip_server":      resourceSIPServer(),
			"thousandeyes_voice":           resourceRTPStream(),
			"thousandeyes_label":           resourceGroupLabel(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thousandeyes_account_group": dataSourceThousandeyesAccountGroup(),
			"thousandeyes_agent":         dataSourceThousandeyesAgent(),
			"thousandeyes_bgp_monitor":   dataSourceThousandeyesBGPMonitor(),
			"thousandeyes_integration":   dataSourceThousandeyesIntegration(),
		},
		ConfigureContextFunc: providerConfigureWithContext,
	}
}

func providerConfigureWithContext(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Initializing ThousandEyes client")
	opts := thousandeyes.ClientOptions{
		AuthToken:   d.Get("token").(string),
		AccountID:   d.Get("account_group_id").(string),
		Timeout:     time.Second * time.Duration(d.Get("timeout").(int)),
		UserAgent:   "ThousandEyes Terraform Provider",
		APIEndpoint: d.Get("api_endpoint").(string),
	}
	var err error
	if opts.AccountID != "" {
		accountGroupId, err = strconv.ParseInt(opts.AccountID, 10, 64)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return thousandeyes.NewClient(&opts), nil
}
