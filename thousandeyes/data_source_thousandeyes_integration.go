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
				Description: "(PagerDuty only) always set to `Auth Token`",
			},
			"auth_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(PagerDuty only) PagerDuty user`",
			},
			"auth_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(PagerDuty only) authentication token",
			},
			"channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "(Slack only) Slack #channel or @user",
			},
			"integration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "unique ID of the integration",
			},
			"integration_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the integration",
			},
			"integration_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of integration",
			},
			"target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "target URL of the integration",
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
	err = d.Set("integration_name", found.IntegrationName)
	if err != nil {
		return err
	}
	err = d.Set("integration_id", found.IntegrationID)
	if err != nil {
		return err
	}

	return nil
}
