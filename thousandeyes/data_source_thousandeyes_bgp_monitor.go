package thousandeyes

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/bgpmonitors"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func dataSourceThousandeyesBGPMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesBGPMonitorsRead,

		Schema: map[string]*schema.Schema{

			"monitor_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The unique ID of BGP monitor.",
				ConflictsWith: []string{"monitor_name"},
			},
			"monitor_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "The display name of the BGP monitor.",
				ConflictsWith: []string{"monitor_id"},
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP address of the BGP monitor.",
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the autonomous system in which the monitor is found.",
			},
			"monitor_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of BGP monitor (either Public or Private).",
			},
		},
		Description: "This data source allows you to configure private and public BGP monitors. For more information, see [BGP Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests/bgp-tests).",
	}
}

func dataSourceThousandeyesBGPMonitorsRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*bgpmonitors.BGPMonitorsAPIService)(&apiClient.Common)

	var found *bgpmonitors.Monitor

	searchName := d.Get("monitor_name").(string)
	searchMonitorID := d.Get("monitor_id").(string)

	req := api.GetBgpMonitors()
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	for _, ar := range resp.GetMonitors() {
		if *ar.MonitorName == searchName {
			found = &ar
			break
		}
		if *ar.MonitorId == searchMonitorID {
			found = &ar
			break
		}

	}
	if found == nil {
		return fmt.Errorf("unable to locate any bgp by name: [%s] or ID: [%s]", searchName, searchMonitorID)
	}

	d.SetId(*found.MonitorId)
	err = d.Set("monitor_name", *found.MonitorName)
	if err != nil {
		return err
	}
	err = d.Set("monitor_id", *found.MonitorId)
	if err != nil {
		return err
	}

	return nil
}
