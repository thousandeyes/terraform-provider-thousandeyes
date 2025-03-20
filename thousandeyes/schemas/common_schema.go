package schemas

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

// Structs used for mapping
var (
	_ = tests.AgentToAgentTestRequest{}
	_ = tests.AgentToAgentTestResponse{}

	_ = tests.AgentToServerTestRequest{}
	_ = tests.AgentToServerTestResponse{}

	_ = tests.BgpTestRequest{}
	_ = tests.UpdateBgpTestRequest{}
	_ = tests.BgpTestResponse{}

	_ = tests.DnsServerTestRequest{}
	_ = tests.DnsServerTestResponse{}

	_ = tests.DnsTraceTestRequest{}
	_ = tests.DnsTraceTestResponse{}

	_ = tests.DnsSecTestRequest{}
	_ = tests.DnsSecTestResponse{}

	_ = tests.FtpServerTestRequest{}
	_ = tests.FtpServerTestResponse{}

	_ = tests.HttpServerTestRequest{}
	_ = tests.HttpServerTestResponse{}

	_ = tests.PageLoadTestRequest{}
	_ = tests.PageLoadTestResponse{}

	_ = tests.SipServerTestRequest{}
	_ = tests.SipServerTestResponse{}

	_ = tests.VoiceTestRequest{}
	_ = tests.VoiceTestResponse{}

	_ = tests.WebTransactionTestRequest{}
	_ = tests.WebTransactionTestResponse{}
)

