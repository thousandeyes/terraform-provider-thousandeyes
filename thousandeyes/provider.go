package thousandeyes

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
		},
		ResourcesMap: map[string]*schema.Resource{
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
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thousandeyes_agent":       dataSourceThousandeyesAgent(),
			"thousandeyes_alert_rule":  dataSourceThousandeyesAlertRule(),
			"thousandeyes_bgp_monitor": dataSourceThousandeyesBGPMonitor(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Thousand Eyes client")
	return thousandeyes.NewClient(d.Get("token").(string), d.Get("account_group_id").(string)), nil
}
