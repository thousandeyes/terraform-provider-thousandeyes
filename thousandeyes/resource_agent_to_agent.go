package thousandeyes

import (
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func resourceAgentToAgent() *schema.Resource {
	agentToAgentSchemasOverride := map[string]*schema.Schema{
		"port": {
			Type:         schema.TypeInt,
			Description:  "The target port.",
			ValidateFunc: validation.IntBetween(1, 65535),
			Default:      49153,
			Optional:     true,
		},
	}
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.AgentToAgentTestRequest{}, schemas.CommonSchema, agentToAgentSchemasOverride),
		Create: resourceAgentAgentCreate,
		Read:   resourceAgentAgentRead,
		Update: resourceAgentAgentUpdate,
		Delete: resourceAgentAgentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create and configure an agent-to-agent test. This test type evaluates the performance of the underlying network between two physical sites. For more information about agent-to-agent tests, see [Agent-to-Agent Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-agent-test).",
	}
	resource.Schema["protocol"] = schemas.CommonSchema["protocol-a2a"]
	return &resource
}

func resourceAgentAgentRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tests.AgentToAgentTestsAPIService)(&apiClient.Common)

		req := api.GetAgentToAgentTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

		resp, _, err := req.Execute()
		return resp, err
	})
}

func resourceAgentAgentUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToAgentTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := ResourceUpdate(d, &tests.AgentToAgentTestRequest{})

	req := api.UpdateAgentToAgentTest(d.Id()).AgentToAgentTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceAgentAgentRead(d, m)
}

func resourceAgentAgentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToAgentTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteAgentToAgentTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAgentAgentCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.AgentToAgentTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAgentAgentStruct(d)

	req := api.CreateAgentToAgentTest().AgentToAgentTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceAgentAgentRead(d, m)
}

func buildAgentAgentStruct(d *schema.ResourceData) *tests.AgentToAgentTestRequest {
	return ResourceBuildStruct(d, &tests.AgentToAgentTestRequest{})
}
