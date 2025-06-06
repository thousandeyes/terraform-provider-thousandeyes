---
page_title: "thousandeyes_agent_to_server Resource - terraform-provider-thousandeyes"
subcategory: ""
description: |-
---

# thousandeyes_agent_to_server (Resource)

This resource allows you to create and configure an agent-to-server test. This test type measures network performance as seen from ThousandEyes agent(s) towards a remote server. For more information about agent-to-server tests, see [Agent-to-Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#agent-to-server-test).

## Example Usage

```terraform
resource "thousandeyes_agent_to_server" "example_agent_to_server_test" {
  test_name      = "Example Agent to Server test set from Terraform provider"
  interval       = 120
  alerts_enabled = false
  server         = "www.thousandeyes.com"
  port           = 443
  agents         = ["3"] # Singapore
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agents` (Set of String) The list of ThousandEyes agent IDs to use.
- `interval` (Number) The interval to run the test on, in seconds.
- `server` (String) The target host.

### Optional

- `alert_rules` (Set of String) List of alert rules IDs to apply to the test (get `ruleId` from `/alerts/rules` endpoint. If `alertsEnabled` is set to `true` and `alertRules` is not included on test creation or update, applicable user default alert rules will be used)
- `alerts_enabled` (Boolean) Set to 'true' to enable alerts, or 'false' to disable alerts. The default value is 'true'.
- `bandwidth_measurements` (Boolean) Set to `true` to measure bandwidth. This only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set. Defaults to 'false'.
- `bgp_measurements` (Boolean) Enable BGP measurements. Set to true for enabled, false for disabled.
- `continuous_mode` (Boolean) To enable continuous monitoring, set this parameter to `true`.  When continuous monitoring is enabled, the following actions occur: * `fixedPacketRate` is enforced * `bandwidthMeasurements` are disabled * If the `protocol` is set to `tcp`, `probeMode` is set to `syn`.
- `description` (String) A description of the alert rule. Defaults to an empty string.
- `dscp_id` (String) The DSCP ID.
- `enabled` (Boolean) Enables or disables the test.
- `fixed_packet_rate` (Number) Sets packets rate sent to measure the network in packets per second.
- `ipv6_policy` (String) [force-ipv4, prefer-ipv6, force-ipv6, or use-agent-policy] IP version policy. Overrides the IPv6 policy configured at the agent level.
- `monitors` (Set of String) Contains list of BGP monitor IDs (get `monitorId` from `/monitors` endpoint)
- `mtu_measurements` (Boolean) Measure MTU sizes on the network from agents to the target.
- `network_measurements` (Boolean) Set to 'true' to enable network measurements.
- `num_path_traces` (Number) The number of path traces.
- `path_trace_mode` (String) [classic or in-session] Choose 'inSession' to perform the path trace within a TCP session. Default value is 'classic'.
- `ping_payload_size` (Number) Payload size (not total packet size) for the end-to-end metrics probes, ping payload size allows values from 0 to 1400 bytes. When set to null, payload sizes are 0 bytes for ICMP-based tests and 1 byte for TCP-based tests.
- `port` (Number) The target port.
- `probe_mode` (String) [auto, sack, or syn] The probe mode used by end-to-end network tests. This is only valid if the protocol is set to TCP. The default value is AUTO.
- `protocol` (String) The protocol used by dependent network tests (end-to-end, path trace, PMTUD). Default value is tcp.
- `randomized_start_time` (Boolean) Indicates whether agents should randomize the start time in each test round.
- `shared_with_accounts` (Set of String) List of accounts
- `test_name` (String) The name of the test.
- `use_public_bgp` (Boolean) Enable to automatically add all available Public BGP Monitors to the test.

### Read-Only

- `created_by` (String) Created by user.
- `created_date` (String) The date of creation.
- `dscp` (String) The Differentiated Services Code Point (DSCP) label.
- `id` (String) The ID of this resource.
- `labels` (Set of String, Deprecated) ["1", "2"] The array of labels.
- `link` (String) Its value is either a URI [RFC3986] or a URI template [RFC6570].
- `live_share` (Boolean) Set to 'true' for a test shared with your account group, or to 'false' for a normal test.
- `modified_by` (String) Last modified by this user.
- `modified_date` (String) The date the test was last modified. Shown in UTC.
- `saved_event` (Boolean) Set to 'true' for a saved event, or to 'false' for a normal test.
- `test_id` (String) The unique ID of the test.
- `type` (String) The type of test.

## Import
In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) providing `resource_id`.
```terraform
import {
  to = thousandeyes_agent_to_server.example_agent_to_server_test
  id = "resource_id"
}
```

Using `terraform import` command.
```shell
terraform import thousandeyes_agent_to_server.example_agent_to_server_test resource_id
```
