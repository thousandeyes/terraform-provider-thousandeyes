package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func dataSourceThousandeyesIntegration() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesIntegrationRead,

		Schema: map[string]*schema.Schema{
			"auth_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(PagerDuty only) The authentication method. Always set to `Auth Token`.",
			},
			"auth_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(PagerDuty only) The PagerDuty user.",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(PagerDuty only) The PagerDuty authentication token.",
			},
			"channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(Slack only) The Slack #channel or @user.",
			},
			"integration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the integration.",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the integration.",
			},
			"integration_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of integration.",
			},
			"target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The target URL of the integration.",
			},
		},
		Description: "This data source allows you to configure a third-party integration for ThousandEyes alerts. Supported integrations include PagerDuty, Slack, AppDynamics, and ServiceNow. For more information, see [Alerts Integrations](https://docs.thousandeyes.com/product-documentation/alerts/integrations).",
	}
}

func dataSourceThousandeyesIntegrationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	var found thousandeyes.Integration

	searchName := d.Get("integration_name").(string)

	integrations, err := client.GetIntegrations()
	if err != nil {
		return err
	}

	if searchName != "" {
		log.Printf("[INFO] ###### Reading Thousandeyes integration by name [%s]", searchName)

		for _, ar := range *integrations {
			if *ar.IntegrationName == searchName {
				found = ar
				break
			}
		}
	} else {
		return fmt.Errorf("must define integration name")
	}

	if found == (thousandeyes.Integration{}) {
		return fmt.Errorf("unable to locate any integration by name: %s", searchName)
	}

	d.SetId(*found.IntegrationID)
	err = ResourceRead(d, &found)
	if err != nil {
		return err
	}

	return nil
}
