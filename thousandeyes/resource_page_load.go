package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourcePageLoad() *schema.Resource {
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
				Description: "interval to run page load test on",
			},
			"http_interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "interval to run http server test on",
			},
			"auth_type": {
				Type:         schema.TypeString,
				Description:  "auth type",
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "BASIC", "NTLM", "KERBEROS"}, false),
			},
			"http_version": {
				Type:         schema.TypeInt,
				Description:  "2 for default (prefer HTTP/2), 1 for HTTP/1.1 only",
				Default:      2,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"url": {
				Type:        schema.TypeString,
				Description: "target for the test",
				Required:    true,
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
		},
		Create: resourcePageLoadCreate,
		Read:   resourcePageLoadRead,
		Update: resourcePageLoadUpdate,
		Delete: resourcePageLoadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourcePageLoadRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetPageLoad(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("auth_type", test.AuthType)
	d.Set("interval", test.Interval)
	d.Set("http_version", test.HTTPVersion)
	d.Set("url", test.URL)
	d.Set("agents", test.Agents)
	d.Set("http_interval", test.HTTPInterval)
	return nil
}

func resourcePageLoadUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.PageLoad
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("auth_type") {
		update.AuthType = d.Get("auth_type").(string)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("http_version") {
		update.HTTPVersion = d.Get("http_version").(int)
	}
	if d.HasChange("url") {
		update.URL = d.Get("url").(string)
	}
	if d.HasChange("http_interval") {
		update.HTTPInterval = d.Get("http_interval").(int)
	}
	_, err := client.UpdatePageLoad(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourcePageLoadRead(d, m)
}

func resourcePageLoadDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeletePageLoad(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourcePageLoadCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	pageLoad := buildPageLoadStruct(d)
	test, err := client.CreatePageLoad(*pageLoad)
	if err != nil {
		return err
	}
	testID := test.TestID
	d.SetId(strconv.Itoa(testID))
	return resourcePageLoadRead(d, m)
}

func buildPageLoadStruct(d *schema.ResourceData) *thousandeyes.PageLoad {
	httpServer := thousandeyes.PageLoad{
		TestName:     d.Get("name").(string),
		AuthType:     d.Get("auth_type").(string),
		URL:          d.Get("url").(string),
		Interval:     d.Get("interval").(int),
		HTTPInterval: d.Get("http_interval").(int),
		Agents:       expandAgents(d.Get("agents").([]interface{})),
	}
	if attr, ok := d.GetOk("http_version"); ok {
		httpServer.HTTPVersion = attr.(int)
	}

	return &httpServer
}
