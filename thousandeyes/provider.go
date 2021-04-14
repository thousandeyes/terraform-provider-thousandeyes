package thousandeyes

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

// Provider for module
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_TOKEN", nil),
			},
			"account_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_AID", nil),
			},
			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TE_TIMEOUT", 0),
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
			"thousandeyes_voice_call":      resourceVoiceCall(),
			"thousandeyes_label":           resourceGroupLabel(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thousandeyes_account_group": dataSourceThousandeyesAccountGroup(),
			"thousandeyes_agent":         dataSourceThousandeyesAgent(),
			"thousandeyes_bgp_monitor":   dataSourceThousandeyesBGPMonitor(),
			"thousandeyes_integration":   dataSourceThousandeyesIntegration(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Thousand Eyes client")
	opts := thousandeyes.ClientOptions{
		AuthToken: d.Get("token").(string),
		AccountID: d.Get("account_group_id").(string),
		Timeout:   time.Second * time.Duration(d.Get("timeout").(int)),
	}
	return thousandeyes.NewClient(&opts), nil
}
