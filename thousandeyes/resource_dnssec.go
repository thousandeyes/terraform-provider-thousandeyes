package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceDNSSec() *schema.Resource {
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
			"domain": {
				Type:        schema.TypeString,
				Description: "target record for test, followed by record type (ie, www.thousandeyes.com A)",
				Required:    true,
			},
			"interval": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "interval to run test on, in seconds",
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
			},
		},
		Create: resourceDNSSecCreate,
		Read:   resourceDNSSecRead,
		Update: resourceDNSSecUpdate,
		Delete: resourceDNSSecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceDNSSecRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetDNSSec(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("domain", test.Domain)
	d.Set("interval", test.Interval)
	return nil
}

func resourceDNSSecUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.DNSSec
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("domain") {
		update.Domain = d.Get("domain").(string)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	_, err := client.UpdateDNSSec(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceDNSSecRead(d, m)
}

func resourceDNSSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteDNSSec(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSSecCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	resourceDNSSec := buildDNSSecStruct(d)
	resourceDNSSecTest, err := client.CreateDNSSec(*resourceDNSSec)
	if err != nil {
		return err
	}
	testID := resourceDNSSecTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceDNSSecRead(d, m)
}

func buildDNSSecStruct(d *schema.ResourceData) *thousandeyes.DNSSec {
	resourceDNSSec := thousandeyes.DNSSec{
		TestName: d.Get("name").(string),
		Agents:   expandAgents(d.Get("agents").([]interface{})),
		Domain:   d.Get("domain").(string),
		Interval: d.Get("interval").(int),
	}

	return &resourceDNSSec
}
