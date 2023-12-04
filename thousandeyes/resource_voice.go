package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceRTPStream() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.RTPStream{}, schemas, nil),
		Create: resourceRTPStreamCreate,
		Read:   resourceRTPStreamRead,
		Update: resourceRTPStreamUpdate,
		Delete: resourceRTPStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a RTP Stream test. This test type measures the quality of real-time protocol (RTP) voice streams between ThousandEyes agents that act as VoIP user agents. For more information, see [RTP Stream Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#rtp-stream-test).",
	}
	return &resource
}

func resourceRTPStreamRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetRTPStream(id)
	})
}

func resourceRTPStreamUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.RTPStream{}).(*thousandeyes.RTPStream)
	_, err := client.UpdateRTPStream(id, *update)
	if err != nil {
		return err
	}
	return resourceRTPStreamRead(d, m)
}

func resourceRTPStreamDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteRTPStream(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceRTPStreamCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildRTPStreamStruct(d)
	remote, err := client.CreateRTPStream(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceRTPStreamRead(d, m)
}

func buildRTPStreamStruct(d *schema.ResourceData) *thousandeyes.RTPStream {
	return ResourceBuildStruct(d, &thousandeyes.RTPStream{}).(*thousandeyes.RTPStream)
}
