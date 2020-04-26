package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceFTPServer() *schema.Resource {
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
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
			},
			"url": {
				Type:        schema.TypeString,
				Description: "target for the test",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "username to be used to authenticate with the destination server",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "password to be used to authenticate with the destination server",
			},
			"request_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Set the type of activity for the test: Download, Upload, or List",
			},
		},
		Create: resourceFTPServerCreate,
		Read:   resourceFTPServerRead,
		Update: resourceFTPServerUpdate,
		Delete: resourceFTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceFTPServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetFTPServer(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("agents", test.Agents)
	d.Set("interval", test.Interval)
	d.Set("password", test.Password)
	d.Set("username", test.Username)
	d.Set("url", test.URL)
	return nil
}

func resourceFTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.FTPServer
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("password") {
		update.Password = d.Get("password").(string)
	}
	if d.HasChange("requestType") {
		update.RequestType = d.Get("request_type").(string)
	}
	if d.HasChange("username") {
		update.Username = d.Get("username").(string)
	}
	if d.HasChange("url") {
		update.URL = d.Get("url").(string)
	}
	_, err := client.UpdateFTPServer(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceFTPServerRead(d, m)
}

func resourceFTPServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteFTPServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceFTPServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	ftpServer := buildFTPServerStruct(d)
	httpTest, err := client.CreateFTPServer(*ftpServer)
	if err != nil {
		return err
	}
	testID := httpTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceFTPServerRead(d, m)
}

func buildFTPServerStruct(d *schema.ResourceData) *thousandeyes.FTPServer {
	ftpServer := thousandeyes.FTPServer{
		Agents:      expandAgents(d.Get("agents").([]interface{})),
		TestName:    d.Get("name").(string),
		Interval:    d.Get("interval").(int),
		Password:    d.Get("password").(string),
		Username:    d.Get("username").(string),
		RequestType: d.Get("request_type").(string),
		URL:         d.Get("url").(string),
	}

	return &ftpServer
}