var CommonSchema = map[string]*schema.Schema{
	// COMMON

	// alertsEnabled
	"alerts_enabled": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to enable alerts, or 'false' to disable alerts. The default value is 'true'.",
		Optional:    true,
		Required:    false,
		Default:     true,
	},
	// enabled
	"enabled": {
		Type:        schema.TypeBool,
		Description: "Enables or disables the test.",
		Optional:    true,
		Required:    false,
		Default:     true,
	},
	// alertRules
	"alert_rules": {
		Description: "List of alert rules IDs to apply to the test (get `ruleId` from `/alerts/rules` endpoint. If `alertsEnabled` is set to `true` and `alertRules` is not included on test creation or update, applicable user default alert rules will be used)",
		Optional:    true,
		Required:    false,
		Type:        schema.TypeList,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// createdBy
	"created_by": {
		Type:        schema.TypeString,
		Description: "Created by user.",
		Computed:    true,
	},
	// createdDate
	"created_date": {
		Type:        schema.TypeString,
		Description: "The date of creation.",
		Computed:    true,
	},
	// description
	"description": {
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		Default:     "",
		Description: "A description of the alert rule. Defaults to an empty string.",
	},
	// liveShare
	"live_share": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' for a test shared with your account group, or to 'false' for a normal test.",
		Computed:    true,
	},
	// modifiedBy
	"modified_by": {
		Type:        schema.TypeString,
		Description: "Last modified by this user.",
		Computed:    true,
	},
	// modifiedDate
	"modified_date": {
		Type:        schema.TypeString,
		Description: "The date the test was last modified. Shown in UTC.",
		Computed:    true,
	},
	// savedEvent
	"saved_event": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' for a saved event, or to 'false' for a normal test.",
		Computed:    true,
	},
	// testId
	"test_id": {
		Type:        schema.TypeString,
		Description: "The unique ID of the test.",
		Computed:    true,
	},
	// testName
	"test_name": {
		Type:        schema.TypeString,
		Description: "The name of the test.",
		Required:    true,
	},
	// type
	"type": {
		Type:        schema.TypeString,
		Description: "The type of test.",
		Computed:    true,
	},
	// _links
	"_links": { // New
		Type:        schema.TypeSet,
		Description: "Self links to the endpoint to pull test metadata, and data links to the endpoint for test data. Read-only, and shows rel and href elements.",
		Computed:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"self": link,
				"test_results": {
					Type:        schema.TypeList,
					Description: "Reference to the test results.",
					Elem:        link,
				},
			},
		},
	},
	// labels
	"labels": {
		Type:        schema.TypeList,
		Description: "[\"1\", \"2\", \"uuid\"] The array of label or tag ids.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	// sharedWithAccounts
	"shared_with_accounts": {
		Type:        schema.TypeList,
		Description: "List of accounts",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// usePublicBgp
	"use_public_bgp": {
		Type:        schema.TypeBool,
		Description: "Enable to automatically add all available Public BGP Monitors to the test.",
		Optional:    true,
		Default:     false,
	},
	// monitors (ex. bgp_monitors)
	"monitors": {
		Type:        schema.TypeList,
		Description: " Contains list of BGP monitor IDs (get `monitorId` from `/monitors` endpoint)",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// agents
	"agents": {
		Type:        schema.TypeSet,
		Description: "The list of ThousandEyes agents to use.",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeString,
					Description: "The unique ID for the ThousandEyes agent.",
					Required:    true,
				},
			},
		},
	},
	// interval
	"interval": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "The interval to run the test on, in seconds.",
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
	},
	// fixedPacketRate
	"fixed_packet_rate": { // New
		Type:        schema.TypeInt,
		Optional:    true,
		Required:    false,
		Description: "Sets packets rate sent to measure the network in packets per second.",
	},
	// numPathTraces
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
	// pathTraceMode
	"path_trace_mode": {
		Type:         schema.TypeString,
		Description:  "[classic or in-session] Choose 'inSession' to perform the path trace within a TCP session. Default value is 'classic'.",
		Optional:     true,
		Required:     false,
		Default:      "classic",
		ValidateFunc: validation.StringInSlice([]string{"classic", "in-session"}, false),
	},
	// dscp
	"dscp": {
		Type:        schema.TypeString,
		Description: "The Differentiated Services Code Point (DSCP) label.",
		Computed:    true,
	},
	// dscpId
	"dscp_id": {
		Type:        schema.TypeString,
		Description: "The DSCP ID.",
		Required:    false,
		Optional:    true,
	},
	// randomizedStartTime
	"randomized_start_time": {
		Type:        schema.TypeBool,
		Description: "Indicates whether agents should randomize the start time in each test round.",
		Optional:    true,
		Required:    false,
		Default:     false,
	},
	// bgpMeasurements
	"bgp_measurements": {
		Type:        schema.TypeBool,
		Description: "Enable BGP measurements. Set to true for enabled, false for disabled.",
		Optional:    true,
		Required:    false,
	},

	// AGENT TO AGENT

	// direction
	"direction": {
		Type:         schema.TypeString,
		Description:  "[to-target, from-target, bidirectional] The direction of the test (affects how results are shown).",
		Optional:     false,
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"to-target", "from-target", "bidirectional"}, false),
	},
	// mss
	"mss": {
		Type:         schema.TypeInt,
		Description:  "The maximum segment size, in bytes. Value can be from 30 to 1400.",
		ValidateFunc: validation.IntBetween(30, 1400),
		Optional:     true,
		Required:     false,
	},
	// targetAgentId
	"target_agent_id": {
		Type:        schema.TypeString,
		Optional:    false,
		Required:    true,
		Description: "The target agent's unique ID. Pulled from the /agents endpoint. Both the 'agents': [] and the targetAgentID cannot be Cloud Agents. Can be Enterprise Agent -> Cloud, Cloud -> Enterprise Agent, or Enterprise Agent -> Enterprise Agent.",
	},
	// throughputMeasurements
	"throughput_measurements": {
		Type:        schema.TypeBool,
		Optional:    true,
		Required:    false,
		Default:     false,
		Description: "Enables or disables throughput measurements. This is not allowed when the source or target of the test is a Cloud Agent. Defaults to disabled.",
	},
	// throughputDuration
	"throughput_duration": {
		Type:         schema.TypeInt,
		Optional:     true,
		Required:     false,
		Default:      10000,
		Description:  "The throughput duration in milliseconds. The default value is 10000.",
		ValidateFunc: validation.IntBetween(5000, 30000),
	},
	// throughputRate
	"throughput_rate": {
		Type:         schema.TypeInt,
		Description:  "Defines the throughput rate. Fo UDP tests only.",
		Optional:     true,
		Required:     false,
		ValidateFunc: validation.IntBetween(8, 1000),
	},
	// protocol
	"protocol-a2a": {
		Type:         schema.TypeString,
		Description:  "[tcp or udp] The protocol for agent to agent tests. Defaults to 'tcp'.",
		Required:     false,
		Optional:     true,
		Default:      "tcp",
		ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
	},

	// AGENT TO SERVER

	// bandwidthMeasurements
	"bandwidth_measurements": {
		Type:        schema.TypeBool,
		Description: "Set to 1 to measure bandwidth. This only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set. Defaults to 'false'.",
		Optional:    true,
		Required:    false,
		Default:     false,
	},
	// continuousMode
	"continuous_mode": { // New
		Type:        schema.TypeBool,
		Description: "To enable continuous monitoring, set this parameter to `true` to.  When continuous monitoring is enabled, the following actions occur: * `fixedPacketRate` is enforced * `bandwidthMeasurements` are disabled * If the `protocol` is set to `tcp`, `probeMode` is set to `syn`.",
		Optional:    true,
		Required:    false,
		Default:     false,
	},
	// mtuMeasurements
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
	// probeMode
	"probe_mode": {
		Type:         schema.TypeString,
		Description:  "[auto, sack, or syn] The probe mode used by end-to-end network tests. This is only valid if the protocol is set to TCP. The default value is AUTO.",
		Optional:     true,
		Required:     false,
		Default:      "auto",
		ValidateFunc: validation.StringInSlice([]string{"auto", "sack", "syn"}, false),
	},
	// server
	"server": {
		Type:        schema.TypeString,
		Description: "The target host.",
		Required:    true,
	},
	// ipv6Policy
	"ipv6_policy": { // New
		Type:        schema.TypeString,
		Description: "[force-ipv4, prefer-ipv6, force-ipv6, or use-agent-policy]", // TO DO describe
		Optional:    true,
		Required:    false,
		ValidateFunc: validation.StringInSlice([]string{
			"force-ipv4",
			"prefer-ipv6",
			"force-ipv6",
			"use-agent-policy",
		}, false),
	},
	// pingPayloadSize
	"ping_payload_size": {
		Type:         schema.TypeInt,
		Description:  "Payload size (not total packet size) for the end-to-end metric's probes, ping payload size allows values from 0 to 1400 bytes. When set to null, payload sizes are 0 bytes for ICMP-based tests and 1 byte for TCP-based tests.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 1400),
	}, // New
	// networkMeasurements
	"network_measurements": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to enable network measurements.",
		Default:     true,
		Optional:    true,
		Required:    false,
	},

	// BGP

	// includeCoveredPrefixes
	"include_covered_prefixes": {
		Type:        schema.TypeBool,
		Description: "Include queries for subprefixes detected under this prefix.",
		Optional:    true,
	},
	// prefix
	"prefix": {
		Type:        schema.TypeString,
		Description: "The BGP network address prefix.",
		Required:    true,
		ForceNew:    true,
		// a.b.c.d is a network address, with the prefix length defined as e.
		// Prefixes can be any length from 8 to 24
		// Can only use private BGP monitors for a local prefix.
	},

	// DNS

	//dnsServers
	"dns_servers": {
		Description: "The array of DNS Server objects (“serverName”: “fqdn of server”).",
		Optional:    false,
		Required:    true,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"server_name": {
					Type:        schema.TypeString,
					Description: "The DNS server name.",
					Optional:    true,
				},
				"server_id": {
					Type:        schema.TypeString,
					Description: "The unique ID of the DNS server.",
					Optional:    true,
				},
			},
		},
	},
	// dnsTransportProtocol
	"dns_transport_protocol": {
		Type:         schema.TypeString,
		Description:  "[udp or tcp] The DNS transport protocol used for DNS requests. Defaults to 'udp'.",
		Optional:     true,
		Required:     false,
		Default:      "udp",
		ValidateFunc: validation.StringInSlice([]string{"udp", "tcp"}, false),
	},
	// domain
	"domain": {
		Type:        schema.TypeString,
		Description: "See notes	target record for test, suffixed by record type (ie, www.thousandeyes.com CNAME). If no record type is specified, the test will default to an ANY record.",
		Optional:    false,
		Required:    true,
		ValidateFunc: validation.StringMatch(
			regexp.MustCompile(`^.* (A|ANY|NS|CNAME|MX|SOA|AAAA|PTR|TXT|NULL|DS|RRSIG|DNSKEY|NSEC)$`),
			"must suffix with record type; check ThousandEyes Developer Reference for more information",
		),
	},
	// protocol
	"protocol": {
		Type:         schema.TypeString,
		Description:  "The protocol used by dependent network tests (end-to-end, path trace, PMTUD). Default value is tcp.",
		Optional:     true,
		Required:     false,
		Default:      "tcp",
		ValidateFunc: validation.StringInSlice([]string{"tcp", "icmp", "udp"}, false),
	},
	// recursiveQueries
	"recursive_queries": {
		Type:        schema.TypeBool,
		Default:     true,
		Description: "Defines whether to run the query with the recursion desired (RD) flag enabled.",
		Optional:    true,
		Required:    false,
	},
	// dnsQueryClass
	"dns_query_class": { // New
		Type:         schema.TypeString,
		Description:  "Domain class used by this test. 'in' stands for Internet, while 'ch' stands for Chaos.",
		Optional:     true,
		Required:     false,
		ValidateFunc: validation.StringInSlice([]string{"in", "ch"}, false),
	},

	// FTP SERVER

	// downloadLimit
	"download_limit": {
		Type:        schema.TypeInt,
		Description: "Specify the maximum number of bytes to download from the target object.",
		Optional:    true,
	},
	// ftpTargetTime
	"ftp_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for operation completion. Specified in milliseconds.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1000, 6000),
	},
	// ftpTimeLimit
	"ftp_time_limit": {
		Type:         schema.TypeInt,
		Description:  "Set the time limit for the test (in seconds). FTP tests default to 10s.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(10, 60),
	},
	// password
	"password-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The password to be used to authenticate with the destination server (required for FTP).",
	},
	// requestType
	"request_type": {
		Type:         schema.TypeString,
		Required:     true,
		Description:  "[download, upload, or list] Sets the type of activity for the test.",
		ValidateFunc: validation.StringInSlice([]string{"download", "upload", "list"}, false),
	},
	// url
	"url": {
		Type:        schema.TypeString,
		Description: "The target URL for the test.",
		Required:    true,
	},
	// useActiveFtp
	"use_active_ftp": {
		Type:        schema.TypeBool,
		Description: "Enables active FTP. If not set, tests default to use passive FTP.",
		Optional:    true,
	},
	// useExplicitFtps
	"use_explicit_ftps": {
		Type:        schema.TypeBool,
		Description: "Enables explicit FTPS (FTP over SSL). By default, tests will autodetect when it is appropriate to use FTPS.",
		Optional:    true,
	},
	// username
	"username-ftp": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "The username to be used to authenticate with the destination server.",
	},

	// HTTP SERVER

	// authType
	"auth_type": {
		Type:        schema.TypeString,
		Description: "[none, basic, ntlm, kerberos, oauth] The HTTP authentication type. Defaults to 'none'.",
		Optional:    true,
		Default:     "none",
		ValidateFunc: validation.StringInSlice([]string{
			"none",
			"basic",
			"ntlm",
			"kerberos",
			"oauth",
		}, false),
	},
	// agentInterfaces
	"agent_interfaces": { // New
		Description: "Agent interfaces",
		Optional:    true,
		Required:    false,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeString,
					Description: "The agent ID of the enterprise agent for the test.",
					Optional:    true,
				},
				"ip_address": {
					Type:        schema.TypeString,
					Description: "IP address of the agent interface.",
					Optional:    true,
				},
			},
		},
	},
	// clientCertificate
	"client_certificate": {
		Type:        schema.TypeString,
		Description: "String representation (containing newline characters) of the client certificate, if used.",
		Optional:    true,
	},
	// contentRegex
	"content_regex": {
		Type:        schema.TypeString,
		Description: "Verify content using a regular expression. This field does not require escaping.",
		Optional:    true,
		Required:    false,
		Default:     ".*",
	},
	// customHeaders
	"custom_headers": {
		Description: "The custom headers.",
		Optional:    true,
		Type:        schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"root": {
					Type:        schema.TypeMap,
					Description: "Use these HTTP headers for root server request.",
					Optional:    true,
				},
				"domains": {
					Type:        schema.TypeMap,
					Description: "Use these HTTP headers for the specified domains.",
					Optional:    true,
				},
				"all": {
					Type:        schema.TypeMap,
					Description: "Use these HTTP headers for all domains.",
					Optional:    true,
				},
			},
		},
	},
	// desiredStatusCode
	"desired_status_code": {
		Type:        schema.TypeString,
		Description: "The valid HTTP response code you’re interested in retrieving.",
		Optional:    true,
	},
	// dnsOverride
	"dns_override": {
		Type:         schema.TypeString,
		Description:  "The IP address to use for DNS override.",
		Optional:     true,
		ValidateFunc: validation.IsIPAddress,
	},
	// httpTargetTime
	"http_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for HTTP server completion, specified in milliseconds.",
		Optional:     true,
		Default:      1000,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	// httpTimeLimit
	"http_time_limit": {
		Type:        schema.TypeInt,
		Description: "The target time for HTTP server limits, specified in seconds.",
		Default:     5,
		Optional:    true,
	},
	// httpVersion
	"http_version": {
		Type:         schema.TypeInt,
		Description:  "Set to 2 for the default HTTP version (prefer HTTP/2), or 1 for HTTP/1.1 only.",
		Default:      2,
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 2),
	},
	// includeHeaders
	"include_headers": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to capture response headers for objects loaded by the test.",
		Optional:    true,
		Default:     true,
	},
	// oAuth
	"o_auth": oauth, // New
	// password
	"password": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The password to be used to authenticate with the destination server.",
	},
	// sslVersion
	"ssl_version": {
		Type:        schema.TypeString,
		Description: "Reflects the verbose ssl protocol version used by a test.",
		Computed:    true,
	},
	// sslVersionId
	"ssl_version_id": {
		Type:        schema.TypeString,
		Description: "Defines the SSL version. 0 for auto, 3 for SSLv3, 4 for TLS v1.0, 5 for TLS v1.1, 6 for TLS v1.2.",
		Optional:    true,
		Default:     "0",
		ValidateFunc: validation.StringInSlice([]string{
			"0",
			"3",
			"4",
			"5",
			"6",
		}, false),
	},
	// useNtlm
	"use_ntlm": {
		Type:        schema.TypeBool,
		Description: "Enable to use basic authentication. Only include this field if you are using authentication. Requires the username and password to be set if enabled.",
		Optional:    true,
	},
	// userAgent
	"user_agent": {
		Type:        schema.TypeString,
		Description: "The user-agent string to be provided during the test.",
		Optional:    true,
	},
	// username
	"username": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "The username to be used to authenticate with the destination server.",
	},
	// verifyCertificate
	"verify_certificate": {
		Type:        schema.TypeBool,
		Description: "Set whether to ignore certificate errors. Set to 'false' to ignore certificate errors. The default value is 'true'.",
		Optional:    true,
		Default:     true,
	},
	// allowUnsafeLegacyRenegotiation
	"allow_unsafe_legacy_renegotiation": { // New
		Type:        schema.TypeBool,
		Description: "Allows TLS renegotiation with servers not supporting RFC 5746. Default Set to true to allow unsafe legacy renegotiation.",
		Optional:    true,
		Default:     true,
	},
	// followRedirects
	"follow_redirects": {
		Type:        schema.TypeBool,
		Description: "Follow HTTP/301 or HTTP/302 redirect directives. Defaults to 'true'.",
		Optional:    true,
		Default:     true,
	},
	// overrideAgentProxy
	"override_agent_proxy": { // New
		Type:        schema.TypeBool,
		Description: "Flag indicating if a proxy other than the default should be used. To override the default proxy for agents, set to `true` and specify a value for `overrideProxyId`.",
		Optional:    true,
	},
	// overrideProxyId
	"override_proxy_id": { // New
		Type:        schema.TypeString,
		Description: "ID of the proxy to be used if the default proxy is overridden.",
		Optional:    true,
	},
	// collectProxyNetworkData
	"collect_proxy_network_data": {
		Type:        schema.TypeBool,
		Description: "Indicates whether network data to the proxy should be collected.",
		Optional:    true,
	},
	// headers
	"headers": {
		Type:        schema.TypeList,
		Description: "[\"header: value\", \"header2: value\"] The array of header strings.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Optional: true,
	},
	// postBody
	"post_body": {
		Type:        schema.TypeString,
		Description: "The POST body content. No escaping is required. If the post body is set to something other than empty, the requestMethod will be set to POST.",
		Optional:    true,
	},

	// PAGE LOAD

	// emulatedDeviceId
	"emulated_device_id": { // New
		Type:        schema.TypeString,
		Description: "ID of the emulated device, if one was given when the test was created.",
		Optional:    true,
	},
	// pageLoadTargetTime
	"page_load_target_time": { // New
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "Target time for page load completion, specified in seconds and cannot exceed the `pageLoadTimeLimit`.",
	},
	// pageLoadTimeLimit
	"page_load_time_limit": { // New
		Type:        schema.TypeInt,
		Required:    false,
		Optional:    true,
		Description: "Page load time limit. Must be larger than the `httpTimeLimit`.",
	},
	// blockDomains
	"block_domains": { // New
		Type:        schema.TypeString,
		Description: "Domains or full object URLs to be excluded from metrics and waterfall data for transaction tests.",
		Optional:    true,
	},
	// disableScreenshot
	"disable_screenshot": { // New
		Type:        schema.TypeBool,
		Description: "Enables or disables screenshots on error. Set true to not capture",
		Optional:    true,
	},
	// allowMicAndCamera
	"allow_mic_and_camera": { // New
		Type:        schema.TypeBool,
		Description: "Set true allow the use of a fake mic and camera in the browser.",
		Optional:    true,
	},
	// allowGeolocation
	"allow_geolocation": { // New
		Type:        schema.TypeBool,
		Description: "Set true to use the agent's geolocation by the web page.",
		Optional:    true,
	},
	// browserLanguage
	"browser_language": { // New
		Type:        schema.TypeString,
		Description: "Set one of the available browser language that you want to use to configure the browser.",
		Optional:    true,
	},
	// pageLoadingStrategy
	"page_loading_strategy": {
		Type:        schema.TypeString,
		Description: "[normal, eager or none] Defines page loading strategy. Defaults to 'none'.",
		Optional:    true,
		Default:     "none",
		ValidateFunc: validation.StringInSlice([]string{
			"normal",
			"eager",
			"none",
		}, false),
	},
	// httpInterval
	"http_interval": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "The interval to run the HTTP server test on.",
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
	},
	// subinterval
	"subinterval": {
		Type:         schema.TypeInt,
		Description:  "The subinterval for round-robin testing (in seconds). The value must be less than or equal to 'interval'.",
		Optional:     true,
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1200, 1800, 3600}),
	},

	// SIP SERVER

	// optionsRegex
	"options_regex": {
		Type:         schema.TypeString,
		Description:  "A regex string of options. This field does not require escaping.",
		Optional:     true,
		ValidateFunc: validation.StringIsValidRegExp,
	},
	// registerEnabled
	"register_enabled": {
		Type:        schema.TypeBool,
		Default:     false,
		Description: "Configure whether to perform SIP registration on the test target with the SIP REGISTER command. Default value is 'false'.",
		Optional:    true,
	},
	// sipTargetTime
	"sip_target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for test completion, specified in milliseconds.",
		Optional:     true,
		Default:      1000,
		ValidateFunc: validation.IntBetween(100, 5000),
	},
	// sipTimeLimit
	"sip_time_limit": {
		Type:         schema.TypeInt,
		Description:  "The test time limit. Can be between 5 and 10 seconds, and defaults to 5 seconds.",
		Optional:     true,
		Default:      10,
		ValidateFunc: validation.IntBetween(5, 10),
	},
	// authUser
	"auth_user": {
		Type:        schema.TypeString,
		Description: "The username for authentication with the SIP server.",
		Required:    true,
	},
	// protocol
	"protocol-sip": {
		Type:         schema.TypeString,
		Description:  "[tcp, tls, or udp] The transport layer for SIP communication. Can be either TCP, TLS (TLS over TCP), or UDP, and defaults to tcp.",
		Required:     true,
		ValidateFunc: validation.StringInSlice([]string{"tcp", "tls", "udp"}, false),
	},
	// sipRegistrar
	"sip_registrar": {
		Type:        schema.TypeString,
		Description: "The SIP server to be tested, specified by domain name or IP address.",
		Required:    true,
	},
	// user
	"user": {
		Type:        schema.TypeString,
		Description: "The username for SIP registration. This should be unique within a ThousandEyes account group.",
		Optional:    true,
	},

	// VOICE

	// codec
	"codec": {
		Type:        schema.TypeString,
		Description: "The label of the codec.",
		Computed:    true,
	},
	// codecId
	"codec_id": {
		Type:        schema.TypeString,
		Description: "The unique ID of the codec to use.",
		Optional:    true,
	},
	// duration
	"duration": {
		Type:         schema.TypeInt,
		Description:  "The duration of the test, in seconds (5 to 30).",
		Optional:     true,
		Default:      5,
		ValidateFunc: validation.IntBetween(5, 30),
	},
	// jitterBuffer
	"jitter_buffer": {
		Type:         schema.TypeInt,
		Description:  "The de-jitter buffer size, in seconds (0 to 150).",
		Optional:     true,
		Default:      40,
		ValidateFunc: validation.IntBetween(0, 150),
	},

	// WEB TRANSACTIONS

	// targetTime
	"target_time": {
		Type:         schema.TypeInt,
		Description:  "The target time for completion. The default value is 50 percent of the time limit, specified in seconds.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(1, 60),
	},
	// timeLimit
	"time_limit": {
		Type:         schema.TypeInt,
		Description:  "The time limit for the transaction. The default value is 30s.",
		Optional:     true,
		Default:      30,
		ValidateFunc: validation.IntBetween(5, 180),
	},
	// transactionScript
	"transaction_script": {
		Type:        schema.TypeString,
		Description: "The full selenium transaction script.",
		Required:    true,
	},
	// credentials
	"credentials": {
		Type:        schema.TypeList,
		Description: "The array of credentialID integers. You can get the credentialId from the /credentials endpoint.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
}
