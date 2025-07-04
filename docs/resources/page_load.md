---
page_title: "thousandeyes_page_load Resource - terraform-provider-thousandeyes"
subcategory: ""
description: |-
---

# thousandeyes_page_load (Resource)

This resource allows you to create a page load test. This test type obtains in-browser site performance metrics. For more information, see [Page Load Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#page-load-test).

## Example Usage

```terraform
resource "thousandeyes_page_load" "example_page_load_test" {
  test_name      = "Example Page Load test set from Terraform provider"
  alerts_enabled = false
  url            = "https://www.thousandeyes.com"
  interval       = 120
  http_interval  = 120
  agents         = ["3"] # Singapore
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agents` (Set of String) The list of ThousandEyes agent IDs to use.
- `http_interval` (Number) The interval to run the HTTP server test on.
- `interval` (Number) The interval to run the test on, in seconds.
- `url` (String) The target URL for the test.

### Optional

- `agent_interfaces` (Block Set) Agent interfaces (see [below for nested schema](#nestedblock--agent_interfaces))
- `alert_rules` (Set of String) List of alert rules IDs to apply to the test (get `ruleId` from `/alerts/rules` endpoint. If `alertsEnabled` is set to `true` and `alertRules` is not included on test creation or update, applicable user default alert rules will be used)
- `alerts_enabled` (Boolean) Set to 'true' to enable alerts, or 'false' to disable alerts. The default value is 'true'.
- `allow_geolocation` (Boolean) Set true to use the agent's geolocation by the web page.
- `allow_mic_and_camera` (Boolean) Set true allow the use of a fake mic and camera in the browser.
- `allow_unsafe_legacy_renegotiation` (Boolean) Allows TLS renegotiation with servers not supporting RFC 5746. Default Set to true to allow unsafe legacy renegotiation.
- `auth_type` (String) [none, basic, ntlm, kerberos, oauth] The HTTP authentication type. Defaults to 'none'.
- `bandwidth_measurements` (Boolean) Set to `true` to measure bandwidth. This only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set. Defaults to 'false'.
- `bgp_measurements` (Boolean) Enable BGP measurements. Set to true for enabled, false for disabled.
- `block_domains` (String) Domains or full object URLs to be excluded from metrics and waterfall data for transaction tests.
- `browser_language` (String) Set one of the available browser language that you want to use to configure the browser.
- `client_certificate` (String) String representation (containing newline characters) of client certificate, the private key must be placed first, then the certificate.
- `collect_proxy_network_data` (Boolean) Indicates whether network data to the proxy should be collected.
- `content_regex` (String) Verify content using a regular expression. This field does not require escaping.
- `custom_headers` (Block Set) The custom headers. (see [below for nested schema](#nestedblock--custom_headers))
- `description` (String) A description of the alert rule. Defaults to an empty string.
- `desired_status_code` (String) The valid HTTP response code you’re interested in retrieving.
- `disable_screenshot` (Boolean) Enables or disables screenshots on error. Set true to not capture
- `distributed_tracing` (Boolean) Adds distributed tracing headers to API requests using B3 and W3C standards.
- `dns_override` (String) The IP address to use for DNS override.
- `download_limit` (Number) Specify the maximum number of bytes to download from the target object.
- `emulated_device_id` (String) ID of the emulated device, if one was given when the test was created.
- `enabled` (Boolean) Enables or disables the test.
- `fixed_packet_rate` (Number) Sets packets rate sent to measure the network in packets per second.
- `follow_redirects` (Boolean) Follow HTTP/301 or HTTP/302 redirect directives. Defaults to 'true'.
- `http_target_time` (Number) The target time for HTTP server completion, specified in milliseconds.
- `http_time_limit` (Number) The target time for HTTP server limits, specified in seconds.
- `http_version` (Number) Set to 2 for the default HTTP version (prefer HTTP/2), or 1 for HTTP/1.1 only.
- `include_headers` (Boolean) Set to 'true' to capture response headers for objects loaded by the test.
- `monitors` (Set of String) Contains list of BGP monitor IDs (get `monitorId` from `/monitors` endpoint)
- `mtu_measurements` (Boolean) Measure MTU sizes on the network from agents to the target.
- `network_measurements` (Boolean) Set to 'true' to enable network measurements.
- `num_path_traces` (Number) The number of path traces.
- `override_agent_proxy` (Boolean) Flag indicating if a proxy other than the default should be used. To override the default proxy for agents, set to `true` and specify a value for `overrideProxyId`.
- `override_proxy_id` (String) ID of the proxy to be used if the default proxy is overridden.
- `page_load_target_time` (Number) Target time for page load completion, specified in seconds and cannot exceed the `pageLoadTimeLimit`.
- `page_load_time_limit` (Number) Page load time limit. Must be larger than the `httpTimeLimit`.
- `page_loading_strategy` (String) [normal, eager or none] Defines page loading strategy. Defaults to 'none'.
- `password` (String, Sensitive) The password to be used to authenticate with the destination server.
- `path_trace_mode` (String) [classic or in-session] Choose 'inSession' to perform the path trace within a TCP session. Default value is 'classic'.
- `probe_mode` (String) [auto, sack, or syn] The probe mode used by end-to-end network tests. This is only valid if the protocol is set to TCP. The default value is AUTO.
- `protocol` (String) The protocol used by dependent network tests (end-to-end, path trace, PMTUD). Default value is tcp.
- `randomized_start_time` (Boolean) Indicates whether agents should randomize the start time in each test round.
- `shared_with_accounts` (Set of String) List of accounts
- `ssl_version_id` (String) Defines the SSL version. 0 for auto, 3 for SSLv3, 4 for TLS v1.0, 5 for TLS v1.1, 6 for TLS v1.2.
- `subinterval` (Number) The subinterval for round-robin testing (in seconds). The value must be less than or equal to 'interval'.
- `test_name` (String) The name of the test.
- `use_ntlm` (Boolean) Enable to use basic authentication. Only include this field if you are using authentication. Requires the username and password to be set if enabled.
- `use_public_bgp` (Boolean) Enable to automatically add all available Public BGP Monitors to the test.
- `user_agent` (String) The user-agent string to be provided during the test.
- `username` (String) The username to be used to authenticate with the destination server.
- `verify_certificate` (Boolean) Set whether to ignore certificate errors. Set to 'false' to ignore certificate errors. The default value is 'true'.

### Read-Only

- `created_by` (String) Created by user.
- `created_date` (String) The date of creation.
- `id` (String) The ID of this resource.
- `labels` (Set of String, Deprecated) ["1", "2"] The array of labels.
- `link` (String) Its value is either a URI [RFC3986] or a URI template [RFC6570].
- `live_share` (Boolean) Set to 'true' for a test shared with your account group, or to 'false' for a normal test.
- `modified_by` (String) Last modified by this user.
- `modified_date` (String) The date the test was last modified. Shown in UTC.
- `saved_event` (Boolean) Set to 'true' for a saved event, or to 'false' for a normal test.
- `ssl_version` (String) Reflects the verbose ssl protocol version used by a test.
- `test_id` (String) The unique ID of the test.
- `type` (String) The type of test.

<a id="nestedblock--agent_interfaces"></a>
### Nested Schema for `agent_interfaces`

Optional:

- `agent_id` (String) The agent ID of the enterprise agent for the test.
- `ip_address` (String) IP address of the agent interface.


<a id="nestedblock--custom_headers"></a>
### Nested Schema for `custom_headers`

Optional:

- `all` (Map of String, Sensitive) Use these HTTP headers for all domains.
- `domains` (Map of String, Sensitive) Use these HTTP headers for the specified domains.
- `root` (Map of String, Sensitive) Use these HTTP headers for root server request.

## Import
In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) providing `resource_id`.
```terraform
import {
  to = thousandeyes_page_load.example_page_load_test
  id = "resource_id"
}
```

Using `terraform import` command.
```shell
terraform import thousandeyes_page_load.example_page_load_test resource_id
```
