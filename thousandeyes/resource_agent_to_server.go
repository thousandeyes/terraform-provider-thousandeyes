package thousandeyes

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceAgentToServer() *schema.Resource {
	agentToServerSchemasOverride := map[string]*schema.Schema{
		"port": {
			Type:         schema.TypeInt,
			Description:  "The target port.",
			ValidateFunc: validation.IntBetween(1, 65535),
			Optional:     true,
			DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
				// allows port to be optional even for TCP tests
				// it uses ThousandEyes' default which is 80
				return newValue == "0"
			},
		},
	}

	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.AgentToServerTest{}, schemas, agentToServerSchemasOverride),
		Create: resourceAgentServerCreate,
		Read:   resourceAgentServerRead,
		Update: resourceAgentServerUpdate,
		Delete: resourceAgentServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create and configure an agent-to-server test. This test type measures network performance as seen from ThousandEyes agent(s) towards a remote server. For more information about agent-to-server tests, see [Agent-to-Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-server-test).",
	}
	return &resource
}

func resourceAgentServerRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.AgentToServerTestsAPIService)(&apiClient.Common)

		req := api.GetAgentToServerTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceAgentServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.AgentToServerTestRequest{})

	req := api.UpdateAgentToServerTest(d.Id()).AgentToServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceAgentServerRead(d, m)
}

func resourceAgentServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteAgentToServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAgentServerStruct(d)

	req := api.CreateAgentToServerTest().AgentToServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceAgentServerRead(d, m)
}

func buildAgentServerStruct(d *schema.ResourceData) *tests.AgentToServerTestRequest {
	return ResourceBuildStruct(d, &tests.AgentToServerTestRequest{})
}
