package thousandeyes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

type accountGroupIdKeyType string

const accountGroupIdKey accountGroupIdKeyType = "aid"

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
				DefaultFunc: schema.EnvDefaultFunc("TE_API_ENDPOINT", "https://api.thousandeyes.com/v7"),
				Description: "The ThousandEyes API Endpoint's URL. E.g. https://api.thousandeyes.com/v7",
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
			"thousandeyes_api_test":        resourceAPI(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thousandeyes_account_group": dataSourceThousandeyesAccountGroup(),
			"thousandeyes_agent":         dataSourceThousandeyesAgent(),
			"thousandeyes_bgp_monitor":   dataSourceThousandeyesBGPMonitor(),
			"thousandeyes_alert_rule":    dataSourceThousandeyesAlertRule(),
		},
		ConfigureContextFunc: providerConfigureWithContext,
	}
}

func providerConfigureWithContext(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Println("[INFO] Initializing ThousandEyes client")

	// set AID to context
	ctx := context.WithValue(context.Background(), accountGroupIdKey, d.Get("account_group_id").(string))

	configuration := &client.Configuration{
		AuthToken:  d.Get("token").(string),
		UserAgent:  "ThousandEyes Terraform Provider",
		ServerURL:  d.Get("api_endpoint").(string),
		HTTPClient: &http.Client{Timeout: time.Second * time.Duration(d.Get("timeout").(int))},
		Context:    ctx,
	}

	apiClient := client.NewAPIClient(configuration)
	return apiClient, nil
}
