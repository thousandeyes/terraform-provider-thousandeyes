package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceDNSSec() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.DNSSec{}, schemas, nil),
		Create: resourceDNSSecCreate,
		Read:   resourceDNSSecRead,
		Update: resourceDNSSecUpdate,
		Delete: resourceDNSSecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a DNSSEC test. This test type verifies the digital signature of DNS resource records and validates the authenticity of those records. For more information, see [DNSSEC Test](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#dnssec-test).",
	}
	return &resource
}

func resourceDNSSecRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetDNSSec(id)
	})
}

func resourceDNSSecUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.DNSSec{}).(*thousandeyes.DNSSec)
	_, err := client.UpdateDNSSec(id, *update)
	if err != nil {
		return err
	}
	return resourceDNSSecRead(d, m)
}

func resourceDNSSecDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteDNSSec(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSSecCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSSecStruct(d)
	remote, err := client.CreateDNSSec(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceDNSSecRead(d, m)
}

func buildDNSSecStruct(d *schema.ResourceData) *thousandeyes.DNSSec {
	return ResourceBuildStruct(d, &thousandeyes.DNSSec{}).(*thousandeyes.DNSSec)
}
