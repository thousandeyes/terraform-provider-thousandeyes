package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourcePageLoad() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.PageLoad{}, schemas, nil),
		Create: resourcePageLoadCreate,
		Read:   resourcePageLoadRead,
		Update: resourcePageLoadUpdate,
		Delete: resourcePageLoadDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a page load test. This test type obtains in-browser site performance metrics. For more information, see [Page Load Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#page-load-test).",
	}
	return &resource
}

func resourcePageLoadRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetPageLoad(id)
	})
}

func resourcePageLoadUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.PageLoad{}).(*thousandeyes.PageLoad)
	_, err := client.UpdatePageLoad(id, *update)
	if err != nil {
		return err
	}
	return resourcePageLoadRead(d, m)
}

func resourcePageLoadDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeletePageLoad(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourcePageLoadCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildPageLoadStruct(d)
	remote, err := client.CreatePageLoad(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourcePageLoadRead(d, m)
}

func buildPageLoadStruct(d *schema.ResourceData) *thousandeyes.PageLoad {
	return ResourceBuildStruct(d, &thousandeyes.PageLoad{}).(*thousandeyes.PageLoad)
}
