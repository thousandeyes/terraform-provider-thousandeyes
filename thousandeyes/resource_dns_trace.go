package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceDNSTrace() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.DNSTrace{}, schemas),
		Create: resourceDNSTraceCreate,
		Read:   resourceDNSTraceRead,
		Update: resourceDNSTraceUpdate,
		Delete: resourceDNSTraceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource provides users with the ability to create a DNS trace test. This test type verifies the delegation of DNS records and ensures the DNS hierarchy is as expected. For more information, see [DNS Trace Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#dns-trace-test).",
	}
	return &resource
}

func resourceDNSTraceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	remote, err := client.GetDNSTrace(id)
	if err != nil {
		return err
	}
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceDNSTraceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.DNSTrace{}).(*thousandeyes.DNSTrace)
	_, err := client.UpdateDNSTrace(id, *update)
	if err != nil {
		return err
	}
	return resourceDNSTraceRead(d, m)
}

func resourceDNSTraceDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteDNSTrace(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSTraceCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSTraceStruct(d)
	remote, err := client.CreateDNSTrace(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceDNSTraceRead(d, m)
}

func buildDNSTraceStruct(d *schema.ResourceData) *thousandeyes.DNSTrace {
	return ResourceBuildStruct(d, &thousandeyes.DNSTrace{}).(*thousandeyes.DNSTrace)
}
