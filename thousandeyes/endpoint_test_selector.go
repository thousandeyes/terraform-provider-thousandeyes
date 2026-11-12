package thousandeyes

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func validateEndpointTestSelector(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
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
	if !endpointTestSelectorValid(mode, agents, tags) {
		return fmt.Errorf("agent_selector_type %q requires only agents for specific-agents, only tag_ids for agent-tags, and neither for all-agents", mode)
	}
	return nil
}

func flattenEndpointTestSelector(d *schema.ResourceData, s *endpointtests.EndpointAgentSelectorConfig) error {
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

func endpointTestSelectorValid(mode string, agents, tags int) bool {
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
