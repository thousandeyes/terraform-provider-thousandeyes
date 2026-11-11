package thousandeyes

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func resourceEndpointHTTPServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointHTTPServerCreate,
		Read:   resourceEndpointHTTPServerRead,
		Update: resourceEndpointHTTPServerUpdate,
		Delete: resourceEndpointHTTPServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description:   "Manages a scheduled Endpoint HTTP Server test. Agent selectors and other create-only settings require replacement when changed.",
		CustomizeDiff: resourceEndpointHTTPServerCustomizeDiff,
		Schema:        schemas.EndpointHTTPServerSchema(),
	}
}

func resourceEndpointHTTPServerCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateEndpointHTTPServerSelector(ctx, d, meta); err != nil {
		return err
	}
	if d.NewValueKnown("ip_version") && d.NewValueKnown("protocol") {
		ipVersion := d.Get("ip_version").(string)
		if ipVersion != "" && ipVersion != "V4_ONLY" && d.Get("protocol").(string) != "icmp" {
			return fmt.Errorf("ip_version %q requires protocol = icmp; other protocols require V4_ONLY", ipVersion)
		}
	}
	return nil
}

func validateEndpointHTTPServerSelector(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	for _, key := range []string{"agent_selector_type", "agents", "tag_ids"} {
		if !d.NewValueKnown(key) || (key != "agent_selector_type" && !d.NewValueKnown(key+".#")) {
			return nil
		}
		// A set containing a newly created agent or tag ID can have a known
		// collection shape while its elements remain unknown until apply.
		raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: key}})
		if !diags.HasError() && !raw.IsWhollyKnown() {
			return nil
		}
	}
	mode := d.Get("agent_selector_type").(string)
	agents, tags := d.Get("agents").(*schema.Set).Len(), d.Get("tag_ids").(*schema.Set).Len()
	if !endpointHTTPServerSelectorValid(mode, agents, tags) {
		return fmt.Errorf("agent_selector_type %q requires only agents for specific-agents, only tag_ids for agent-tags, and neither for all-agents", mode)
	}
	return nil
}

func buildEndpointHTTPServerCreate(d *schema.ResourceData) *endpointtests.EndpointHttpServerTestRequest {
	req := ResourceBuildStruct(d, &endpointtests.EndpointHttpServerTestRequest{})
	req.Url = schemas.NormalizeEndpointHTTPServerURL(req.Url)
	return req
}

func buildEndpointHTTPServerUpdate(d *schema.ResourceData) *endpointtests.EndpointHttpTestUpdate {
	req := ResourceBuildStruct(d, &endpointtests.EndpointHttpTestUpdate{})
	if req.Url != nil {
		url := schemas.NormalizeEndpointHTTPServerURL(*req.Url)
		req.Url = &url
	}
	return req
}

func endpointHTTPServerAPI(m interface{}) *endpointtests.HTTPServerEndpointScheduledTestsAPIService {
	return (*endpointtests.HTTPServerEndpointScheduledTestsAPIService)(&m.(*client.APIClient).Common)
}

func resourceEndpointHTTPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	req := endpointHTTPServerAPI(m).CreateHttpServerEndpointScheduledTest().
		EndpointHttpServerTestRequest(*buildEndpointHTTPServerCreate(d))
	resp, _, err := SetAidFromContext(apiClient.GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	if resp.GetTestId() == "" {
		return fmt.Errorf("created endpoint HTTP server test returned no ID")
	}
	d.SetId(resp.GetTestId())
	// The create API always enables the test and does not accept is_enabled.
	// Disable it through the update API when the configuration requests it.
	if !d.Get("is_enabled").(bool) {
		return resourceEndpointHTTPServerUpdate(d, m)
	}
	return resourceEndpointHTTPServerRead(d, m)
}

func resourceEndpointHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).UpdateHttpServerEndpointScheduledTest(d.Id()).
		EndpointHttpTestUpdate(*buildEndpointHTTPServerUpdate(d))
	_, _, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	return resourceEndpointHTTPServerRead(d, m)
}

func resourceEndpointHTTPServerDelete(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).DeleteHttpServerEndpointScheduledTest(d.Id())
	httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		return err
	}
	d.SetId("")
	return nil
}

func resourceEndpointHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).GetHttpServerEndpointScheduledTest(d.Id())
	resp, httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if httpResp != nil && httpResp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	return flattenEndpointHTTPServer(d, resp)
}

func flattenEndpointHTTPServer(d *schema.ResourceData, resp *endpointtests.EndpointHttpServerTest) error {
	if err := ResourceRead(context.Background(), d, resp); err != nil {
		return err
	}
	// The create API calls this targetResponseTime, while reads return httpTargetTime.
	if resp.HttpTargetTime != nil {
		if err := d.Set("target_response_time", *resp.HttpTargetTime); err != nil {
			return err
		}
	}
	return flattenEndpointHTTPServerSelector(d, resp.AgentSelectorConfig)
}

func flattenEndpointHTTPServerSelector(d *schema.ResourceData, s *endpointtests.EndpointAgentSelectorConfig) error {
	if s == nil {
		return nil
	}
	var mode string
	var max *int32
	var agents, tags []string
	switch {
	case s.EndpointAllAgentsSelectorConfig != nil:
		v := s.EndpointAllAgentsSelectorConfig
		mode = v.AgentSelectorType
		max = v.MaxMachines
	case s.EndpointSpecificAgentsSelectorConfig != nil:
		v := s.EndpointSpecificAgentsSelectorConfig
		mode = v.AgentSelectorType
		max = v.MaxMachines
		agents = v.Agents
	case s.EndpointAgentTagsSelectorConfig != nil:
		v := s.EndpointAgentTagsSelectorConfig
		mode = v.AgentSelectorType
		max = v.MaxMachines
		tags = v.TagIds
	default:
		return fmt.Errorf("unsupported endpoint agent selector; only all-agents, specific-agents, and agent-tags are supported")
	}
	values := map[string]interface{}{"agent_selector_type": mode, "max_machines": 25}
	// Leave absent optional sets null. Writing empty sets on every refresh
	// would produce null-to-empty drift after Terraform applies the state.
	if mode == "specific-agents" || d.Get("agents").(*schema.Set).Len() > 0 {
		values["agents"] = agents
	}
	if mode == "agent-tags" || d.Get("tag_ids").(*schema.Set).Len() > 0 {
		values["tag_ids"] = tags
	}
	if max != nil {
		values["max_machines"] = *max
	}
	for k, v := range values {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}

	return nil
}

func endpointHTTPServerSelectorValid(mode string, agents, tags int) bool {
	switch mode {
	case "all-agents":
		return agents == 0 && tags == 0
	case "specific-agents":
		return agents > 0 && tags == 0
	case "agent-tags":
		return tags > 0 && agents == 0
	}
	return false
}
