package thousandeyes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

const (
	endpointAgentSelectorAllAgents       = schemas.EndpointAgentSelectorAllAgents
	endpointAgentSelectorSpecificAgents  = schemas.EndpointAgentSelectorSpecificAgents
	endpointAgentSelectorAgentTags       = schemas.EndpointAgentSelectorAgentTags
	endpointAgentToServerTestIDAttribute = "test_id"
)

func resourceEndpointAgentToServer() *schema.Resource {
	return &schema.Resource{
		Create:        resourceEndpointAgentToServerCreate,
		Read:          resourceEndpointAgentToServerRead,
		Update:        resourceEndpointAgentToServerUpdate,
		Delete:        resourceEndpointAgentToServerDelete,
		CustomizeDiff: validateEndpointTestSelector,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema:      schemas.EndpointAgentToServerSchema,
		Description: "Manages a scheduled Endpoint Agent-to-Server test. Use `agent_selector_type` to choose all Endpoint Agents, specific agents, or agents selected by tags.",
	}
}

func resourceEndpointAgentToServerCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointtests.AgentToServerEndpointScheduledTestsAPIService)(&apiClient.Common)
	request, err := buildEndpointAgentToServerCreateRequest(d)
	if err != nil {
		return err
	}
	log.Printf("[INFO] Creating ThousandEyes Endpoint Agent-to-Server test")
	req := api.CreateAgentToServerEndpointScheduledTest().EndpointAgentToServerTestRequest(*request)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	response, _, err := req.Execute()
	if err != nil {
		return err
	}
	if response == nil || response.TestId == nil || *response.TestId == "" {
		return fmt.Errorf("endpoint agent-to-server create response did not include a test ID")
	}
	d.SetId(*response.TestId)

	// These two properties are only accepted by the SDK update model. Apply
	// explicitly configured values after creation, since the API defaults are
	// otherwise retained.
	createUpdate := endpointAgentToServerCreateUpdate(d)
	if endpointAgentToServerUpdateHasFields(createUpdate) {
		if err := applyEndpointAgentToServerUpdate(d, meta, createUpdate); err != nil {
			return err
		}
		return nil
	}
	return resourceEndpointAgentToServerRead(d, meta)
}

func buildEndpointAgentToServerCreateRequest(d *schema.ResourceData) (*endpointtests.EndpointAgentToServerTestRequest, error) {
	selector := d.Get("agent_selector_type").(string)
	if err := validateEndpointAgentSelector(selector, d.Get("agents").(*schema.Set), d.Get("tag_ids").(*schema.Set)); err != nil {
		return nil, err
	}
	request := &endpointtests.EndpointAgentToServerTestRequest{
		AgentSelectorType: endpointtests.EndpointTestAgentSelectorType(selector).Ptr(),
		TestName:          d.Get("test_name").(string),
		Server:            d.Get("server").(string),
	}
	if value, ok := d.GetOk("max_machines"); ok {
		maxMachines := int32(value.(int))
		request.MaxMachines = &maxMachines
	}
	switch selector {
	case endpointAgentSelectorSpecificAgents:
		request.Agents = resourceDataStringSet(d.Get("agents").(*schema.Set))
	case endpointAgentSelectorAgentTags:
		request.TagIds = resourceDataStringSet(d.Get("tag_ids").(*schema.Set))
	}
	if value, ok := d.GetOk("port"); ok {
		port := int32(value.(int))
		request.Port = &port
	}
	if value, ok := d.GetOk("interval"); ok {
		interval := endpointtests.TestInterval(value.(int))
		request.Interval = &interval
	}
	if value, ok := d.GetOk("protocol"); ok {
		protocol := endpointtests.EndpointTestProtocol(value.(string))
		request.Protocol = &protocol
	}
	if value, ok := d.GetOk("ip_version"); ok {
		ipVersion := endpointtests.EndpointIpVersionIn(value.(string))
		request.IpVersion = &ipVersion
	}
	if value, ok := d.GetOk("is_prioritized"); ok {
		prioritized := value.(bool)
		request.IsPrioritized = &prioritized
	}
	return request, nil
}

func resourceDataStringSet(set *schema.Set) []string {
	values := make([]string, 0, set.Len())
	for _, value := range set.List() {
		values = append(values, value.(string))
	}
	return values
}

func validateEndpointAgentSelector(selector string, agents, tags *schema.Set) error {
	switch selector {
	case endpointAgentSelectorAllAgents:
		if agents.Len() != 0 || tags.Len() != 0 {
			return fmt.Errorf("agents and tag_ids must be empty when agent_selector_type is %q", selector)
		}
	case endpointAgentSelectorSpecificAgents:
		if agents.Len() == 0 || tags.Len() != 0 {
			return fmt.Errorf("agent_selector_type %q requires agents and does not accept tag_ids", selector)
		}
	case endpointAgentSelectorAgentTags:
		if tags.Len() == 0 || agents.Len() != 0 {
			return fmt.Errorf("agent_selector_type %q requires tag_ids and does not accept agents", selector)
		}
	default:
		return fmt.Errorf("unsupported agent_selector_type %q", selector)
	}
	return nil
}

func resourceEndpointAgentToServerRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointtests.AgentToServerEndpointScheduledTestsAPIService)(&apiClient.Common)
	log.Printf("[INFO] Reading ThousandEyes Endpoint Agent-to-Server test %s", d.Id())
	req := api.GetAgentToServerEndpointScheduledTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	response, httpResponse, err := req.Execute()
	if err != nil {
		if httpResponse != nil && httpResponse.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return err
	}
	if response == nil {
		return fmt.Errorf("received an empty response while reading Endpoint Agent-to-Server test %s", d.Id())
	}
	if response.TestId != nil && *response.TestId != "" {
		d.SetId(*response.TestId)
	}
	if err := ResourceRead(apiClient.GetConfig().Context, d, response); err != nil {
		return err
	}
	server := ""
	if response.Server != nil {
		server = *response.Server
	}
	if err := d.Set("server", server); err != nil {
		return err
	}
	return setEndpointAgentSelectorState(d, response.AgentSelectorConfig)
}

func setEndpointAgentSelectorState(d *schema.ResourceData, config *endpointtests.EndpointAgentSelectorConfig) error {
	selector := ""
	var agents, tags []string
	var maxMachines *int32
	if config != nil {
		switch {
		case config.EndpointAllAgentsSelectorConfig != nil:
			selector = endpointAgentSelectorAllAgents
			maxMachines = config.EndpointAllAgentsSelectorConfig.MaxMachines
		case config.EndpointSpecificAgentsSelectorConfig != nil:
			selector = endpointAgentSelectorSpecificAgents
			agents = config.EndpointSpecificAgentsSelectorConfig.Agents
			maxMachines = config.EndpointSpecificAgentsSelectorConfig.MaxMachines
		case config.EndpointAgentTagsSelectorConfig != nil:
			selector = endpointAgentSelectorAgentTags
			tags = config.EndpointAgentTagsSelectorConfig.TagIds
			maxMachines = config.EndpointAgentTagsSelectorConfig.MaxMachines
		case config.EndpointAgentLabelsSelectorConfig != nil:
			return fmt.Errorf("endpoint agent selector type %q is not supported; use all-agents, specific-agents, or agent-tags", "agent-labels")
		}
	}
	if selector != "" {
		if err := d.Set("agent_selector_type", selector); err != nil {
			return err
		}
	}
	if err := d.Set("agents", agents); err != nil {
		return err
	}
	if err := d.Set("tag_ids", tags); err != nil {
		return err
	}
	if maxMachines != nil {
		if err := d.Set("max_machines", *maxMachines); err != nil {
			return err
		}
	}
	return nil
}

func endpointAgentToServerCreateUpdate(d *schema.ResourceData) endpointtests.EndpointNetworkTestUpdate {
	update := endpointtests.EndpointNetworkTestUpdate{}
	// GetOk omits an explicitly configured false value, which must be applied
	// because the API enables newly created tests by default.
	//nolint:staticcheck
	if value, ok := d.GetOkExists("is_enabled"); ok && !value.(bool) {
		isEnabled := value.(bool)
		update.IsEnabled = &isEnabled
	}
	if value, ok := d.GetOk("tcp_probe_mode"); ok {
		tcpProbeMode := endpointtests.TestProbeMode(value.(string))
		update.TcpProbeMode = &tcpProbeMode
	}
	return update
}

func resourceEndpointAgentToServerUpdate(d *schema.ResourceData, meta interface{}) error {
	update := endpointtests.EndpointNetworkTestUpdate{}
	if d.HasChange("test_name") {
		value := d.Get("test_name").(string)
		update.TestName = &value
	}
	if d.HasChange("server") {
		value := d.Get("server").(string)
		update.Server = &value
	}
	if d.HasChange("port") {
		value := int32(d.Get("port").(int))
		update.Port = &value
	}
	if d.HasChange("protocol") {
		value := endpointtests.EndpointTestProtocol(d.Get("protocol").(string))
		update.Protocol = &value
	}
	if d.HasChange("interval") {
		value := endpointtests.TestInterval(d.Get("interval").(int))
		update.Interval = &value
	}
	if d.HasChange("is_enabled") {
		value := d.Get("is_enabled").(bool)
		update.IsEnabled = &value
	}
	if d.HasChange("tcp_probe_mode") {
		value := endpointtests.TestProbeMode(d.Get("tcp_probe_mode").(string))
		update.TcpProbeMode = &value
	}
	return applyEndpointAgentToServerUpdate(d, meta, update)
}

func applyEndpointAgentToServerUpdate(d *schema.ResourceData, meta interface{}, update endpointtests.EndpointNetworkTestUpdate) error {
	if !endpointAgentToServerUpdateHasFields(update) {
		return resourceEndpointAgentToServerRead(d, meta)
	}
	apiClient := meta.(*client.APIClient)
	api := (*endpointtests.AgentToServerEndpointScheduledTestsAPIService)(&apiClient.Common)
	log.Printf("[INFO] Updating ThousandEyes Endpoint Agent-to-Server test %s", d.Id())
	req := api.UpdateAgentToServerEndpointScheduledTest(d.Id()).EndpointNetworkTestUpdate(update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceEndpointAgentToServerRead(d, meta)
}

func endpointAgentToServerUpdateHasFields(update endpointtests.EndpointNetworkTestUpdate) bool {
	return update.Interval != nil || update.TestName != nil || update.Protocol != nil || update.IsEnabled != nil || update.TcpProbeMode != nil || update.Port != nil || update.Server != nil
}

func resourceEndpointAgentToServerDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointtests.AgentToServerEndpointScheduledTestsAPIService)(&apiClient.Common)
	log.Printf("[INFO] Deleting ThousandEyes Endpoint Agent-to-Server test %s", d.Id())
	req := api.DeleteAgentToServerEndpointScheduledTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	response, err := req.Execute()
	if err != nil && (response == nil || response.StatusCode != http.StatusNotFound) {
		return err
	}
	d.SetId("")
	return nil
}
