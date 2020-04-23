package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

// voiceAuthElem is the used for schema for SIP server tests.
// Similar to the schema used in sip-server, but some
// option requirements are different.
var voiceCallAuthElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"auth_user": {
			Type:        schema.TypeString,
			Description: "username for authentication with SIP server",
			Required:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "password for authentication with SIP server",
			Required:    true,
		},
		"port": {
			Type:         schema.TypeInt,
			Description:  "port number for the chosen protocol",
			Required:     true,
			ValidateFunc: validation.IntBetween(1024, 65535),
		},
		"protocol": {
			Type:         schema.TypeString,
			Description:  "transport layer for SIP communication: TCP, TLS (TLS over TCP), or UDP. Defaults to TCP",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"TCP", "TLS", "UDP"}, false),
		},
		"user": {
			Type:        schema.TypeString,
			Description: "username for SIP registration; should be unique within a ThousandEyes Account Group",
			Required:    true,
		},
	},
}

func resourceVoiceCall() *schema.Resource {
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
				ValidateFunc: validation.IntBetween(0, 56),
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
			"source_sip_credentials": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "source SIP server parameters",
				Elem:        voiceCallAuthElem,
			},
			"target_agent_id": {
				Type:        schema.TypeInt,
				Description: "target agent to run test on",
				Required:    true,
			},
			"target_sip_credentials": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "target SIP server parameters",
				Elem:        voiceCallAuthElem,
			},
		},
		Create: resourceVoiceCallCreate,
		Read:   resourceVoiceCallRead,
		Update: resourceVoiceCallUpdate,
		Delete: resourceVoiceCallDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceVoiceCallRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetVoiceCall(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("codec_id", test.CodecID)
	d.Set("dscp_id", test.DscpID)
	d.Set("duration", test.Duration)
	d.Set("interval", test.Interval)
	d.Set("jitter_buffer", test.JitterBuffer)
	d.Set("source_sip_credentials", test.SourceSipCredentials)
	d.Set("target_agent_id", test.TargetAgentID)
	d.Set("target_sip_credentials", test.TargetSipCredentials)
	return nil
}

func resourceVoiceCallUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.VoiceCall
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
	if d.HasChange("source_sip_credentials") {
		update.SourceSipCredentials = unpackSIPAuthData(d.Get("source_sip_credentials"))
	}
	if d.HasChange("target_agent_id") {
		update.TargetAgentID = d.Get("target_agent_id").(int)
	}
	if d.HasChange("target_sip_credentials") {
		update.TargetSipCredentials = unpackSIPAuthData(d.Get("target_sip_credentials"))
	}
	_, err := client.UpdateVoiceCall(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceVoiceCallRead(d, m)
}

func resourceVoiceCallDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteVoiceCall(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceVoiceCallCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	resourceVoiceCall := buildVoiceCallStruct(d)
	resourceVoiceCallTest, err := client.CreateVoiceCall(*resourceVoiceCall)
	if err != nil {
		return err
	}
	testID := resourceVoiceCallTest.TestID
	d.SetId(strconv.Itoa(testID))
	log.Printf("[DEBUG] Creating VoiceCall:\n%+v", resourceVoiceCall)
	return resourceVoiceCallRead(d, m)
}

func buildVoiceCallStruct(d *schema.ResourceData) *thousandeyes.VoiceCall {
	resourceVoiceCall := thousandeyes.VoiceCall{
		TestName:             d.Get("name").(string),
		Agents:               expandAgents(d.Get("agents").([]interface{})),
		CodecID:              d.Get("codec_id").(int),
		DscpID:               d.Get("dscp_id").(int),
		Duration:             d.Get("duration").(int),
		Interval:             d.Get("interval").(int),
		JitterBuffer:         d.Get("jitter_buffer").(int),
		SourceSipCredentials: unpackSIPAuthData(d.Get("source_sip_credentials")),
		TargetAgentID:        d.Get("target_agent_id").(int),
		TargetSipCredentials: unpackSIPAuthData(d.Get("target_sip_credentials")),
	}

	return &resourceVoiceCall
}
