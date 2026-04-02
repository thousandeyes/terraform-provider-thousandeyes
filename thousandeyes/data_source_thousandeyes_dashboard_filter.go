package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/dashboards"
)

func dataSourceThousandeyesDashboardFilter() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDashboardFilterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the dashboard filter.",
			},
			"filter_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the dashboard filter.",
			},
		},
		Description: "This data source allows you to search for Dashboard Filters by exact match. For more information, see [What is a Dashboard Filter](https://docs.thousandeyes.com/product-documentation/dashboards/dashboard-filters).",
	}
}

func dataSourceDashboardFilterRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*dashboards.DashboardsFiltersAPIService)(&apiClient.Common)

	name := d.Get("name").(string)
	log.Printf("[INFO] Reading ThousandEyes Dashboard Filters %s", name)
	req := api.GetDashboardsFilters()
	ctx := apiClient.GetConfig().Context
	req = SetAidFromContext(ctx, req)
	req = req.SearchPattern(name)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	var found *dashboards.ApiContextFilterResponse

	for _, filter := range resp.GetDashboardFilters() {
		if filter.Name == name {
			found = &filter
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any filter with the name: %s", name)
	}

	d.SetId(found.GetId())
	err = d.Set("name", found.GetName())
	if err != nil {
		return err
	}
	err = d.Set("filter_id", found.GetId())
	if err != nil {
		return err
	}

	return nil
}
