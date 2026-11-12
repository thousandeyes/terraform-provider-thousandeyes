package schemas

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// EndpointHTTPServerSchema returns a new schema owned by the scheduled Endpoint HTTP Server resource.
func EndpointHTTPServerSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
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
		"max_machines": {
			Type:         schema.TypeInt,
			Description:  "Maximum number of agents that can execute the test. Defaults to 25.",
			Optional:     true,
			ForceNew:     true,
			Default:      25,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"test_name": {
			Type:         schema.TypeString,
			Description:  "Name of the test.",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		"ip_version": {
			Type:         schema.TypeString,
			Description:  "Endpoint IP preference: V4_ONLY, V6_ONLY, V6_PREFER, or OS_DEFAULT.",
			Optional:     true,
			ForceNew:     true,
			Computed:     true,
			ValidateFunc: validation.StringInSlice([]string{"V4_ONLY", "V6_ONLY", "V6_PREFER", "OS_DEFAULT"}, false),
		},
		"auth_type": {
			Type:         schema.TypeString,
			Description:  "Auth type.",
			Optional:     true,
			ForceNew:     true,
			Default:      "none",
			ValidateFunc: validation.StringInSlice([]string{"none", "basic", "ntlm"}, false),
		},
		"has_path_trace_in_session": {
			Type:        schema.TypeBool,
			Description: "Enables \"in session\" path trace. When enabled, this option initiates a TCP session with the target server and sends path trace packets within the established TCP session.",
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"http_time_limit": {
			Type:         schema.TypeInt,
			Description:  "HTTP request timeout in milliseconds.",
			Optional:     true,
			ForceNew:     true,
			Default:      5000,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"protocol": {
			Type:         schema.TypeString,
			Description:  "Protocol.",
			Optional:     true,
			Default:      "icmp",
			ValidateFunc: validation.StringInSlice([]string{"icmp", "icmp-with-tcp-connect", "tcp", "prefer-tcp", "autodetect"}, false),
		},
		"username": {
			Type:        schema.TypeString,
			Description: "Username for Basic/NTLM authentication.",
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"ssl_version_id": {
			Type:         schema.TypeString,
			Description:  "Ssl version id.",
			Optional:     true,
			ForceNew:     true,
			Default:      "0",
			ValidateFunc: validation.StringInSlice([]string{"0", "3", "4", "5", "6", "7"}, false),
		},
		"tcp_probe_mode": {
			Type:         schema.TypeString,
			Description:  "Tcp probe mode.",
			Optional:     true,
			Default:      "auto",
			ValidateFunc: validation.StringInSlice([]string{"auto", "sack", "syn"}, false),
		},
		"verify_certificate": {
			Type:        schema.TypeBool,
			Description: "Flag indicating if a certificate should be verified.",
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},
		"url": {
			Type:         schema.TypeString,
			Description:  "The test target URL. HTTPS is used when the scheme is omitted. Changes require replacement when auth_type is basic or ntlm.",
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				return NormalizeEndpointHTTPServerURL(old) == NormalizeEndpointHTTPServerURL(new)
			},
			DiffSuppressOnRefresh: true,
		},
		"network_measurements": {
			Type:        schema.TypeBool,
			Description: "Enable or disable network measurements. Set to true to enable or false to disable network measurements.",
			Optional:    true,
			ForceNew:    true,
			Default:     true,
		},
		"target_response_time": {
			Type:         schema.TypeInt,
			Description:  "Response time target in milliseconds.",
			Optional:     true,
			ForceNew:     true,
			Default:      1000,
			ValidateFunc: validation.IntBetween(1, 2147483647),
		},
		"password": {
			Type:        schema.TypeString,
			Description: "Password for Basic or NTLM authentication.",
			Optional:    true,
			ForceNew:    true,
			Sensitive:   true,
		},
		"is_prioritized": {
			Type:        schema.TypeBool,
			Description: "Indicates whether the test should be prioritized when the number of tests assigned to an agent exceeds the license limit.",
			Optional:    true,
			ForceNew:    true,
			Default:     false,
		},
		"interval": {
			Type:         schema.TypeInt,
			Description:  "Test execution interval in seconds.",
			Optional:     true,
			Default:      60,
			ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
		},
		"is_enabled": {
			Type:        schema.TypeBool,
			Description: "Whether the test is enabled.",
			Optional:    true,
			Default:     true,
		},

		"link": {
			Type:        schema.TypeString,
			Description: "URL of the test API resource.",
			Computed:    true,
		},
		"labels": {
			Type:        schema.TypeSet,
			Description: "IDs of legacy test labels assigned to the test. These are not Endpoint Agent labels and cannot be used for agent selection.",
			Computed:    true,
			Elem:        &schema.Schema{Type: schema.TypeString},
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
		"is_saved_event": {
			Type:        schema.TypeBool,
			Description: "Whether the test is a private snapshot.",
			Computed:    true,
		},
		"modified_date": {
			Type:        schema.TypeString,
			Description: "Last modification timestamp in RFC 3339 format.",
			Computed:    true,
		},
		"server": {
			Type:        schema.TypeString,
			Description: "Target domain name or IP address.",
			Computed:    true,
		},
		"test_id": {
			Type:        schema.TypeString,
			Description: "Test ID.",
			Computed:    true,
		},
		"type": {
			Type:        schema.TypeString,
			Description: "Type of test being queried.",
			Computed:    true,
		},
		"port": {
			Type:        schema.TypeInt,
			Description: "Port number.",
			Computed:    true,
		},
		"follow_redirects": {
			Type:        schema.TypeBool,
			Description: "Whether the test follows HTTP redirects.",
			Computed:    true,
		},
		"http_target_time": {
			Type:        schema.TypeInt,
			Description: "HTTP response time target in milliseconds.",
			Computed:    true,
		},
		"http_version": {
			Type:        schema.TypeInt,
			Description: "HTTP protocol version used by the test.",
			Computed:    true,
		},
		"ssl_version": {
			Type:        schema.TypeString,
			Description: "SSL/TLS protocol version used by the test.",
			Computed:    true,
		},
		"use_ntlm": {
			Type:        schema.TypeBool,
			Description: "Whether HTTP authentication uses NTLM.",
			Computed:    true,
		},
	}
}

// NormalizeEndpointHTTPServerURL applies the Endpoint API's default HTTPS scheme.
func NormalizeEndpointHTTPServerURL(value string) string {
	if value == "" || strings.Contains(value, "://") {
		return value
	}
	return "https://" + value
}
