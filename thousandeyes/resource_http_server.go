package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceHTTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.HTTPServer{}, schemas),
		Create: resourceHTTPServerCreate,
		Read:   resourceHTTPServerRead,
		Update: resourceHTTPServerUpdate,
		Delete: resourceHTTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetHTTPServer(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.HTTPServer{}).(*thousandeyes.HTTPServer)
	_, err := client.UpdateHTTPServer(id, *update)
	if err != nil {
		return err
	}
	return resourceHTTPServerRead(d, m)
}

func resourceHTTPServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteHTTPServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceHTTPServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildHTTPServerStruct(d)
	remote, err := client.CreateHTTPServer(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceHTTPServerRead(d, m)
}

func buildHTTPServerStruct(d *schema.ResourceData) *thousandeyes.HTTPServer {
	return ResourceBuildStruct(d, &thousandeyes.HTTPServer{}).(*thousandeyes.HTTPServer)
}
