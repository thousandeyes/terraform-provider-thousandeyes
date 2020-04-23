package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceRTPStream() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the test",
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
			"codec_id": {
				Type:         schema.TypeInt,
				Description:  "codec to use",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 8),
			},
			"dscp_id": {
				Type:         schema.TypeInt,
				Description:  "DSCP to use",
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 44, 46, 48, 56}),
			},
			"duration": {
				Type:         schema.TypeInt,
				Description:  "duration of test, in seconds (5 to 30)",
				Optional:     true,
				ValidateFunc: validation.IntBetween(5, 30),
			},
			"interval": {
				Type:         schema.TypeInt,
				Description:  "interval to run test on, in seconds",
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
			},
			"jitter_buffer": {
				Type:         schema.TypeInt,
				Description:  "de-jitter buffer size, in seconds (0 to 150)",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 150),
			},
			"target_agent_id": {
				Type:        schema.TypeInt,
				Description: "target agent to run test on",
				Required:    true,
			},
		},
		Create: resourceRTPStreamCreate,
		Read:   resourceRTPStreamRead,
		Update: resourceRTPStreamUpdate,
		Delete: resourceRTPStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceRTPStreamRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetRTPStream(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("interval", test.Interval)
	return nil
}

func resourceRTPStreamUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.RTPStream
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("codec_id") {
		update.CodecID = d.Get("codec_id").(int)
	}
	if d.HasChange("dscp_id") {
		update.DscpID = d.Get("dscp_id").(int)
	}
	if d.HasChange("duration") {
		update.Duration = d.Get("duration").(int)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("jitter_buffer") {
		update.JitterBuffer = d.Get("jitter_buffer").(int)
	}
	if d.HasChange("target_agent_id") {
		update.TargetAgentID = d.Get("target_agent_id").(int)
	}
	_, err := client.UpdateRTPStream(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceRTPStreamRead(d, m)
}

func resourceRTPStreamDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteRTPStream(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceRTPStreamCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	resourceRTPStream := buildRTPStreamStruct(d)
	resourceRTPStreamTest, err := client.CreateRTPStream(*resourceRTPStream)
	if err != nil {
		return err
	}
	testID := resourceRTPStreamTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceRTPStreamRead(d, m)
}

func buildRTPStreamStruct(d *schema.ResourceData) *thousandeyes.RTPStream {
	resourceRTPStream := thousandeyes.RTPStream{
		TestName:      d.Get("name").(string),
		Agents:        expandAgents(d.Get("agents").([]interface{})),
		CodecID:       d.Get("codec_id").(int),
		DscpID:        d.Get("dscp_id").(int),
		Duration:      d.Get("duration").(int),
		Interval:      d.Get("interval").(int),
		JitterBuffer:  d.Get("jitter_buffer").(int),
		TargetAgentID: d.Get("target_agent_id").(int),
	}

	return &resourceRTPStream
}
