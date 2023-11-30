package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceBGP() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.BGP{}, schemas, nil),
		Create: resourceBGPCreate,
		Read:   resourceBGPRead,
		Update: resourceBGPUpdate,
		Delete: resourceBGPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a ThousandEyes BGP test. This test type collects BGP routing related information and presents a visualization of the route and relevants events on the timeline. For more information, see [BGP Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#bgp-test).",
	}
	return &resource
}

func resourceBGPRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetBGP(id)
	})
}

func resourceBGPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.BGP{}).(*thousandeyes.BGP)
	_, err := client.UpdateBGP(id, *update)
	if err != nil {
		return err
	}
	return resourceBGPRead(d, m)
}

func resourceBGPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteBGP(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceBGPCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildBGPStruct(d)
	remote, err := client.CreateBGP(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceBGPRead(d, m)
}

func buildBGPStruct(d *schema.ResourceData) *thousandeyes.BGP {
	return ResourceBuildStruct(d, &thousandeyes.BGP{}).(*thousandeyes.BGP)
}
