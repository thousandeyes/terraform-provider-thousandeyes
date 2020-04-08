package thousandeyes

import (
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceAgentToServer() *schema.Resource {

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the test",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Description: "test description",
			},

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
			"alert_rules": {
				Description: "get ruleId from /alert-rules endpoint. If alertsEnabled is set to 1 and alertRules is not included in a creation/update query, applicable defaults will be used.",
				Optional:    true,
				Required:    false,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type:        schema.TypeInt,
							Description: "If alertsEnabled is set to 1 and alertRules is not included in a creation/update query, applicable defaults will be used.",
							Optional:    true,
						},
					},
				},
			},
			"alerts_enabled": {
				Type:         schema.TypeInt,
				Description:  "choose 1 to enable alerts, or 0 to disable alerts. Defaults to 1",
				Optional:     true,
				Required:     false,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
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
			"bgp_monitors": {
				Type:        schema.TypeList,
				Description: "bgp monitors to use ",
				Optional:    true,
				Required:    false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type:        schema.TypeInt,
							Description: "monitor id",
							Optional:    true,
						},
					},
				},
			},
			"enabled": {
				Type:         schema.TypeInt,
				Description:  "choose 1 to enable the test, 0 to disable the test",
				Optional:     true,
				Required:     false,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "interval to run test on",
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
			"path_trace_mode": {
				Type:         schema.TypeString,
				Description:  "choose inSession to perform the path trace within a TCP session; defaults to classic",
				Optional:     true,
				Required:     false,
				Default:      "classic",
				ValidateFunc: validation.StringInSlice([]string{"classic", "inSession"}, false),
			},
			"port": {
				Type:         schema.TypeInt,
				Description:  "target port for agent",
				Default:      80,
				ValidateFunc: validation.IntBetween(1, 65535),
				Optional:     true,
				Required:     false,
			},
			"probe_mode": {
				Type:         schema.TypeString,
				Description:  "probe mode used by End-to-end Network Test; only valid if protocol is set to TCP; defaults to AUTO",
				Optional:     true,
				Required:     false,
				Default:      "AUTO",
				ValidateFunc: validation.StringInSlice([]string{"AUTO", "SACK", "SYN"}, false),
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "protocol used by dependent Network tests (End-to-end, Path Trace, PMTUD); defaults to TCP",
				Optional:     true,
				Required:     false,
				Default:      "TCP",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "ICMP"}, false),
			},
			"server": {
				Type:        schema.TypeString,
				Description: "target host",
				Required:    true,
			},
		},
		Create: resourceAgentToServerCreate,
		Read:   resourceAgentToServerRead,
		Update: resourceAgentToServerUpdate,
		Delete: resourceAgentToServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceAgentToServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] ## Reading Thousandeyes Agent to Server Test: %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	agent, err := client.GetAgentServer(id)
	if err != nil {
		return err
	}
	// The API returns port and server together int he response.
	// We need to split it up in orcer to avoid TF seeing it as a change
	serverParsed := strings.Split(agent.Server, ":")
	d.Set("port", serverParsed[1])
	d.Set("server", serverParsed[0])

	d.Set("agents", agent.Agents)
	d.Set("alerts_enabled", agent.AlertsEnabled)
	d.Set("alert_rules", agent.AlertRules)
	d.Set("description", agent.Description)
	d.Set("enabled", agent.Enabled)
	d.Set("name", agent.TestName)
	d.Set("bandwidth_measurements", agent.BandwidthMeasurements)
	d.Set("bgp_beasurements", agent.BgpMeasurements)
	d.Set("bgp_monitors", agent.BgpMonitors)
	d.Set("interval", agent.Interval)
	d.Set("mtu_measurements", agent.MtuMeasurements)
	d.Set("num_path_traces", agent.NumPathTraces)
	d.Set("probe_mode", agent.ProbeMode)
	d.Set("protocol", agent.Protocol)

	return nil
}

func resourceAgentToServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] ###### Updating ThousandEyes Agent to Server Test %s", d.Id())

	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.AgentServer

	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}

	if d.HasChange("alerts_enabled") {
		update.AlertsEnabled = d.Get("alerts_enabled").(int)
	}
	if d.HasChange("alert_rules") {
		update.AlertRules = expandAlertRules(d.Get("alert_rules").([]interface{}))
	}

	if d.HasChange("description") {
		update.Description = d.Get("description").(string)
	}
	if d.HasChange("enabled") {
		update.Enabled = d.Get("enabled").(int)
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("bandwidth_measurements") {
		update.BandwidthMeasurements = d.Get("bandwidth_measurements").(int)
	}
	if d.HasChange("bgp_measurements") {
		update.BgpMeasurements = d.Get("bgp_measurements").(int)
	}
	if d.HasChange("bgp_monitors") {
		update.BgpMonitors = expandBGPMonitors(d.Get("bgp_monitors").([]interface{}))
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
	if d.HasChange("protocol") {
		update.Protocol = d.Get("probe_mode").(string)
	}
	if d.HasChange("server") {
		update.Server = d.Get("server").(string)
	}

	_, err := client.UpdateAgentServer(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceAgentToServerRead(d, m)
}

func resourceAgentToServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] ###### Deleting ThousandEyes Agent to Server Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteAgentServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentToServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] ###### Creating ThousandEyes Agent to Server Test %s", d.Id())

	agentServer := buildAgentToServerStruct(d)
	agentToServer, err := client.CreateAgentServer(*agentServer)
	if err != nil {
		return err
	}
	testID := agentToServer.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceAgentToServerRead(d, m)
}

func buildAgentToServerStruct(d *schema.ResourceData) *thousandeyes.AgentServer {
	transaction := thousandeyes.AgentServer{
		Agents:                expandAgents(d.Get("agents").([]interface{})),
		AlertsEnabled:         d.Get("alerts_enabled").(int),
		Description:           d.Get("description").(string),
		Enabled:               d.Get("enabled").(int),
		TestName:              d.Get("name").(string),
		BandwidthMeasurements: d.Get("bandwidth_measurements").(int),
		BgpMeasurements:       d.Get("bgp_measurements").(int),
		BgpMonitors:           expandBGPMonitors(d.Get("bgp_monitors").([]interface{})),
		Interval:              d.Get("interval").(int),
		MtuMeasurements:       d.Get("mtu_measurements").(int),
		NumPathTraces:         d.Get("num_path_traces").(int),
		Port:                  d.Get("port").(int),
		ProbeMode:             d.Get("probe_mode").(string),
		Protocol:              d.Get("protocol").(string),
		Server:                d.Get("server").(string),
	}
	if attr, ok := d.GetOk("alerts_enabled"); ok {
		transaction.AlertsEnabled = attr.(int)
	}
	if attr, ok := d.GetOk("alert_rules"); ok {
		transaction.AlertRules = expandAlertRules(attr.([]interface{}))
	}
	// if attr, ok := d.GetOk("alert_rules"); ok {
	// 	transaction.AlertRules = attr.(int)
	// }
	if attr, ok := d.GetOk("description"); ok {
		transaction.Description = attr.(string)
	}
	if attr, ok := d.GetOk("enabled"); ok {
		transaction.Enabled = attr.(int)
	}
	if attr, ok := d.GetOk("name"); ok {
		transaction.TestName = attr.(string)
	}
	if attr, ok := d.GetOk("bandwidth_measurements"); ok {
		transaction.BandwidthMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("bgp_measurements"); ok {
		transaction.BgpMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("bgp_monitors"); ok {
		transaction.BgpMonitors = expandBGPMonitors(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("interval"); ok {
		transaction.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("mtu_measurements"); ok {
		transaction.MtuMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("num_path_traces"); ok {
		transaction.NumPathTraces = attr.(int)
	}

	if attr, ok := d.GetOk("probe_mode"); ok {
		transaction.ProbeMode = attr.(string)
	}
	if attr, ok := d.GetOk("protocol"); ok {
		transaction.Protocol = attr.(string)
	}
	if attr, ok := d.GetOk("port"); ok {
		transaction.Port = attr.(int)
	}

	if attr, ok := d.GetOk("server"); ok {
		transaction.Server = attr.(string)
	}

	return &transaction
}
