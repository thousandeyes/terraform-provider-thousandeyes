package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// EndpointAgentToServerDynamicSchema returns a new schema owned by the dynamic Endpoint Agent-to-Server resource.
func EndpointAgentToServerDynamicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"application": {
			Type:         schema.TypeString,
			Description:  "Application to monitor: webex, zoom, or microsoft-teams.",
			Required:     true,
			ValidateFunc: validation.StringInSlice([]string{"webex", "zoom", "microsoft-teams"}, false),
		},
		"test_name": {
			Type:         schema.TypeString,
			Description:  "Name of the test.",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"agent_selector_type": {
			Type:         schema.TypeString,
			Description:  "Agent selector type.",
			Optional:     true,
			ForceNew:     true,
			Default:      "all-agents",
			ValidateFunc: validation.StringInSlice([]string{"all-agents", "specific-agents", "agent-tags"}, false),
		},
		"agents": {
			Type:        schema.TypeSet,
			Description: "Endpoint Agent IDs. Required only for specific-agents selection.",
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"has_path_trace_in_session": {
			Type:        schema.TypeBool,
			Description: "Enables \"in session\" path trace. When enabled, this option initiates a TCP session with the target server and sends path trace packets within the established TCP session.",
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"interval": {
			Type:         schema.TypeInt,
			Description:  "Test execution interval in seconds.",
			Optional:     true,
			Default:      60,
			ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
		},
		"ip_version": {
			Type:         schema.TypeString,
			Description:  "Endpoint IP preference: V4_ONLY, V6_ONLY, V6_PREFER, or OS_DEFAULT.",
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"V4_ONLY", "V6_ONLY", "V6_PREFER", "OS_DEFAULT"}, false),
		},
		"is_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether the test is enabled.",
			Optional:    true,
			Default:     true,
		},
		"is_prioritized": {
			Type:        schema.TypeBool,
			Description: "Indicates whether the test should be prioritized when the number of tests assigned to an agent exceeds the license limit.",
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"max_machines": {
			Type:         schema.TypeInt,
			Description:  "Maximum number of agents that can execute the test. Defaults to 25.",
			Optional:     true,
			ForceNew:     true,
			Default:      25,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"protocol": {
			Type:         schema.TypeString,
			Description:  "Protocol.",
			Optional:     true,
			Default:      "icmp",
			ValidateFunc: validation.StringInSlice([]string{"icmp", "icmp-with-tcp-connect", "tcp", "prefer-tcp", "autodetect"}, false),
		},
		"tag_ids": {
			Type:        schema.TypeSet,
			Description: "Dynamic Endpoint Agent tag IDs. Required only for agent-tags selection.",
			Optional:    true,
			ForceNew:    true,
			Elem: &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
		"tcp_probe_mode": {
			Type:         schema.TypeString,
			Description:  "Tcp probe mode.",
			Optional:     true,
			Default:      "auto",
			ValidateFunc: validation.StringInSlice([]string{"auto", "sack", "syn"}, false),
		},

		"aid": {
			Type:        schema.TypeString,
			Description: "Account group ID.",
			Computed:    true,
		},
		"created_date": {
			Type:        schema.TypeString,
			Description: "Creation timestamp in RFC 3339 format.",
			Computed:    true,
		},
		"has_ping": {
			Type:        schema.TypeBool,
			Description: "Whether the test runs ping.",
			Computed:    true,
		},
		"has_traceroute": {
			Type:        schema.TypeBool,
			Description: "Whether the test runs traceroute.",
			Computed:    true,
		},
		"labels": {
			Type:        schema.TypeSet,
			Description: "IDs of legacy test labels assigned to the test. These are not Endpoint Agent labels and cannot be used for agent selection.",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"link": {
			Type:        schema.TypeString,
			Description: "URL of the test API resource.",
			Computed:    true,
		},
		"modified_date": {
			Type:        schema.TypeString,
			Description: "Last modification timestamp in RFC 3339 format.",
			Computed:    true,
		},
		"network_measurements": {
			Type:        schema.TypeBool,
			Description: "Whether network measurements are enabled.",
			Computed:    true,
		},
		"test_id": {
			Type:        schema.TypeString,
			Description: "Test ID.",
			Computed:    true,
		},
	}
}
