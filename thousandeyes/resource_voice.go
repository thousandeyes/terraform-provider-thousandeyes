package thousandeyes

import (
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceRTPStream() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.VoiceTestRequest{}, schemas.CommonSchema, nil),
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
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.VoiceTestsAPIService)(&apiClient.Common)

		req := api.GetVoiceTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceRTPStreamUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.VoiceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.VoiceTestRequest{})

	req := api.UpdateVoiceTest(d.Id()).VoiceTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceRTPStreamRead(d, m)
}

func resourceRTPStreamDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.VoiceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteVoiceTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceRTPStreamCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.VoiceTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildRTPStreamStruct(d)

	req := api.CreateVoiceTest().VoiceTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceRTPStreamRead(d, m)
}

func buildRTPStreamStruct(d *schema.ResourceData) *tests.VoiceTestRequest {
	return ResourceBuildStruct(d, &tests.VoiceTestRequest{})
}
