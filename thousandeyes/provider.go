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
		},
		ResourcesMap: map[string]*schema.Resource{
			"thousandeyes_http_server":     resourceHttpServer(),
			"thousandeyes_page_load":       resourcePageLoad(),
			"thousandeyes_agent_to_server": resourceAgentServer(),
			"thousandeyes_web_transaction": resourceWebTransaction(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"thousandeyes_agent": dataSourceThousandeyesAgent(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing PagerDuty client")
	return thousandeyes.NewClient(d.Get("token").(string)), nil
}
