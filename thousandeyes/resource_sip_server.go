package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

// resourceSIPServerAuthElem is the used for schema for SIP server tests.
// Similar to the schema used in voice-call, but some
// option requirements are different.
var sipAuthElem = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"auth_user": {
			Type:        schema.TypeString,
			Description: "username for authentication with SIP server",
			Optional:    true,
		},
		"password": {
			Type:        schema.TypeString,
			Description: "password for authentication with SIP server",
			Optional:    true,
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
	},
}

func resourceSIPServer() *schema.Resource {
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
			"interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "interval to run test on, in seconds",
				ValidateFunc: validation.IntInSlice([]int{120, 300, 600, 900, 1800, 3600}),
			},
			"target_sip_credentials": {
				Type:        schema.TypeMap,
				Required:    true,
				Description: "target SIP server parameters",
				Elem:        sipAuthElem,
			},
		},
		Create: resourceSIPServerCreate,
		Read:   resourceSIPServerRead,
		Update: resourceSIPServerUpdate,
		Delete: resourceSIPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceSIPServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetSIPServer(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("interval", test.Interval)
	d.Set("target_sip_credentials", test.TargetSipCredentials)
	return nil
}

func resourceSIPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.SIPServer
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("target_sip_credentials") {
		update.TargetSipCredentials = unpackSIPAuthData(d.Get("target_sip_credentials"))
	}
	_, err := client.UpdateSIPServer(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceSIPServerRead(d, m)
}

func resourceSIPServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteSIPServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceSIPServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	resourceSIPServer := buildSIPServerStruct(d)
	resourceSIPServerTest, err := client.CreateSIPServer(*resourceSIPServer)
	if err != nil {
		return err
	}
	testID := resourceSIPServerTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceSIPServerRead(d, m)
}

func buildSIPServerStruct(d *schema.ResourceData) *thousandeyes.SIPServer {
	resourceSIPServer := thousandeyes.SIPServer{
		TestName:             d.Get("name").(string),
		Agents:               expandAgents(d.Get("agents").([]interface{})),
		Interval:             d.Get("interval").(int),
		TargetSipCredentials: unpackSIPAuthData(d.Get("target_sip_credentials")),
	}

	return &resourceSIPServer
}
