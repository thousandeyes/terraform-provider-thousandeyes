package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	EndpointAgentSelectorAllAgents      = "all-agents"
	EndpointAgentSelectorSpecificAgents = "specific-agents"
	EndpointAgentSelectorAgentTags      = "agent-tags"
)

var EndpointAgentToServerSchema = map[string]*schema.Schema{
	"agent_selector_type": {
		Type:         schema.TypeString,
		Required:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{EndpointAgentSelectorAllAgents, EndpointAgentSelectorSpecificAgents, EndpointAgentSelectorAgentTags}, false),
		Description:  "How Endpoint Agents are selected for this test: `all-agents`, `specific-agents`, or `agent-tags`.",
	},
	"agents": {
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Endpoint Agent IDs to use when `agent_selector_type` is `specific-agents`.",
	},
	"tag_ids": {
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Elem:        &schema.Schema{Type: schema.TypeString},
		Description: "Tag IDs used to select Endpoint Agents when `agent_selector_type` is `agent-tags`.",
	},
	"max_machines": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ForceNew:     true,
		ValidateFunc: validation.IntAtLeast(1),
		Description:  "Maximum number of Endpoint Agents that can execute the test. If omitted, the API default of 25 applies.",
	},
	"test_name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the test.",
	},
	"server": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Target domain name or IP address.",
	},
	"port": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntBetween(1, 65535),
		Description:  "Target port. The API default is 443.",
	},
	"protocol": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"icmp", "icmp-with-tcp-connect", "tcp", "prefer-tcp", "autodetect"}, false),
		Description:  "Protocol requested for the test.",
	},
	"interval": {
		Type:         schema.TypeInt,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
		Description:  "Interval between test runs in seconds.",
	},
	"ip_version": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ForceNew:     true,
		ValidateFunc: validation.StringInSlice([]string{"V4_ONLY", "V6_ONLY", "V6_PREFER", "OS_DEFAULT"}, false),
		Description:  "IP version used for the test. This setting requires replacement when changed.",
	},
	"is_prioritized": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		ForceNew:    true,
		Description: "Whether the test is prioritized when the Endpoint Agent's test limit is exceeded. This setting requires replacement when changed.",
	},
	"is_enabled": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Whether the test is enabled. The API enables tests by default.",
	},
	"tcp_probe_mode": {
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		ValidateFunc: validation.StringInSlice([]string{"auto", "sack", "syn"}, false),
		Description:  "Probe mode used when the protocol is TCP.",
	},
	"network_measurements": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether network measurements are enabled.",
	},
	"has_path_trace_in_session": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether path tracing runs within a TCP session.",
	},
	"is_saved_event": {
		Type:        schema.TypeBool,
		Computed:    true,
		Description: "Whether the test is a saved event (private snapshot).",
	},
	"created_date": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UTC date and time when the test was created.",
	},
	"modified_date": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "UTC date and time when the test was last modified.",
	},
	"type": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The type of test returned by the API.",
	},
	"aid": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Account group ID associated with the test.",
	},
	"link": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Link to the test resource returned by the API.",
	},
	"test_id": {
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The unique ID of the test.",
	},
}
