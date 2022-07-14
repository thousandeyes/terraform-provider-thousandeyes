package thousandeyes

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
					Description: "Account group for roles",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"aid": {
								Type:        schema.TypeInt,
								Description: "Account group ID",
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
								Description: "Role ID",
								Optional:    true,
							},
						},
					},
				},
			},
		},
	},
	"agents": {
		Type:        schema.TypeSet,
		Description: "The list of ThousandEyes agents to use.",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:     schema.TypeInt,
					Description: "The unique ID for the ThousandEyes agent.",
					Optional: true,
				},
				"agent_name": {
					Type:     schema.TypeString,
					Description: "The name of the agent.",
					Optional: true,
				},
				"agent_state": {
					Type:     schema.TypeString,
					Description: "Defines whether the agent's status is online, offline, or disabled.",
					Optional: true,
				},
				"agent_type": {
					Type:     schema.TypeString,
					Description: "The type of ThousandEyes agent. Default value is enterprise.",
					Computed: true,
				},
				"country_id": {
					Type:     schema.TypeString,
					Description: "The two-digit ISO country code of the agent.",
					Optional: true,
				},
				"cluster_members": {
					Type:     schema.TypeList,
					Description: "Detailed information about each cluster member, shown as an array. This field is not shown for Enterprise Agents in standalone mode, or for Cloud Agents.",
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"agent_state": {
								Type:     schema.TypeString,
								Description: "Defines whether the agent's status is online, offline, or disabled.",
								Optional: true,
							},
							"ip_addresses": {
								Type:     schema.TypeList,
								Description: "The array of ipAddress entries.",
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"last_seen": {
								Type:     schema.TypeString,
								Description: "The last time the agent connected with ThousandEyes. Uses UTC (yyyy-MM-dd hh:mm:ss).",
								Optional: true,
							},
							"member_id": {
								Type:     schema.TypeInt,
								Description: "The unique ID of the cluster member.",
								Optional: true,
							},
							"name": {
								Type:     schema.TypeString,
								Description: "The name of the cluster member.",
								Optional: true,
							},
							"network": {
								Type:     schema.TypeString,
								Description: "The name of the autonomous system in which the Enterprise Agent is found (Enterprise Agents only).",
								Optional: true,
							},
							"prefix": {
								Type:     schema.TypeString,
								Description: "The network prefix, in CIDR format (Enterprise Agents only).",
								Optional: true,
							},
							"public_ip_addresses": {
								Type:     schema.TypeList,
								Description: "The array of public ipAddress entries.",
								Optional: true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
							"target_for_tests": {
								Type:     schema.TypeString,
								Description: "The target IP address or domain name. Represents the test's destination when the agent is acting as a test target in an agent-to-agent test.",
								Optional: true,
							},
							"utilization": {
								Type:     schema.TypeInt,
								Description: "Shows the overall utilization percentage of a cluster member.",
								Optional: true,
							},
						},
					},
				},
				"created_date": {
					Type:     schema.TypeString,
					Description: "The date the agent was created. Expressed in UTC (yyyy-MM-dd hh:mm:ss).",
					Optional: true,
				},
				"enabled": {
					Type:     schema.TypeBool,
					Description: "Shows whether the agent is enabled or disabled.",
					Optional: true,
				},
				"error_details": {
					Type:     schema.TypeList,
					Description: "If one or more errors present in the agent, the error details are shown for each as an array.",
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"code": {
								Type:     schema.TypeString,
								Description: "[AGENT_VERSION_OUTDATED, APPLIANCE_VERSION_OUTDATED, BROWSERBOT_VERSION_OUTDATED, CLOCK_OFFSET, NAT_TRAVERSAL_ERROR, OS_END_OF_INSTALLATION_SUPPORT, OS_END_OF_SUPPORT, or OS_END_OF_LIFE] The error code.",
								Optional: true,
							},
							"description": {
								Type:     schema.TypeString,
								Description: "A detailed explanation of the error code.",
								Optional: true,
							},
						},
					},
				},
				"groups": {
					Type:     schema.TypeList,
					Description: "An array of label objects.",
					Optional: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"builtin": {
								Type:     schema.TypeInt,
								Description: "Shows whether you are using built-in (1) labels or user-created (2) labels. Built-in labels are read-only.",
								Optional: true,
							},
							"group_id": {
								Type:     schema.TypeInt,
								Description: "The unique ID of the label. This number is negative for built-in labels. Query the /groups/{id} endpoint to see a list of agents/tests with this label.",
								Optional: true,
							},
							"name": {
								Type:     schema.TypeString,
								Description: "The name of the label.",
								Optional: true,
							},
							"type": {
								Type:     schema.TypeString,
								Description: "[tests, agents, endpoint_tests or endpoint_agents] The type of label.",
								Optional: true,
							},
						},
					},
				},
				"hostname": {
					Type:     schema.TypeString,
					Description: "Fully qualified domain name of the agent.",
					Optional: true,
				},
				"ip_addresses": {
					Type:     schema.TypeList,
					Description: "An array of the ipAddress entries.",
					Optional: true,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				"ipv6_policy": {
					Type:     schema.TypeString,
					Description: "[FORCE_IPV4, PREFER_IPV6 or FORCE_IPV6] The IP version policy.",
					Optional: true,
				},
				"keep_browser_cache": {
					Type:     schema.TypeBool,
					Description: "Defines whether the browser cache should be kept. Either 1 for enabled or 0 for disabled.",
					Optional: true,
				},
				"last_seen": {
					Type:     schema.TypeString,
					Description: "The last time the agent connected with ThousandEyes. Shown in UTC (yyyy-MM-dd hh:mm:ss).",
					Optional: true,
				},
				"location": {
					Type:     schema.TypeString,
					Description: "The location of the agent.",
					Optional: true,
				},
				"network": {
					Type:     schema.TypeString,
					Description: "The name of the autonomous system in which the agent is found.",
					Optional: true,
				},
				"prefix": {
					Type:     schema.TypeString,
					Description: "The network prefix, expressed in CIDR format.",
					Optional: true,
				},
				"target_for_tests": {
					Type:     schema.TypeString,
					Description: "The target IP address or domain name representing the test destination when the agent is acting as a test target in an agent-to-agent test.",
					Optional: true,
				},
				"utilization": {
					Type:     schema.TypeInt,
					Description: "Shows the overall utilization percentage.",
					Optional: true,
				},
				"verify_ssl_certificate": {
					Type:     schema.TypeBool,
					Description: "Shows whether the SSL certificate needs to be verified. 1 for enabled and 0 for disabled.",
					Optional: true,
				},
			},
		},
	},
	"agents-label": {
		Type:        schema.TypeList,
		Description: "The list of agents to use.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeInt,
					Description: "The list of unique agent IDs that the label is applied to.",
					Optional:    true,
				},
			},
		},
	},
	"alert_rule_id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of the alert rule.",
		Computed:    true,
	},
	"alert_rules": {
		Description: "Gets the ruleId from the /alert-rules endpoint. If alertsEnabled is set to 'true' and alertRules is not included in a creation/update query, the applicable defaults will be used.",
		Optional:    true,
		Required:    false,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"rule_id": {
					Type:        schema.TypeInt,
					Description: "The unique ID of the alert rule.",
					Optional:    true,
				},
			},
		},
	},
	"alert_type": {
		Description:  "[Page Load, HTTP Server, End-to-End (Server), End-to-End (Agent), DNS+ Domain, DNS+ Server, DNS Server, DNS Trace, DNSSEC, Transactions, Web Transactions, BGP, Path Trace, FTP, SIP Server] The type of alert rule. Acceptable values include the verbose names of supported tests.",
		Type:         schema.TypeString,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"Page Load", "HTTP Server", "End-to-End (Server)", "End-to-End (Agent)", "Voice", "DNS+ Domain", "DNS+ Server", "DNS Server", "DNS Trace", "DNSSEC", "Transactions", "Web Transactions", "BGP", "Path Trace", "FTP", "SIP Server"}, false),
	},
	"alerts_enabled": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to enable alerts, or 'false' to disable alerts. The default value is 'true'.",
		Optional:    true,
		Required:    false,
		Default:     true,
	},
	"all_account_group_roles": {
		Type:        schema.TypeList,
		Description: "The configured role for all account groups.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_id": {
					Type:        schema.TypeInt,
					Description: "The unique ID of the role.",
					Optional:    true,
				},
			},
		},
	},
	"api_links": {
		Type:        schema.TypeList,
		Description: "Self links to the endpoint to pull test metadata, and data links to the endpoint for test data. Read-only, and shows rel and href elements.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"href": {
					Type:        schema.TypeString,
					Description: "The href link.",
					Computed:    true,
				},
				"rel": {
					Type:        schema.TypeString,
					Description: "The rel link.",
					Computed:    true,
				},
			},
		},
	},
	"auth_type": {
		Type:         schema.TypeString,
		Description:  "[NONE, BASIC, NTLM, KERBEROS] The HTTP authentication type. Defaults to NONE.",
		Optional:     true,
		Default:      "NONE",
		ValidateFunc: validation.StringInSlice([]string{"NONE", "BASIC", "NTLM", "KERBEROS"}, false),
	},
	"auth_user": {
		Type:        schema.TypeString,
		Description: "The username for authentication with the SIP server.",
		Required:    true,
	},
	"bandwidth_measurements": {
		Type:        schema.TypeBool,
		Description: "Set to 1 to measure bandwidth. This only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set. Defaults to 'false'.",
		Optional:    true,
		Required:    false,
		Default:     false,
	},
	"bgp_measurements": {
		Type:        schema.TypeBool,
		Description: "Enable BGP measurements. Set to true for enabled, false for disabled.",
		Optional:    true,
		Required:    false,
	},
	"bgp_monitors": {
		Type:        schema.TypeList,
		Optional:    true,
		Required:    false,
		Description: "The array of BGP monitor object IDs. The monitorIDs can be sourced from the /bgp-monitors endpoint.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"ip_address": {
					Type:     schema.TypeString,
					Description: "The IP address of the BGP monitor.",
					Optional: true,
				},
				"monitor_id": {
					Type:        schema.TypeInt,
					Description: "The unique ID of the BGP monitor.",
					Required:    true,
				},
				"monitor_name": {
					Type:     schema.TypeString,
					Description: "The name of the BGP monitor.",
					Optional: true,
				},
				"monitor_type": {
					Type:     schema.TypeString,
					Description: "[Public or Private] Shows the type of BGP monitor.",
					Optional: true,
				},
				"network": {
					Type:     schema.TypeString,
					Description: "The name of the autonomous system in which the BGP monitor is found.",
					Optional: true,
				},
			},
		},
	},
	"builtin": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' for built-in labels, or to 'false' for user-created labels. Built-in labels are read-only.",
		Computed:    true,
	},
	"client_certificate": {
		Type:        schema.TypeString,
		Description: "String representation (containing newline characters) of the client certificate, if used.",
		Optional:    true,
	},
	"codec": {
		Type:        schema.TypeString,
		Description: "The label of the codec.",
		Computed:    true,
	},
	"codec_id": {
		Type:         schema.TypeInt,
		Description:  "The unique ID of the codec to use.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 8),
	},
	"content_regex": {
		Type:        schema.TypeString,
		Description: "Verify content using a regular expression. This field does not require escaping.",
		Optional:    true,
		Required:    false,
		Default:     ".*",
	},
	"created_by": {
		Type:        schema.TypeString,
		Description: "Created by user.",
		Computed:    true,
	},
	"created_date": {
		Type:        schema.TypeString,
		Description: "The date of creation.",
		Computed:    true,
	},
	"credentials": {
		Type:        schema.TypeList,
		Description: "The array of credentialID integers. You can get the credentialId from the /credentials endpoint.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"custom_headers": {
		Type:     schema.TypeMap,
		Description: "The custom headers.",
		Optional: true,
	},
	"default": {
		Type:        schema.TypeBool,
		Description: "Alert rules allow up to 1 alert rule to be selected as a default for each type. By marking an alert rule as default, the rule will be automatically included in subsequently created tests that test a metric used in the alert rule.",
		Optional:    true,
		Default:     false,
	},
	"description": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Default:     "",
		Description: "A description of the alert rule. Defaults to an empty string.",
	},
	"desired_status_code": {
		Type: schema.TypeString,
		Description: "The valid HTTP response code you’re interested in retrieving.",
		Optional: true,
	},
	"direction": {
		Type: schema.TypeString,
		Description: "[TO_TARGET, FROM_TARGET, BIDIRECTIONAL] The direction of the test (affects how results are shown).",
		Optional:     false,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"TO_TARGET", "FROM_TARGET", "BIDIRECTIONAL"}, false),
	},
	"direction-alert_rule": {
		Type: schema.TypeString,
		Description: "[TO_TARGET, FROM_TARGET, BIDIRECTIONAL] The direction of the test (affects how results are shown).",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"TO_TARGET", "FROM_TARGET", "BIDIRECTIONAL"}, false),
	},
	"dns_override": {
		Type:         schema.TypeString,
		Description:  "The IP address to use for DNS override.",
		Optional:     true,
		ValidateFunc: validation.IsIPAddress,
	},
	"dns_servers": {
		Description: "The array of DNS Server objects (“serverName”: “fqdn of server”).",
		Optional:    false,
		Required:    true,
		Type:        schema.TypeList,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_name": {
					Type:        schema.TypeString,
					Description: "The DNS server name.",
					Optional:    true,
				},
				"server_id": {
					Type:     schema.TypeInt,
					Description: "The unique ID of the DNS server.",
					Optional: true,
				},
			},
		},
	},
	"dns_transport_protocol": {
		Type: schema.TypeString,
		Description: "[UDP or TCP] The DNS transport protocol used for DNS requests. Defaults to UDP.",
		Optional:     true,
		Required:     false,
		Default:      "UDP",
		ValidateFunc: validation.StringInSlice([]string{"UDP", "TCP"}, false),
	},
	"domain": {
		Type: schema.TypeString,
		Description: "See notes	target record for test, suffixed by record type (ie, www.thousandeyes.com CNAME). If no record type is specified, the test will default to an ANY record.",
		Optional: false,
		Required: true,
	},
	"download_limit": {
		Type:        schema.TypeInt,
		Description: "Specify the maximum number of bytes to download from the target object.",
		Optional:    true,
	},
	"dscp": {
		Type:        schema.TypeString,
		Description: "The Differentiated Services Code Point (DSCP) label.",
		Computed:    true,
	},
	"dscp_id": {
		Type:         schema.TypeInt,
		Description:  "The DSCP ID.",
		Required:     false,
		Default:      0,
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{0, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 44, 46, 48, 56}),
	},
	"duration": {
		Type:         schema.TypeInt,
		Description:  "The duration of the test, in seconds (5 to 30).",
		Optional:     true,
		ValidateFunc: validation.IntBetween(5, 30),
	},
	"email": {
		Type:        schema.TypeString,
		Description: "The user's email address.",
		Required:    true,
	},
	"enabled": {
		Type:        schema.TypeBool,
		Description: "Enables or disables the test.",
		Optional:    true,
		Required:    false,
		Default:     true,
	},
	"expression": {
		Type:         schema.TypeString,
		Description:  "The alert rule evaluation expression.",
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	},
	"follow_redirects": {
		Type:        schema.TypeBool,
		Description: "Follow HTTP/301 or HTTP/302 redirect directives. Defaults to 'true'.",
		Optional:    true,
		Default:     true,
	},
	"ftp_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for operation completion. Specified in milliseconds.",
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
		Description: "The array of label objects.",
		Optional:    true,
		Elem: &schema.Resource{
			// Schema definition here is to support group objects returned from
			// reads of test resources.
			Schema: map[string]*schema.Schema{
				"agents": {
					// See `agents-label` rather than `agents`
					Type:        schema.TypeList,
					Description: "Define the ThousandEyes agents to use.",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"agent_id": {
								Type:        schema.TypeInt,
								Description: "The unique ThousandEyes agent ID.",
								Optional:    true,
							},
						},
					},
				},
				"builtin": {
					Type:        schema.TypeBool,
					Description: "Shows whether you are using built-in (true) labels or user-created (false) labels. Built-in labels are read-only.",
					Computed:    true,
				},
				"group_id": {
					Type:        schema.TypeInt,
					Description: "The unique ID of the label",
					Required:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the label.",
					Optional:    true,
				},
				"tests": {
					Type:        schema.TypeList,
					Description: "The list of tests.",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"test_id": {
								Type:        schema.TypeInt,
								Description: "The unique ID of the test.",
								Optional:    true,
							},
						},
					},
				},
				"type": {
					// See `type-label` rather than `type`
					Type:        schema.TypeString,
					Description: "[tests, agents, endpoint_tests, or endpoint_agents] The type of label.",
					Optional:    true,
				},
			},
		},
	},
	"group_id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of the label. For built-in labels, this number is a negative.",
		Computed:    true,
	},
	"headers": {
		Type:        schema.TypeList,
		Description: "[\"header: value\", \"header2: value\"] The array of header strings.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	"http_interval": {
		Type:        schema.TypeInt,
		Required:    true,
		Description: "The interval to run the HTTP server test on.",
	},
	"http_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for HTTP server completion, specified in milliseconds.",
		Optional:     true,
		Default:      1000,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	"http_time_limit": {
		Type:        schema.TypeInt,
		Description: "The target time for HTTP server limits, specified in seconds.",
		Default:     5,
		Optional:    true,
	},
	"http_version": {
		Type:         schema.TypeInt,
		Description:  "Set to 2 for the default HTTP version (prefer HTTP/2), or 1 for HTTP/1.1 only.",
		Default:      2,
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 2),
	},
	"include_covered_prefixes": {
		Type:        schema.TypeBool,
		Description: "Include queries for subprefixes detected under this prefix.",
		Optional:    true,
	},
	"include_headers": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to capture response headers for objects loaded by the test.",
		Optional:    true,
	},
	"interval": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "The interval to run the test on, in seconds.",
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
	},
	"jitter_buffer": {
		Type:         schema.TypeInt,
		Description:  "The de-jitter buffer size, in seconds (0 to 150).",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 150),
	},
	"live_share": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' for a test shared with your account group, or to 'false' for a normal test.",
		Computed:    true,
	},
	"login_account_group": {
		Type:        schema.TypeMap,
		Description: "The default account group.",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aid": {
					Type:        schema.TypeInt,
					Description: "The unique ID of the account group.",
					Optional:    true,
				},
			},
		},
	},
	"minimum_sources": {
		Type:         schema.TypeInt,
		Description:  "The minimum number of agents or monitors that must meet the specified criteria in order to trigger an alert. This option is mutually exclusive with 'minimum_sources_pct'.",
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	"minimum_sources_pct": {
		Type:         schema.TypeInt,
		Description:  "The minimum percentage of agents or monitors that must meet the specified criteria in order to trigger an alert. This option is mutually exclusive with 'minimum_sources'.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 100),
	},
	"modified_by": {
		Type:        schema.TypeString,
		Description: "Last modified by this user.",
		Computed:    true,
	},
	"modified_date": {
		Type:        schema.TypeString,
		Description: "The date the test was last modified. Shown in UTC.",
		Computed:    true,
	},
	"mss": {
		Type: schema.TypeInt,
		Description: "The maximum segment size, in bytes. Value can be from 30 to 1400.",
		ValidateFunc: validation.IntBetween(30, 1400),
		Optional:     true,
		Required:     false,
	},
	"mtu_measurements": {
		Type:        schema.TypeBool,
		Description: "Measure MTU sizes on the network from agents to the target.",
		Optional:    true,
		Required:    false,
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// If this field isn't set, it will change when
			// network_measurements change. If we're not explicitly
			// setting this field, then ignore the diff.
			return !d.HasChange(k)
		},
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The name of the test.",
	},
	"network_measurements": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to enable network measurements.",
		Default:     true,
		Optional:    true,
		Required:    false,
	},
	"notifications": {
		Type:        schema.TypeSet,
		Description: "The list of notifications for the alert rule.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"email": {
					Type:        schema.TypeSet,
					Description: "The email notification.",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"message": {
								Type:        schema.TypeString,
								Description: "The contents of the email, as a string.",
								Optional:    true,
							},
							"recipient": {
								Type:        schema.TypeList,
								Description: "The email addresses to send the notification to.",
								Optional:    true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
			},
		},
	},
	"notify_on_clear": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to trigger the notification when the alert clears.",
		Optional:    true,
		Default:     true,
	},
	"num_path_traces": {
		Type:         schema.TypeInt,
		Description:  "The number of path traces.",
		Optional:     true,
		Required:     false,
		ValidateFunc: validation.IntBetween(1, 10),
		DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
			// If this field isn't set, it will change when
			// network_measurements change. If we're not explicitly
			// setting this field, then ignore the diff.
			return !d.HasChange(k)
		},
	},
	"options_regex": {
		Type:         schema.TypeString,
		Description:  "A regex string of options. This field does not require escaping.",
		Optional:     true,
		ValidateFunc: validation.StringIsValidRegExp,
	},
	"page_load_target_time": {
		Type:        schema.TypeInt,
		Description: "The target time for page load completion, specified in seconds (1 to 30). The value cannot exceed the pageLoadTimeLimit value.",
		Optional:    true,
	},
	"page_load_time_limit": {
		Type:         schema.TypeInt,
		Description:  "The page load time limit. This value must be larger than httpTimeLimit, and defaults to 10 seconds.",
		Optional:     true,
		Default:      10,
		ValidateFunc: validation.IntBetween(5, 60),
	},
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The password to be used to authenticate with the destination server.",
	},
	"password-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The password to be used to authenticate with the destination server (required for FTP).",
	},
	"path_trace_mode": {
		Type:         schema.TypeString,
		Description:  "[classic or inSession] Choose 'inSession' to perform the path trace within a TCP session. Default value is 'classic'.",
		Optional:     true,
		Required:     false,
		Default:      "classic",
		ValidateFunc: validation.StringInSlice([]string{"classic", "inSession"}, false),
	},
	"port": {
		Type:         schema.TypeInt,
		Description:  "The target port.",
		ValidateFunc: validation.IntBetween(1, 65535),
		Optional:     true,
		Required:     false,
	},
	"post_body": {
		Type:        schema.TypeString,
		Description: "The POST body content. No escaping is required. If the post body is set to something other than empty, the requestMethod will be set to POST.",
		Optional:    true,
	},
	"prefix": {
		Type:        schema.TypeString,
		Description: "The BGP network address prefix.",
		Required:    true,
		ForceNew:    true,
		// a.b.c.d is a network address, with the prefix length defined as e.
		// Prefixes can be any length from 8 to 24
		// Can only use private BGP monitors for a local prefix.
	},
	"probe_mode": {
		Type:         schema.TypeString,
		Description:  "[AUTO, SACk, or SYN] The probe mode used by end-to-end network tests. This is only valid if the protocol is set to TCP. The default value is AUTO.",
		Optional:     true,
		Required:     false,
		Default:      "AUTO",
		ValidateFunc: validation.StringInSlice([]string{"AUTO", "SACK", "SYN"}, false),
	},
	"protocol": {
		Type:         schema.TypeString,
		Description:  "The protocol used by dependent network tests (end-to-end, path trace, PMTUD). Default value is TCP.",
		Optional:     true,
		Required:     false,
		Default:      "TCP",
		ValidateFunc: validation.StringInSlice([]string{"TCP", "ICMP"}, false),
	},
	"protocol-agent_to_agent": {
		Type:         schema.TypeString,
		Description:  "[TCP or UDP] The protocol for agent to agent tests. Defaults to TCP.",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"TCP", "UDP"}, false),
	},
	"recursive_queries": {
		Type:        schema.TypeBool,
		Default:     true,
		Description: "Defines whether to run the query with the recursion desired (RD) flag enabled.",
		Optional:    true,
		Required:    false,
	},
	"register_enabled": {
		Type:         schema.TypeBool,
		Default:      0,
		Description:  "Configure whether to perform SIP registration on the test target with the SIP REGISTER command. Default value is 'false'.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"request_type": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "[Download, Upload, or List] Sets the type of activity for the test.",
	},
	"rounds_violating_mode": {
		Type:         schema.TypeString,
		Description:  "[ANY or EXACT] Defines whether the same agent(s) must meet the EXACT same threshold in consecutive rounds or not. The default value is ANY.",
		Default:      "ANY",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"ANY", "EXACT"}, false),
	},
	"rounds_violating_required": {
		Type:         schema.TypeInt,
		Description:  "Specifies the numerator (X value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10. Must be less than or equal to 'roundsViolatingOutOf'.",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"rounds_violating_out_of": {
		Type:         schema.TypeInt,
		Description:  "Specifies the divisor (Y value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10.",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	"rule_id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of the alert rule.",
		Computed:    true,
	},
	"rule_name": {
		Type:        schema.TypeString,
		Description: "The name of the alert rule.",
		Required:    true,
	},
	"saved_event": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' for a saved event, or to 'false' for a normal test.",
		Computed:    true,
	},
	"server": {
		Type:        schema.TypeString,
		Description: "The target host.",
		Required:    true,
	},
	"shared_with_accounts": {
		Type:        schema.TypeList,
		Description: "[“serverName”: “fqdn of server”] The array of DNS Server objects.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"aid": {
					Type:        schema.TypeInt,
					Description: "The account group ID.",
					Required:    true,
				},
				"name": {
					Type:        schema.TypeString,
					Description: "The name of account.",
					Optional:    true,
				},
			},
		},
	},
	"sip_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for test completion, specified in milliseconds.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	"sip_time_limit": {
		Type:         schema.TypeInt,
		Description:  "The test time limit. Can be between 5 and 10 seconds, and defaults to 5 seconds.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(5, 10),
	},
	"source_sip_credentials": {
		Type:     schema.TypeMap,
		Description: "The SIP credentials.",
		Required: true,
	},
	"ssl_version": {
		Type:        schema.TypeString,
		Description: "Reflects the verbose ssl protocol version used by a test.",
		Computed:    true,
	},
	"ssl_version_id": {
		Type:         schema.TypeInt,
		Description:  "Defines the SSL version. 0 for auto, 3 for SSLv3, 4 for TLS v1.0, 5 for TLS v1.1, 6 for TLS v1.2.",
		Optional:     true,
		Default:      0,
		ValidateFunc: validation.IntInSlice([]int{0, 3, 4, 5, 6}),
	},
	"subinterval": {
		Type:         schema.TypeInt,
		Description:  "The subinterval for round-robin testing (in seconds). The value must be less than or equal to 'interval'.",
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1200, 1800, 3600}),
	},
	"target_agent_id": {
		Type:     schema.TypeInt,
		Optional: false,
		Required: true,
		Description: "The target agent's unique ID. Pulled from the /agents endpoint. Both the 'agents': [] and the targetAgentID cannot be Cloud Agents. Can be Enterprise Agent -> Cloud, Cloud -> Enterprise Agent, or Enterprise Agent -> Enterprise Agent.",
	},
	"target_sip_credentials": {
		Type:     schema.TypeList,
		Description: "The Target SIP server credentials.",
		MaxItems: 1,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auth_user": {
					Type:        schema.TypeString,
					Description: "The username for authentication with the SIP server.",
					Required:    true,
				},
				"password": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The password to be used to authenticate with the destination server.",
				},
				"port": {
					Type:         schema.TypeInt,
					Description:  "The target port.",
					ValidateFunc: validation.IntBetween(1, 65535),
					Optional:     true,
					Required:     false,
				},
				"protocol": {
					Type:         schema.TypeString,
					Description:  "[TCP, TLS, or UDP] The transport layer for SIP communication. Can be either TCP, TLS (TLS over TCP), or UDP, and defaults to TCP.",
					Required:     true,
					ValidateFunc: validation.StringInSlice([]string{"TCP", "TLS", "UDP"}, false),
				},
				"sip_proxy": {
					Type:        schema.TypeString,
					Description: "The SIP proxy. This is distinct from the SIP server, and is specified as a domain name or IP address.",
					Optional:    true,
				},
				"sip_registrar": {
					Type:        schema.TypeString,
					Description: "The SIP server to be tested, specified by domain name or IP address.",
					Required:    true,
				},
				"user": {
					Type:        schema.TypeString,
					Description: "The username for SIP registration. This should be unique within a ThousandEyes account group.",
					Optional:    true,
				},
			},
		},
	},
	"target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for completion. The default value is 50 percent of the time limit, specified in seconds.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 60),
	},
	"test_id": {
		Type:        schema.TypeInt,
		Description: "The unique ID of the test.",
		Computed:    true,
	},
	"test_ids": {
		Type:        schema.TypeList,
		Description: "The valid test IDs.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeInt,
		},
	},
	"test_name": {
		Type:        schema.TypeString,
		Description: "The name of the test.",
		Required:    true,
	},
	"tests": {
		Type:        schema.TypeList,
		Description: "The list of included tests.",
		Optional:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"test_id": {
					Type:        schema.TypeInt,
					Description: "The list of unique test IDs.",
					Optional:    true,
				},
			},
		},
	},
	"time_limit": {
		Type:         schema.TypeInt,
		Description:  "The time limit for the transaction. The default value is 30s.",
		Optional:     true,
		Default:      30,
		ValidateFunc: validation.IntBetween(1, 60),
	},
	"throughput_duration": {
		Type:         schema.TypeInt,
		Optional:     true,
		Required:     false,
		Default:      10000,
		Description:  "The throughput duration in milliseconds. The default value is 10000.",
		ValidateFunc: validation.IntBetween(5000, 30000),
	},
	"throughput_measurements": {
		Type:        schema.TypeBool,
		Optional:    true,
		Required:    false,
		Default:     false,
		Description: "Enables or disables throughput measurements. This is not allowed when the source or target of the test is a Cloud Agent. Defaults to disabled.",
	},
	"throughput_rate": {
		Type:         schema.TypeInt,
		Description:  "Defines the throughput rate. Fo UDP tests only.",
		Optional:     true,
		Required:     false,
		Default:      0,
		ValidateFunc: validation.IntBetween(8, 1000),
	},
	"transaction_script": {
		Type:        schema.TypeString,
		Description: "The full selenium transaction script.",
		Required:    true,
	},
	"type": {
		Type:        schema.TypeString,
		Description: "The type of test.",
		Computed:    true,
	},
	"type-label": {
		Type:         schema.TypeString,
		Description:  "[tests, agents, endpoint_tests, or endpoint_agents] The type of label.",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"tests", "agents", "endpoint_tests", "endpoint_agents"}, false),
	},
	"url": {
		Type:        schema.TypeString,
		Description: "The target URL for the test.",
		Required:    true,
	},
	"use_active_ftp": {
		Type:         schema.TypeInt,
		Description:  "Enables active FTP. If not set, tests default to use passive FTP.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"use_explicit_ftps": {
		Type:         schema.TypeInt,
		Description:  "Enables explicit FTPS (FTP over SSL). By default, tests will autodetect when it is appropriate to use FTPS.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"use_public_bgp": {
		Type:        schema.TypeBool,
		Description: "Enable to automatically add all available Public BGP Monitors to the test.",
		Optional:    true,
		Default:     true,
	},
	"use_ntlm": {
		Type:        schema.TypeBool,
		Description: "Enable to use basic authentication. Only include this field if you are using authentication. Requires the username and password to be set if enabled.",
		Optional:    true,
	},
	"user": {
		Type:        schema.TypeString,
		Description: "The username for SIP registration. This username should be unique within the ThousandEyes account group.",
		Optional:    true,
	},
	"user_agent": {
		Type:        schema.TypeString,
		Description: "The user-agent string to be provided during the test.",
		Optional:    true,
	},
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The username to be used to authenticate with the destination server.",
	},
	"username-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The username to be used to authenticate with the destination server.",
	},
	"verify_certificate": {
		Type:        schema.TypeBool,
		Description: "Set whether to ignore certificate errors. Set to 'false' to ignore certificate errors. The default value is 'true'.",
		Optional:    true,
		Default:     true,
	},
}
