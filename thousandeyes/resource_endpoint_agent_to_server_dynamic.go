package thousandeyes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func resourceEndpointAgentToServerDynamic() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointAgentToServerDynamicCreate,
		Read:   resourceEndpointAgentToServerDynamicRead,
		Update: resourceEndpointAgentToServerDynamicUpdate,
		Delete: resourceEndpointAgentToServerDynamicDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description:   "Manages a dynamic Endpoint Agent-to-Server test. Agent selectors and other create-only settings require replacement when changed.",
		CustomizeDiff: resourceEndpointAgentToServerDynamicCustomizeDiff,
		Schema:        schemas.EndpointAgentToServerDynamicSchema(),
	}
}

func resourceEndpointAgentToServerDynamicCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	return validateEndpointTestSelector(ctx, d, meta)
}

func buildEndpointAgentToServerDynamicCreate(d *schema.ResourceData) *endpointtests.DynamicTestRequest {
	return ResourceBuildStruct(d, &endpointtests.DynamicTestRequest{})
}

func buildEndpointAgentToServerDynamicUpdate(d *schema.ResourceData) *endpointtests.EndpointDynamicTestUpdate {
	return ResourceBuildStruct(d, &endpointtests.EndpointDynamicTestUpdate{})
}

func endpointAgentToServerDynamicAPI(m interface{}) *endpointtests.AgentToServerEndpointDynamicTestsAPIService {
	return (*endpointtests.AgentToServerEndpointDynamicTestsAPIService)(&m.(*client.APIClient).Common)
}

func resourceEndpointAgentToServerDynamicCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	req := endpointAgentToServerDynamicAPI(m).CreateAgentToServerEndpointDynamicTest().
		DynamicTestRequest(*buildEndpointAgentToServerDynamicCreate(d))
	resp, _, err := SetAidFromContext(apiClient.GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	if resp.GetTestId() == "" {
		return fmt.Errorf("created endpoint agent-to-server dynamic test returned no ID")
	}
	d.SetId(resp.GetTestId())
	// The create API always enables the test and does not accept is_enabled.
	// Disable it through the update API when the configuration requests it.
	if !d.Get("is_enabled").(bool) {
		return resourceEndpointAgentToServerDynamicUpdate(d, m)
	}
	return resourceEndpointAgentToServerDynamicRead(d, m)
}

func resourceEndpointAgentToServerDynamicUpdate(d *schema.ResourceData, m interface{}) error {
	req := endpointAgentToServerDynamicAPI(m).UpdateAgentToServerEndpointDynamicTest(d.Id()).
		EndpointDynamicTestUpdate(*buildEndpointAgentToServerDynamicUpdate(d))
	_, _, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	return resourceEndpointAgentToServerDynamicRead(d, m)
}

func resourceEndpointAgentToServerDynamicDelete(d *schema.ResourceData, m interface{}) error {
	req := endpointAgentToServerDynamicAPI(m).DeleteAgentToServerEndpointDynamicTest(d.Id())
	httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		return err
	}
	d.SetId("")
	return nil
}

func resourceEndpointAgentToServerDynamicRead(d *schema.ResourceData, m interface{}) error {
	req := endpointAgentToServerDynamicAPI(m).GetAgentToServerEndpointDynamicTest(d.Id())
	resp, httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if httpResp != nil && httpResp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	return flattenEndpointAgentToServerDynamic(d, resp)
}

func flattenEndpointAgentToServerDynamic(d *schema.ResourceData, resp *endpointtests.DynamicTest) error {
	if err := ResourceRead(context.Background(), d, resp); err != nil {
		return err
	}
	return flattenEndpointTestSelector(d, resp.AgentSelectorConfig)
}
