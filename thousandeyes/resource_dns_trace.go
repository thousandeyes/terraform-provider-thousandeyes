package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceDNSTrace() *schema.Resource {
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
			"dns_transport_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "transport protocol used for DNS requests; defaults to UDP",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "transport protocol used for DNS requests; defaults to UDP",
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "interval to run dns trace test on",
			},
		},
		Create: resourceDNSTraceCreate,
		Read:   resourceDNSTraceRead,
		Update: resourceDNSTraceUpdate,
		Delete: resourceDNSTraceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDNSTraceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetDNSTrace(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("dns_transport_protocol", test.DNSTransportProtocol)
	d.Set("domain", test.Domain)
	d.Set("interval", test.Interval)
	return nil
}

func resourceDNSTraceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.DNSTrace
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("dns_transport_protocol") {
		update.DNSTransportProtocol = d.Get("dns_transport_protocol").(string)
	}
	if d.HasChange("domain") {
		update.Domain = d.Get("domain").(string)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	_, err := client.UpdateDNSTrace(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceDNSTraceRead(d, m)
}

func resourceDNSTraceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteDNSTrace(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSTraceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	dnsTrace := buildDNSTraceStruct(d)
	dnsTest, err := client.CreateDNSTrace(*dnsTrace)
	if err != nil {
		return err
	}
	testID := dnsTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceDNSTraceRead(d, m)
}

func buildDNSTraceStruct(d *schema.ResourceData) *thousandeyes.DNSTrace {
	dnsTrace := thousandeyes.DNSTrace{
		TestName:             d.Get("name").(string),
		Agents:               expandAgents(d.Get("agents").([]interface{})),
		DNSTransportProtocol: d.Get("dns_transport_protocol").(string),
		Domain:               d.Get("domain").(string),
		Interval:             d.Get("interval").(int),
	}

	return &dnsTrace
}
