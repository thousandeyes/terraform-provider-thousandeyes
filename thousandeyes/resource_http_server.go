package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceHTTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.HTTPServer{}, schemas, nil),
		Create: resourceHTTPServerCreate,
		Read:   resourceHTTPServerRead,
		Update: resourceHTTPServerUpdate,
		Delete: resourceHTTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create an HTTP server test. This test type measures the availability and performance of an HTTP service. For more information, see [HTTP Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#http-server-test).",
	}
	return &resource
}

func resourceHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetHTTPServer(id)
	})
}

func resourceHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceHTTPServerRead(d, m)
}

func buildHTTPServerStruct(d *schema.ResourceData) *thousandeyes.HTTPServer {
	return ResourceBuildStruct(d, &thousandeyes.HTTPServer{}).(*thousandeyes.HTTPServer)
}
