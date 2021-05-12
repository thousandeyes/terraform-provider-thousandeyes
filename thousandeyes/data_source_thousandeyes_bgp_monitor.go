package thousandeyes

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/william20111/go-thousandeyes"
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
	}
}

func dataSourceThousandeyesBGPMonitorsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	var found *thousandeyes.BGPMonitor

	searchName := d.Get("monitor_name").(string)
	searchMonitorID := d.Get("monitor_id").(int)

	BGPMonitors, err := client.GetBPGMonitors()
	if err != nil {
		return err
	}

	for _, ar := range *BGPMonitors {
		if ar.MonitorName == searchName {
			found = &ar
			break
		}
		if ar.MonitorID == searchMonitorID {
			found = &ar
			break
		}

	}
	if found == nil {
		return fmt.Errorf("unable to locate any bgp by name: [%s] or ID: [%d]", searchName, searchMonitorID)
	}

	d.SetId(strconv.Itoa(found.MonitorID))
	d.Set("monitor_name", found.MonitorName)
	d.Set("monitor_id", found.MonitorID)

	return nil
}
