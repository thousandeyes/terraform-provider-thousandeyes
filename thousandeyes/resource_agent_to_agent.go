package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceAgentToAgent() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"description": {
				Type:        schema.TypeString,
				Required:    false,
				Optional:    true,
				Default:     "",
				Description: "defaults to empty string",
			},

			"direction": {
				Type: schema.TypeString,
				Description: "[TO_TARGET, FROM_TARGET, BIDIRECTIONAL]	Direction of the test (affects how results are shown)",
				Optional:     false,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"TO_TARGET", "FROM_TARGET", "BIDIRECTIONAL"}, false),
			},
			"dscp_id": {
				Type:         schema.TypeInt,
				Description:  "A Differentiated Services Code Point (DSCP) is a value found in an IP packet header which is used to request a level of priority for delivery (Defined in RFC 2474 https://www.ietf.org/rfc/rfc2474.txt). It is one of the Quality of Service management tools used in router configuration to protect real-time and high priority data applications.",
				Required:     false,
				Default:      0,
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 8, 16, 24, 32, 40, 48, 56, 10, 12, 14, 18, 20, 22, 26, 28, 30, 34, 36, 38, 46, 44}),
			},
			"enabled": {
				Type: schema.TypeInt,
				Description: "0 or 1	choose 1 to enable the test, 0 to disable the test",
				Optional:     true,
				Required:     false,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"interval": {
				Type:     schema.TypeInt,
				Required: true,
				Description: "[120, 300, 600, 900, 1800, 3600]	value in seconds",
				ValidateFunc: validation.IntInSlice([]int{120, 300, 600, 900, 1800, 3600}),
			},
			"mss": {
				Type: schema.TypeInt,
				Description: "(30..1400)	Maximum Segment Size, in bytes.",
				ValidateFunc: validation.IntBetween(30, 1400),
				Optional:     true,
				Required:     false,
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
			"target_agent_id": {
				Type:     schema.TypeInt,
				Optional: false,
				Required: true,
				Description: "pull from /agents endpoint	Both the 'agents': [] and the targetAgentID cannot be cloud agents. Can be Enterprise Agent -> Cloud, Cloud -> Enterprise Agent, or Enterprise Agent -> Enterprise Agent",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Optional:    false,
				Description: "Test name must be unique",
			},
			"throughput_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Required:     false,
				Default:      10000,
				Description:  "Defaults to 10000",
				ValidateFunc: validation.IntBetween(5000, 10000),
			},
			"throughput_measurements": {
				Type:         schema.TypeInt,
				ValidateFunc: validation.IntBetween(0, 1),
				Optional:     true,
				Required:     false,
				Default:      0,
				Description: "0 or 1	defaults to 0 (disabled), not allowed when source (or target) of the test is a cloud agent",
			},
			"throughput_rate": {
				Type:         schema.TypeInt,
				Description:  "for UDP only",
				Optional:     true,
				Required:     false,
				Default:      0,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
		},
		Create: resourceAgentToAgentCreate,
		Read:   resourceAgentToAgentRead,
		Update: resourceAgentToAgentUpdate,
		Delete: resourceAgentToAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceAgentToAgentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Agent to Agent - testID: [%s]", d.Id())
	id, _ := strconv.Atoi(d.Id())
	agent, err := client.GetAgentAgent(id)
	if err != nil {
		return err
	}

	d.Set("agents", agent.Agents)
	d.Set("alert_rules", agent.AlertRules)
	d.Set("alertsEnabled", agent.AlertsEnabled)
	d.Set("bgp_measurements", agent.BgpMeasurements)
	d.Set("bgp_monitors", agent.BgpMonitors)
	d.Set("description", agent.Description)
	d.Set("direction", agent.Direction)
	d.Set("dscp", agent.Dscp)
	d.Set("dscp_id", agent.DscpID)
	d.Set("enabled", agent.Enabled)
	d.Set("interval", agent.Interval)
	d.Set("id", agent.TestID)
	d.Set("name", agent.TestName)
	d.Set("mss", agent.Mss)
	d.Set("num_path_traces", agent.NumPathTraces)
	d.Set("port", agent.Port)
	d.Set("protocol", agent.Protocol)
	d.Set("target_agent_id", agent.TargetAgentID)
	d.Set("throughput_duration", agent.ThroughputDuration)
	d.Set("throughput_measurements", agent.ThroughputMeasurements)
	d.Set("throughput_rate", agent.ThroughputRate)

	return nil
}

func resourceAgentToAgentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Agent to Agent Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.AgentAgent

	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("alert_rules") {
		update.AlertRules = expandAlertRules(d.Get("alert_rules").([]interface{}))
	}
	if d.HasChange("alertsEnabled") {
		update.AlertsEnabled = d.Get("alerts_enabled").(int)
	}
	if d.HasChange("bgp_measurements") {
		update.BgpMeasurements = d.Get("bgp_measurements").(int)
	}
	if d.HasChange("bgp_monitors") {
		update.BgpMonitors = expandBGPMonitors(d.Get("bgp_monitors").([]interface{}))
	}
	if d.HasChange("description") {
		update.Description = d.Get("description").(string)
	}
	if d.HasChange("dscp_id") {
		update.DscpID = d.Get("dscp_id").(int)
	}
	if d.HasChange("enabled") {
		update.Enabled = d.Get("enabled").(int)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("mss") {
		update.Mss = d.Get("mss").(int)
	}
	if d.HasChange("num_path_traces") {
		update.NumPathTraces = d.Get("num_path_traces").(int)
	}
	if d.HasChange("port") {
		update.Port = d.Get("port").(int)
	}
	if d.HasChange("protocol") {
		update.Protocol = d.Get("protocol").(string)
	}
	if d.HasChange("throughput_duration") {
		update.ThroughputDuration = d.Get("throughput_duration").(int)
	}
	if d.HasChange("throughput_measurements") {
		update.ThroughputMeasurements = d.Get("throughput_measurements").(int)
	}
	if d.HasChange("throughput_rate") {
		update.ThroughputRate = d.Get("throughput_rate").(int)
	}

	_, err := client.UpdateAgentAgent(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceAgentToAgentRead(d, m)
}

func resourceAgentToAgentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Agent to Agent Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteAgentAgent(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentToAgentCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Agent to Agent Test %s", d.Id())
	webTrans := buildAgentToAgentStruct(d)
	httpTest, err := client.CreateAgentAgent(*webTrans)
	if err != nil {
		return err
	}
	testID := httpTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceAgentToAgentRead(d, m)
}

func buildAgentToAgentStruct(d *schema.ResourceData) *thousandeyes.AgentAgent {
	transaction := thousandeyes.AgentAgent{
		Agents:                 expandAgents(d.Get("agents").([]interface{})),
		AlertsEnabled:          d.Get("alerts_enabled").(int),
		BgpMeasurements:        d.Get("bgp_measurements").(int),
		Description:            d.Get("description").(string),
		Direction:              d.Get("direction").(string),
		DscpID:                 d.Get("dscp_id").(int),
		Enabled:                d.Get("enabled").(int),
		Interval:               d.Get("interval").(int),
		Mss:                    d.Get("mss").(int),
		NumPathTraces:          d.Get("num_path_traces").(int),
		Port:                   d.Get("port").(int),
		Protocol:               d.Get("protocol").(string),
		TestName:               d.Get("name").(string),
		TargetAgentID:          d.Get("target_agent_id").(int),
		ThroughputDuration:     d.Get("throughput_duration").(int),
		ThroughputMeasurements: d.Get("throughput_measurements").(int),
		ThroughputRate:         d.Get("throughput_rate").(int),
	}

	if attr, ok := d.GetOk("alerts_enabled"); ok {
		transaction.AlertsEnabled = attr.(int)
	}
	if attr, ok := d.GetOk("alert_rules"); ok {
		transaction.AlertRules = expandAlertRules(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("description"); ok {
		transaction.Description = attr.(string)
	}
	if attr, ok := d.GetOk("enabled"); ok {
		transaction.Enabled = attr.(int)
	}
	if attr, ok := d.GetOk("name"); ok {
		transaction.TestName = attr.(string)
	}
	if attr, ok := d.GetOk("bgp_measurements"); ok {
		transaction.BgpMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("bgp_monitors"); ok {
		transaction.BgpMonitors = expandBGPMonitors(attr.([]interface{}))
	}
	if attr, ok := d.GetOk("direction"); ok {
		transaction.Direction = attr.(string)
	}
	if attr, ok := d.GetOk("dscp_id"); ok {
		transaction.DscpID = attr.(int)
	}
	if attr, ok := d.GetOk("interval"); ok {
		transaction.Interval = attr.(int)
	}
	if attr, ok := d.GetOk("mss"); ok {
		transaction.Mss = attr.(int)
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
	if attr, ok := d.GetOk("target_agent_id"); ok {
		transaction.TargetAgentID = attr.(int)
	}
	if attr, ok := d.GetOk("throughput_duration"); ok {
		transaction.ThroughputDuration = attr.(int)
	}
	if attr, ok := d.GetOk("throughput_measurements"); ok {
		transaction.ThroughputMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("throughput_rate"); ok {
		transaction.ThroughputRate = attr.(int)
	}

	return &transaction
}
