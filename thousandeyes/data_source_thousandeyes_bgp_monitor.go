package thousandeyes

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func dataSourceThousandeyesBGPMonitor() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesBGPMonitorsRead,

		Schema: map[string]*schema.Schema{

			"monitor_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "unique ID of BGP monitor",
				ConflictsWith: []string{"monitor_name"},
			},
			"monitor_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "display name of the BGP monitor",
				ConflictsWith: []string{"monitor_id"},
			},
			"ip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP address of the BGP monitor",
			},
			"network": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name of the autonomous system in which the monitor is found",
			},
			"monitor_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "either Public or Private, shows the type of monitor",
			},
		},
		Description: "This data source allows you to configure private and public BGP monitors. For more information, see [BGP Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests/bgp-tests).",
	}
}

func dataSourceThousandeyesBGPMonitorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	var found *thousandeyes.BGPMonitor

	searchName := d.Get("monitor_name").(string)
	searchMonitorID := int64(d.Get("monitor_id").(int))

	BGPMonitors, err := client.GetBPGMonitors()
	if err != nil {
		return err
	}

	for _, ar := range *BGPMonitors {
		if *ar.MonitorName == searchName {
			found = &ar
			break
		}
		if *ar.MonitorID == searchMonitorID {
			found = &ar
			break
		}

	}
	if found == nil {
		return fmt.Errorf("unable to locate any bgp by name: [%s] or ID: [%d]", searchName, searchMonitorID)
	}

	d.SetId(strconv.FormatInt(*found.MonitorID, 10))
	err = d.Set("monitor_name", *found.MonitorName)
	if err != nil {
		return err
	}
	err = d.Set("monitor_id", *found.MonitorID)
	if err != nil {
		return err
	}

	return nil
}
