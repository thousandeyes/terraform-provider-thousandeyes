package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceDNSTrace() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.DnsTraceTestRequest{}, schemas.CommonSchema, nil),
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
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.DNSTraceTestsAPIService)(&apiClient.Common)

		req := api.GetDnsTraceTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceDNSTraceUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSTraceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildDNSTraceStruct(d)

	req := api.UpdateDnsTraceTest(d.Id()).DnsTraceTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceDNSTraceRead(d, m)
}

func resourceDNSTraceDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSTraceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteDnsTraceTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSTraceCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSTraceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSTraceStruct(d)

	req := api.CreateDnsTraceTest().DnsTraceTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceDNSTraceRead(d, m)
}

func buildDNSTraceStruct(d *schema.ResourceData) *tests.DnsTraceTestRequest {
	return ResourceBuildStruct(d, &tests.DnsTraceTestRequest{})
}
