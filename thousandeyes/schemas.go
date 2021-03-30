package thousandeyes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var schemas = map[string]*schema.Schema{
	"account_group_roles": {
		Type:        schema.TypeList,
		Description: "List of roles for user",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"account_group": {
					Type:        schema.TypeMap,
					Description: "account group for roles",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"aid": {
								Type:        schema.TypeInt,
								Description: "account group id",
								Optional:    true,
							},
						},
					},
				},
				"roles": {
					Type:        schema.TypeList,
					Description: "List of roles for user",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"role_id": {
								Type:        schema.TypeInt,
								Description: "role id",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	},
	"agents": {
		Type:        schema.TypeList,
		Description: "agents to use ",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeInt,
					Description: "agent id",
					Optional:    true,
				},
			},
		},
	},
	"agents-label": {
		Type:        schema.TypeList,
		Description: "agents to use ",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeInt,
					Description: "agent id",
					Optional:    true,
				},
			},
		},
	},
	"alert_rules": {
		Description: "get ruleId from /alert-rules endpoint. If alertsEnabled is set to 1 and alertRules is not included in a creation/update query, applicable defaults will be used.",
		Optional:    true,
		Required:    false,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rule_id": {
					Type:        schema.TypeInt,
					Description: "If alertsEnabled is set to 1 and alertRules is not included in a creation/update query, applicable defaults will be used.",
					Optional:    true,
				},
			},
		},
	},
	"alert_type": {
		Description:  "Acceptable test types, verbose names",
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Page Load", "HTTP Server", "End-to-End (Server)", "End-to-End (Agent)", "Voice", "DNS+ Domain", "DNS+ Server", "DNS Server", "DNS Trace", "DNSSEC", "Transactions", "Web Transactions", "BGP", "Path Trace", "FTP", "SIP Server"}, false),
	},
	"alerts_enabled": {
		Type:         schema.TypeInt,
		Description:  "choose 1 to enable alerts, or 0 to disable alerts. Defaults to 1",
		Optional:     true,
		Required:     false,
		Default:      1,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"all_account_group_roles": {
		Type:        schema.TypeList,
		Description: "Role for all account groups",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_id": {
					Type:        schema.TypeInt,
					Description: "role id",
					Optional:    true,
				},
			},
		},
	},
	"api_links": {
		Type:        schema.TypeList,
		Description: "self links to endpoint to pull test metadata, and data links to endpoint for test data",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:        schema.TypeString,
					Description: "href link",
					Computed:    true,
				},
				"rel": {
					Type:        schema.TypeString,
					Description: "rel link",
					Computed:    true,
				},
			},
		},
	},
	"auth_type": {
		Type:         schema.TypeString,
		Description:  "auth type",
		Optional:     true,
		Default:      "NONE",
		ValidateFunc: validation.StringInSlice([]string{"NONE", "BASIC", "NTLM", "KERBEROS"}, false),
	},
	"auth_user": {
		Type:        schema.TypeString,
		Description: "username for authentication with SIP server",
		Required:    true,
	},
	"bandwidth_measurements": {
		Type:        schema.TypeInt,
		Description: "set to 1 to measure bandwidth; defaults to 0. Only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set.",
		Optional:    true,
		Required:    false,
		Default:     1,
	},
	"bgp_measurements": {
		Type:         schema.TypeInt,
		Description:  "choose 1 to enable bgp measurements, 0 to disable; defaults to 1",
		Optional:     true,
		Required:     false,
		Default:      1,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"bgp_monitors": {
		Type:        schema.TypeList,
		Optional:    true,
		Required:    false,
		Description: "array of BGP Monitor objects",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"monitor_id": {
					Type:        schema.TypeInt,
					Description: "monitor id",
					Optional:    true,
				},
			},
		},
	},
	"client_certificate": {
		Type:        schema.TypeString,
		Description: "String representation (containing newline characters) of client certificate, if used",
		Optional:    true,
	},
	"codec": {
		Type:        schema.TypeString,
		Description: "codec label",
		Computed:    true,
	},
	"codec_id": {
		Type:         schema.TypeInt,
		Description:  "codec to use",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 8),
	},
	"content_regex": {
		Type: schema.TypeString,
		Description: "regular Expressions	Verify content using a regular expression. This field does not require escaping",
		Optional: true,
		Default:  "NONE",
	},
	"created_by": {
		Type:        schema.TypeString,
		Description: "created by user",
		Computed:    true,
	},
	"created_date": {
		Type:        schema.TypeString,
		Description: "date of creation",
		Computed:    true,
	},
	"credentials": {
		Type:        schema.TypeMap,
		Description: "Array of credentialID integers.  Get credentialId from /credentials endpoint.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeList,
			Elem: schema.TypeInt,
		},
	},
	"custom_headers": {
		Type: schema.TypeMap,
		Elem: &schema.Schema{
			Type: schema.TypeMap,
			Elem: schema.TypeString,
		},
		Optional: true,
	},
	"default": {
		Type:         schema.TypeInt,
		Description:  "to set the rule as a default, set this value to 1.",
		Optional:     true,
		Default:      0,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"description": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Default:     "",
		Description: "defaults to empty string",
	},
	"desired_status_code": {
		Type: schema.TypeString,
		Description: "A valid HTTP response code	Set to the value you’re interested in retrieving",
		Optional: true,
	},
	"direction": {
		Type: schema.TypeString,
		Description: "[TO_TARGET, FROM_TARGET, BIDIRECTIONAL]	Direction of the test (affects how results are shown)",
		Optional:     false,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"TO_TARGET", "FROM_TARGET", "BIDIRECTIONAL"}, false),
	},
	"direction-alert_rule": {
		Type: schema.TypeString,
		Description: "[TO_TARGET, FROM_TARGET, BIDIRECTIONAL]	Direction of the test (affects how results are shown)",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"TO_TARGET", "FROM_TARGET", "BIDIRECTIONAL"}, false),
	},
	"dns_override": {
		Type:         schema.TypeString,
		Description:  "IP address to use for DNS override",
		Optional:     true,
		ValidateFunc: validation.IsIPAddress,
	},
	"dns_servers": {
		Description: "array of DNS Server objects {“serverName”: “fqdn of server”}",
		Optional:    false,
		Required:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_name": {
					Type:        schema.TypeString,
					Description: "DNS Server name",
					Optional:    true,
				},
			},
		},
	},
	"dns_transport_protocol": {
		Type: schema.TypeString,
		Description: "string	UDP or TCP	transport protocol used for DNS requests; defaults to UDP",
		Optional:     true,
		Required:     false,
		Default:      "UDP",
		ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP"}, false),
	},
	"domain": {
		Type: schema.TypeString,
		Description: "see notes	target record for test, suffixed by record type (ie, www.thousandeyes.com CNAME). If no record type is specified, the test will default to an ANY record.",
		Optional: false,
		Required: true,
	},
	"download_limit": {
		Type:        schema.TypeInt,
		Description: "specify maximum number of bytes to download from the target object",
		Optional:    true,
	},
	"dscp": {
		Type:        schema.TypeString,
		Description: "dscp  label",
		Computed:    true,
	},
	"dscp_id": {
		Type:         schema.TypeInt,
		Description:  "A Differentiated Services Code Point (DSCP) is a value found in an IP packet header which is used to request a level of priority for delivery (Defined in RFC 2474 https://www.ietf.org/rfc/rfc2474.txt). It is one of the Quality of Service management tools used in router configuration to protect real-time and high priority data applications.",
		Required:     false,
		Default:      0,
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{0, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 44, 46, 48, 56}),
	},
	"duration": {
		Type:         schema.TypeInt,
		Description:  "duration of test, in seconds (5 to 30)",
		Optional:     true,
		ValidateFunc: validation.IntBetween(5, 30),
	},
	"email": {
		Type:        schema.TypeString,
		Description: "User email address",
		Required:    true,
	},
	"enabled": {
		Type: schema.TypeInt,
		Description: "0 or 1	choose 1 to enable the test, 0 to disable the test",
		Optional:     true,
		Required:     false,
		Default:      1,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"expression": {
		Type:        schema.TypeString,
		Description: "Alert rule evaluation expression",
		Optional:    true,
	},
	"follow_redirects": {
		Type:        schema.TypeInt,
		Description: "set to 0 to not follow HTTP/301 or HTTP/302 redirect directives. Default is 1",
		Optional:    true,
		Default:     1,
	},
	"ftp_target_time": {
		Type:         schema.TypeInt,
		Description:  "target time for operation completion; specified in milliseconds",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1000, 6000),
	},
	"ftp_time_limit": {
		Type:         schema.TypeInt,
		Description:  "Set the time limit for the test (in seconds). FTP tests default to 10s.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(10, 60),
	},
	"groups": {
		Type:        schema.TypeList,
		Description: "array of label objects",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"group_id": {
					Type:        schema.TypeInt,
					Description: "Unique ID of the label",
					Required:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "Name of the label",
					Optional:    true,
				},
			},
		},
	},
	"headers": {
		Type:        schema.TypeList,
		Description: "array of header strings [\"header: value\", \"header2: value\"]",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"http_interval": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "interval to run http server test on",
	},
	"http_target_time": {
		Type:         schema.TypeInt,
		Description:  "target time for HTTP server completion; specified in milliseconds",
		Optional:     true,
		Default:      1000,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	"http_time_limit": {
		Type:        schema.TypeInt,
		Description: "target time for HTTP server limits; specified in seconds",
		Default:     5,
		Optional:    true,
	},
	"http_version": {
		Type:         schema.TypeInt,
		Description:  "2 for default (prefer HTTP/2), 1 for HTTP/1.1 only",
		Default:      2,
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 2),
	},
	"include_covered_prefixes": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to include queries for subprefixes detected under this prefix",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"include_headers": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to capture response headers for objects loaded by the test.Default is 1.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"interval": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "interval to run test on, in seconds",
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
	},
	"jitter_buffer": {
		Type:         schema.TypeInt,
		Description:  "de-jitter buffer size, in seconds (0 to 150)",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 150),
	},
	"live_share": {
		Type:        schema.TypeInt,
		Description: "indicates 1 for a test shared with your account group, 0 for a normal test (does not apply to DNS+ tests)",
		Computed:    true,
	},
	"login_account_group": {
		Type:        schema.TypeMap,
		Description: "default account group",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aid": {
					Type:        schema.TypeInt,
					Description: "account group id",
					Optional:    true,
				},
			},
		},
	},
	"minimum_sources": {
		Type:        schema.TypeInt,
		Description: "The minimum number of agents or monitors that must meet the specified criteria in order to trigger an alert",
		Optional:    true,
	},
	"minimum_sources_pct": {
		Type:        schema.TypeInt,
		Description: "The minimum percentage of agents or monitors that must meet the specified criteria in order to trigger an alert",
		Optional:    true,
	},
	"modified_by": {
		Type:        schema.TypeString,
		Description: "Last modified by user",
		Computed:    true,
	},
	"modified_date": {
		Type:        schema.TypeString,
		Description: "Last modified by date; shown in UTC",
		Computed:    true,
	},
	"mss": {
		Type: schema.TypeInt,
		Description: "(30..1400)	Maximum Segment Size, in bytes.",
		ValidateFunc: validation.IntBetween(30, 1400),
		Optional:     true,
		Required:     false,
	},
	"mtu_measurements": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to measure MTU sizes on network from agents to the target.",
		Optional:     true,
		Required:     false,
		Default:      1,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "name of the test",
	},
	"network_measurements": {
		Type: schema.TypeInt,
		Description: "integer	0 or 1	choose 1 to enable network measurements, 0 to disable; defaults to 1",
		Default:      1,
		Optional:     true,
		Required:     false,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"notifications": {
		Description: "Map of webhook and third party integrations",
		Optional:    true,
	},
	"notify_on_clear": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to trigger the notification when the alert clears.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"num_path_traces": {
		Type:         schema.TypeInt,
		Description:  "number of path traces. default 3.",
		Default:      3,
		Optional:     true,
		Required:     false,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"options_regex": {
		Type:         schema.TypeString,
		Description:  "regex string. This field does not require escaping.",
		Optional:     true,
		ValidateFunc: validation.StringIsValidRegExp,
	},
	"page_load_target_time": {
		Type:        schema.TypeInt,
		Description: "target time for Page Load completion; specified in seconds (1 to 30); cannot exceed pageLoadTimeLimit value",
		Optional:    true,
	},
	"page_load_time_limit": {
		Type:         schema.TypeInt,
		Description:  "must be larger than httpTimeLimit; defaults to 10 seconds",
		Optional:     true,
		Default:      10,
		ValidateFunc: validation.IntBetween(5, 60),
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "password to be used to authenticate with the destination server",
	},
	"password-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "password to be used to authenticate with the destination server (required for FTP)",
	},
	"path_trace_mode": {
		Type:         schema.TypeString,
		Description:  "choose inSession to perform the path trace within a TCP session; defaults to classic",
		Optional:     true,
		Required:     false,
		Default:      "classic",
		ValidateFunc: validation.StringInSlice([]string{"classic", "inSession"}, false),
	},
	"port": {
		Type:         schema.TypeInt,
		Description:  "target port",
		ValidateFunc: validation.IntBetween(1, 65535),
		Optional:     true,
		Required:     false,
	},
	"post_body": {
		Type:        schema.TypeString,
		Description: "Enter the post body in this field. No escaping is required. If the post body is set to something other than empty, the requestMethod will be set to POST.",
		Optional:    true,
	},
	"prefix": {
		Type:        schema.TypeString,
		Description: "BGP network address prefix",
		Required:    true,
		// a.b.c.d is a network address, with the prefix length defined as e.
		// Prefixes can be any length from 8 to 24
		// Can only use private BGP monitors for a local prefix.
	},
	"probe_mode": {
		Type:         schema.TypeString,
		Description:  "probe mode used by End-to-end Network Test; only valid if protocol is set to TCP; defaults to AUTO",
		Optional:     true,
		Required:     false,
		Default:      "AUTO",
		ValidateFunc: validation.StringInSlice([]string{"AUTO", "SACK", "SYN"}, false),
	},
	"protocol": {
		Type:         schema.TypeString,
		Description:  "protocol used by dependent Network tests (End-to-end, Path Trace, PMTUD); defaults to TCP",
		Optional:     true,
		Required:     false,
		Default:      "TCP",
		ValidateFunc: validation.StringInSlice([]string{"TCP", "ICMP"}, false),
	},
	"protocol-agent_to_agent": {
		Type:         schema.TypeString,
		Description:  "Protocol for agent to agent tests, TCP or UDP.  Defaults to TCP",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
	},
	"protocol-sip": {
		Type:         schema.TypeString,
		Description:  "transport layer for SIP communication: TCP, TLS (TLS over TCP), or UDP. Defaults to TCP",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"TCP", "TLS", "UDP"}, false),
	},
	"recursive_queries": {
		Type:         schema.TypeInt,
		Default:      1,
		ValidateFunc: validation.IntBetween(0, 1),
		Description: "0 or 1	set to 1 to run query with RD (recursion desired) flag enabled",
		Optional: true,
		Required: false,
	},
	"register_enabled": {
		Type:         schema.TypeInt,
		Default:      0,
		Description:  "1 to perform SIP registration on the test target with the SIP REGISTER command, defaults to 0",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"request_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "Set the type of activity for the test: Download, Upload, or List",
	},
	"rounds_violating_mode": {
		Type:         schema.TypeString,
		Description:  "ANY or EXACT.  EXACT requires that the same agent(s) meet the threshold in consecutive rounds; default is ANY",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"ANY", "EXACT"}, false),
	},
	"rounds_violating_required": {
		Type:         schema.TypeInt,
		Description:  "specifies the numerator (X value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10. Must be less than or equal to roundsViolatingOutOf",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"rounds_violating_out_of": {
		Type:         schema.TypeInt,
		Description:  "specifies the divisor (Y value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10.",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"rule_id": {
		Type:        schema.TypeInt,
		Description: "ID of alert rule",
		Computed:    true,
	},
	"rule_name": {
		Type:        schema.TypeString,
		Description: "name of the alert rule",
		Required:    true,
	},
	"saved_event": {
		Type:        schema.TypeInt,
		Description: "indicates 1 for a saved event, 0 for a normal test",
		Computed:    true,
	},
	"server": {
		Type:        schema.TypeString,
		Description: "target host",
		Required:    true,
	},
	"shared_with_accounts": {
		Type:        schema.TypeList,
		Description: "array of DNS Server objects {“serverName”: “fqdn of server”}",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aid": {
					Type:        schema.TypeInt,
					Description: "Account group ID",
					Required:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "Account group name",
					Optional:    true,
				},
			},
		},
	},
	"sip_target_time": {
		Type:         schema.TypeInt,
		Description:  "target time for test completion; specified in milliseconds",
		Optional:     true,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	"sip_time_limit": {
		Type:         schema.TypeInt,
		Description:  "defaults to 5 seconds",
		Optional:     true,
		ValidateFunc: validation.IntBetween(5, 10),
	},
	"source_sip_credentials": {
		Type:     schema.TypeMap,
		Required: true,
	},
	"ssl_version": {
		Type:        schema.TypeString,
		Description: "Reflects the verbose ssl protocol version used by a test",
		Computed:    true,
	},
	"ssl_version_id": {
		Type:         schema.TypeInt,
		Description:  "0 for auto, 3 for SSLv3, 4 for TLS v1.0, 5 for TLS v1.1, 6 for TLS v1.2",
		Optional:     true,
		Default:      0,
		ValidateFunc: validation.IntInSlice([]int{0, 3, 4, 5, 6}),
	},
	"subinterval": {
		Type:         schema.TypeInt,
		Description:  "subinterval for round-robin testing (in seconds), must be less than or equal to interval",
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1200, 1800, 3600}),
	},
	"target_agent_id": {
		Type:     schema.TypeInt,
		Optional: false,
		Required: true,
		Description: "pull from /agents endpoint	Both the 'agents': [] and the targetAgentID cannot be cloud agents. Can be Enterprise Agent -> Cloud, Cloud -> Enterprise Agent, or Enterprise Agent -> Enterprise Agent",
	},
	"target_sip_credentials": {
		Type:     schema.TypeMap,
		Required: true,
	},
	"target_time": {
		Type:         schema.TypeInt,
		Description:  "target time for completion, defaults to 50% of time limit; specified in seconds",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 60),
	},
	"test_id": {
		Type:        schema.TypeInt,
		Description: "Unique ID of test",
		Computed:    true,
	},
	"test_ids": {
		Type:        schema.TypeList,
		Description: "Valid test IDs",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"test_name": {
		Type:        schema.TypeString,
		Description: "Name of Test",
		Required:    true,
	},
	"tests": {
		Type:        schema.TypeList,
		Description: "list of tests",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"test_id": {
					Type:        schema.TypeInt,
					Description: "test id",
					Optional:    true,
				},
			},
		},
	},
	"time_limit": {
		Type:         schema.TypeInt,
		Description:  "time limit for transaction; defaults to 30s",
		Optional:     true,
		Default:      30,
		ValidateFunc: validation.IntBetween(1, 60),
	},
	"throughput_duration": {
		Type:         schema.TypeInt,
		Optional:     true,
		Required:     false,
		Default:      10000,
		Description:  "Defaults to 10000",
		ValidateFunc: validation.IntBetween(5000, 30000),
	},
	"throughput_measurements": {
		Type:         schema.TypeInt,
		ValidateFunc: validation.IntBetween(0, 1),
		Optional:     true,
		Required:     false,
		Default:      0,
		Description: "0 or 1	defaults to 0 (disabled), not allowed when source (or target) of the test is a cloud agent",
	},
	"throughput_rate": {
		Type:         schema.TypeInt,
		Description:  "for UDP only",
		Optional:     true,
		Required:     false,
		Default:      0,
		ValidateFunc: validation.IntBetween(8, 1000),
	},
	"transaction_script": {
		Type:        schema.TypeString,
		Description: "selenium transaction script",
		Required:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "Type of test",
		Computed:    true,
	},
	"type-label": {
		Type:         schema.TypeString,
		Description:  "Type of label (tests, agents, endpoint_tests, or endpoint_agents",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"tests", "agents", "endpoint_tests", "endpoint_agents"}, false),
	},
	"url": {
		Type:        schema.TypeString,
		Description: "target for the test",
		Required:    true,
	},
	"use_active_ftp": {
		Type:         schema.TypeInt,
		Description:  "explicitly set the flag to use active FTP. Tests are set to use passive FTP by default",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"use_explicit_ftps": {
		Type:         schema.TypeInt,
		Description:  "use explicit FTPS (ftp over SSL). By default, tests will autodetect when it is appropriate to use FTPS.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"use_public_bgp": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to automatically add all available Public BGP Monitors",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"use_ntlm": {
		Type:        schema.TypeInt,
		Description: "choose 0 to use Basic Authentication, or omit field.Requires username/password to be set",
		Optional:    true,
	},
	"user": {
		Type:        schema.TypeString,
		Description: "username for SIP registration; should be unique within a ThousandEyes Account Group",
		Optional:    true,
	},
	"user_agent": {
		Type:        schema.TypeString,
		Description: "user-agent string to be provided during the test",
		Optional:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "username to be used to authenticate with the destination server",
	},
	"username-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "username to be used to authenticate with the destination server",
	},
	"verify_certificate": {
		Type:        schema.TypeInt,
		Description: "set to 0 to ignore certificate errors (defaults to 1)",
		Optional:    true,
		Default:     1,
	},
}
