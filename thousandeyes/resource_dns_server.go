package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func resourceDNSServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.DNSServer{}, schemas),
		Create: resourceDNSServerCreate,
		Read:   resourceDNSServerRead,
		Update: resourceDNSServerUpdate,
		Delete: resourceDNSServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	return &resource
}

func resourceDNSServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetDNSServer(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceDNSServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.DNSServer{}).(*thousandeyes.DNSServer)
	_, err := client.UpdateDNSServer(id, *update)
	if err != nil {
		return err
	}
	return resourceDNSServerRead(d, m)
}

func resourceDNSServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteDNSServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSServerStruct(d)
	remote, err := client.CreateDNSServer(*local)
	if err != nil {
		return err
	}
	id := remote.TestID
	d.SetId(strconv.Itoa(id))
	return resourceDNSServerRead(d, m)
}

func buildDNSServerStruct(d *schema.ResourceData) *thousandeyes.DNSServer {
	return ResourceBuildStruct(d, &thousandeyes.DNSServer{}).(*thousandeyes.DNSServer)
}
