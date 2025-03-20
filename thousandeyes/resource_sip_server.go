package thousandeyes

import (
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceSIPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.SipServerTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceSIPServerCreate,
		Read:   resourceSIPServerRead,
		Update: resourceSIPServerUpdate,
		Delete: resourceSIPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create a SIP server test. This test type checks for the availability and performance of a VoIP SIP server, confirms the ability to perform SIP Register with a target server, and observes the requests and responses. For more information, see [SIP Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#sip-server-test).",
	}
	resource.Schema["protocol"] = schemas.CommonSchema["protocol-sip"]
	return &resource
}

func resourceSIPServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.SIPServerTestsAPIService)(&apiClient.Common)

		req := api.GetSipServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceSIPServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.SIPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.SipServerTestRequest{})
	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, SIP Server tests actually require a few fields
	// within the targetSipCredentials object to be retained on update.
	// Calls without port, protocol, or sipRegistrar will fail, whereas sipProxy
	// being absent will cause the update to remove the  value.
	// Unlike other cases, we can send all non-updated values within targetSipCredentials
	// without being rejected.
	fullUpdate := buildSIPServerStruct(d)
	update.TargetSipCredentials = fullUpdate.TargetSipCredentials
	req := api.UpdateSipServerTest(d.Id()).SipServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, _, err := req.Execute(); err != nil {
		return err
	}
	return resourceSIPServerRead(d, m)
}

func resourceSIPServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.SIPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteSipServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceSIPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.SIPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildSIPServerStruct(d)

	req := api.CreateSipServerTest().SipServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceSIPServerRead(d, m)
}

func buildSIPServerStruct(d *schema.ResourceData) *tests.SipServerTestRequest {
	return ResourceBuildStruct(d, &tests.SipServerTestRequest{})
}
