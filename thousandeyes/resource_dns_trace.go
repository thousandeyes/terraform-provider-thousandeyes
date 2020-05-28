package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
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
	}
	return &resource
}

func resourceDNSTraceRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetDNSTrace(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceDNSTraceUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
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
	id, _ := strconv.Atoi(d.Id())
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
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceDNSTraceRead(d, m)
}

func buildDNSTraceStruct(d *schema.ResourceData) *thousandeyes.DNSTrace {
	return ResourceBuildStruct(d, &thousandeyes.DNSTrace{}).(*thousandeyes.DNSTrace)
}
