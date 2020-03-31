package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceAgentServer() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the test",
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "interval to run test on",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "test description",
			},
			"server": {
				Type:        schema.TypeString,
				Description: "target host",
				Required:    true,
			},
			"bandwidth_measurements": {
				Type:        schema.TypeInt,
				Description: "set to 1 to measure bandwidth; defaults to 0. Only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set.",
				Optional:    true,
				Required:    false,
				Default:     1,
			},
			"bgp_measurements": {
				Type:         schema.TypeInt,
				Description:  "choose 1 to enable bgp measurements, 0 to disable; defaults to 1",
				Optional:     true,
				Required:     false,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"mtu_measurements": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to measure MTU sizes on network from agents to the target.",
				Optional:     true,
				Required:     false,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"num_path_traces": {
				Type:         schema.TypeInt,
				Description:  "number of path traces. default 3.",
				Default:      3,
				Optional:     true,
				Required:     false,
				ValidateFunc: validation.IntBetween(3, 10),
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "target port for agent",
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
				Required:     false,
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "protocol used by dependent Network tests (End-to-end, Path Trace, PMTUD); defaults to TCP",
				Optional:     true,
				Required:     false,
				Default:      "TCP",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "ICMP"}, false),
			},
			"probe_mode": {
				Type:         schema.TypeString,
				Description:  "probe mode used by End-to-end Network Test; only valid if protocol is set to TCP; defaults to AUTO",
				Optional:     true,
				Required:     false,
				Default:      "AUTO",
				ValidateFunc: validation.StringInSlice([]string{"AUTO", "SACK", "SYN"}, false),
			},
			// "bgp_monitors": {
			// 	Type:        schema.TypeList,
			// 	Description: "bgp monitors to use ",
			// 	Optional:    true,
			// 	Required:    false,
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"monitor_id": {
			// 				Type:        schema.TypeInt,
			// 				Description: "monitor id",
			// 				Optional:    true,
			// 			},
			// 		},
			// 	},
			// },
			"agents": {
				Type:        schema.TypeList,
				Description: "agents to use ",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": {
							Type:        schema.TypeInt,
							Description: "agent id",
							Optional:    true,
						},
					},
				},
			},
		},
		Create: resourceAgentServerCreate,
		Read:   resourceAgentServerRead,
		Update: resourceAgentServerUpdate,
		Delete: resourceAgentServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceAgentServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Agent %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	agent, err := client.GetAgentServer(id)
	if err != nil {
		return err
	}

	d.Set("name", agent.TestName)
	d.Set("id", agent.TestId)
	d.Set("interval", agent.Interval)
	d.Set("description", agent.Description)
	d.Set("bandwidth_measurements", agent.BandwidthMeasurements)
	d.Set("bgp_measurements", agent.BgpMeasurements)
	d.Set("mtu_measurements", agent.MtuMeasurements)
	d.Set("num_path_traces", agent.NumPathTraces)
	d.Set("port", agent.Port)
	d.Set("protocol", agent.Protocol)
	d.Set("server", agent.Server)
	// d.Set("bgp_monitors", agent.BgpMonitors)
	d.Set("agents", agent.Agents)

	return nil
}

func resourceAgentServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.AgentServer

	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}

	if d.HasChange("description") {
		update.Description = d.Get("description").(string)
	}

	if d.HasChange("bandwidth_measurements") {
		update.BandwidthMeasurements = d.Get("bandwidth_measurements").(int)
	}
	if d.HasChange("protocol") {
		update.Protocol = d.Get("protocol").(string)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("mtu_measurements") {
		update.MtuMeasurements = d.Get("mtu_measurements").(int)
	}
	if d.HasChange("num_path_traces") {
		update.NumPathTraces = d.Get("num_path_traces").(int)
	}

	if d.HasChange("port") {
		update.Port = d.Get("port").(int)
	}

	if d.HasChange("probe_mode") {
		update.ProbeMode = d.Get("probe_mode").(string)
	}

	if d.HasChange("bgp_measurements") {
		update.BgpMeasurements = d.Get("bgp_measurements").(int)
	}
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}

	// Need to figure this out

	// if d.HasChange("bgp_monitors") {
	// 	update.BgpMonitors = expandMonitors(d.Get("bgp_monitors").([]interface{}))
	// }

	if d.HasChange("server") {
		update.Server = d.Get("server").(string)
	}

	_, err := client.UpdateAgentServer(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceAgentServerRead(d, m)
}

func resourceAgentServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteAgentServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	webTrans := buildAgentServerStruct(d)
	httpTest, err := client.CreateAgentServer(*webTrans)
	if err != nil {
		return err
	}
	testID := httpTest.TestId
	d.SetId(strconv.Itoa(testID))
	return resourceAgentServerRead(d, m)
}

func buildAgentServerStruct(d *schema.ResourceData) *thousandeyes.AgentServer {
	transaction := thousandeyes.AgentServer{
		TestName:              d.Get("name").(string),
		Interval:              d.Get("interval").(int),
		Description:           d.Get("description").(string),
		BandwidthMeasurements: d.Get("bandwidth_measurements").(int),
		BgpMeasurements:       d.Get("bgp_measurements").(int),
		MtuMeasurements:       d.Get("mtu_measurements").(int),
		Port:                  d.Get("port").(int),
		Protocol:              d.Get("protocol").(string),
		ProbeMode:             d.Get("probe_mode").(string),
		Server:                d.Get("server").(string),
		Agents:                expandAgents(d.Get("agents").([]interface{})),
	}
	if attr, ok := d.GetOk("name"); ok {
		transaction.TestName = attr.(string)
	}
	if attr, ok := d.GetOk("interval"); ok {
		transaction.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("description"); ok {
		transaction.Description = attr.(string)
	}
	if attr, ok := d.GetOk("bandwidth_measurements"); ok {
		transaction.BandwidthMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("bgp_measurements"); ok {
		transaction.BgpMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("mtu_measurements"); ok {
		transaction.MtuMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("num_path_traces"); ok {
		transaction.NumPathTraces = attr.(int)
	}
	if attr, ok := d.GetOk("port"); ok {
		transaction.Port = attr.(int)
	}
	if attr, ok := d.GetOk("protocol"); ok {
		transaction.Protocol = attr.(string)
	}
	if attr, ok := d.GetOk("probe_mode"); ok {
		transaction.ProbeMode = attr.(string)
	}
	if attr, ok := d.GetOk("server"); ok {
		transaction.Server = attr.(string)
	}

	return &transaction
}
