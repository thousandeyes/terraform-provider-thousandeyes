package thousandeyes

import (
	"context"
	"fmt"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceDNSSec() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.DnsSecTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceDNSSecCreate,
		Read:   resourceDNSSecRead,
		Update: resourceDNSSecUpdate,
		Delete: resourceDNSSecDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a DNSSEC test. This test type verifies the digital signature of DNS resource records and validates the authenticity of those records. For more information, see [DNSSEC Test](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#dnssec-test).",
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    schemas.LegacyTestSchema().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.LegacyTestStateUpgrade,
				Version: 0,
			},
		},
		SchemaVersion: 1,
	}
	return &resource
}

func resourceDNSSecRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.DNSSECTestsAPIService)(&apiClient.Common)

		req := api.GetDnsSecTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceDNSSecUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSSECTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildDNSSecStruct(d)
	// it will keep previous record type if isn't set
	if !checkDomainRecordTypeExists(update.Domain) {
		update.Domain = fmt.Sprintf("%s ANY", update.Domain)
	}

	req := api.UpdateDnsSecTest(d.Id()).DnsSecTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceDNSSecRead(d, m)
}

func resourceDNSSecDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSSECTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteDnsSecTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceDNSSecCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.DNSSECTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildDNSSecStruct(d)

	req := api.CreateDnsSecTest().DnsSecTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceDNSSecRead(d, m)
}

func buildDNSSecStruct(d *schema.ResourceData) *tests.DnsSecTestRequest {
	return ResourceBuildStruct(d, &tests.DnsSecTestRequest{})
}
