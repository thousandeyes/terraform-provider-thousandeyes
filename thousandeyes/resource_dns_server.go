package thousandeyes

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceDNSServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.DnsServerTest{}, schemas, nil),
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
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.DNSServerTestsAPIService)(&apiClient.Common)

		req := api.GetDnsServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceDNSServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.DnsServerTestRequest{})

	req := api.UpdateDnsServerTest(d.Id()).DnsServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceDNSServerRead(d, m)
}

func resourceDNSServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteDnsServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSServerStruct(d)

	req := api.CreateDnsServerTest().DnsServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceDNSServerRead(d, m)
}

func buildDNSServerStruct(d *schema.ResourceData) *tests.DnsServerTestRequest {
	return ResourceBuildStruct(d, &tests.DnsServerTestRequest{})
}
