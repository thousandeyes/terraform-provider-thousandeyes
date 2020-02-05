package thousandeyes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
	"log"
	"strconv"
)

func resourceHttpServer() *schema.Resource {
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
		Create: resourceHttpServerCreate,
		Read:   resourceHttpServerRead,
		Update: resourceHttpServerUpdate,
		Delete: resourceHttpServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceHttpServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetHttpServer(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("auth_type", test.AuthType)
	d.Set("interval", test.Interval)
	d.Set("http_version", test.HttpVersion)
	d.Set("url", test.Url)
	d.Set("agents", test.Agents)
	return nil
}

func resourceHttpServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.HttpServer
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
		update.HttpVersion = d.Get("http_version").(int)
	}
	if d.HasChange("url") {
		update.Url = d.Get("url").(string)
	}
	_, err := client.UpdateHttpServer(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceHttpServerRead(d, m)
}

func resourceHttpServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteHttpServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceHttpServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	httpServer := buildHttpServerStruct(d)
	httpTest, err := client.CreateHttpServer(*httpServer)
	if err != nil {
		return err
	}
	testId := httpTest.TestId
	d.SetId(strconv.Itoa(testId))
	return resourceHttpServerRead(d, m)
}

func buildHttpServerStruct(d *schema.ResourceData) *thousandeyes.HttpServer {
	httpServer := thousandeyes.HttpServer{
		TestName:    d.Get("name").(string),
		AuthType:    d.Get("auth_type").(string),
		HttpVersion: d.Get("http_version").(int),
		Url:         d.Get("url").(string),
		Interval:    d.Get("interval").(int),
		Agents:      expandAgents(d.Get("agents").([]interface{})),
	}
	if attr, ok := d.GetOk("http_version"); ok {
		httpServer.HttpVersion = attr.(int)
	}

	return &httpServer
}

func expandAgents(v interface{}) thousandeyes.Agents {
	var agents thousandeyes.Agents

	for _, er := range v.([]interface{}) {
		rer := er.(map[string]interface{})
		agent := &thousandeyes.Agent{
			AgentId: rer["agent_id"].(int),
		}
		agents = append(agents, *agent)
	}

	return agents
}
