package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceDNSServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.DNSServer{}, schemas, nil),
		Create: resourceDNSServerCreate,
		Read:   resourceDNSServerRead,
		Update: resourceDNSServerUpdate,
		Delete: resourceDNSServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows users to create a DNS server test. This test type validates DNS records and provides service performance metrics. For more information, see [DNS Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#dns-server-test).",
	}
	return &resource
}

func resourceDNSServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		return client.GetDNSServer(id)
	})
}

func resourceDNSServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
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
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceDNSServerRead(d, m)
}

func buildDNSServerStruct(d *schema.ResourceData) *thousandeyes.DNSServer {
	return ResourceBuildStruct(d, &thousandeyes.DNSServer{}).(*thousandeyes.DNSServer)
}
