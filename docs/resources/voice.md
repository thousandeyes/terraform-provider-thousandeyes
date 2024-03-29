---
page_title: "thousandeyes_voice Resource - terraform-provider-thousandeyes"
subcategory: ""
description: |-
---

# thousandeyes_voice (Resource)

This resource allows you to create a RTP Stream test. This test type measures the quality of real-time protocol (RTP) voice streams between ThousandEyes agents that act as VoIP user agents. For more information, see [RTP Stream Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#rtp-stream-test).

## Example Usage

```terraform
resource "thousandeyes_voice" "example_voice_test" {
  test_name        = "Example RTP stream test set from Terraform provider"
  interval         = 120
  alerts_enabled   = false

  bgp_measurements = true
  use_public_bgp   = true

  target_agent_id = 4 # Tokyo

  agents {
    agent_id = 3 # Singapore
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `agents` (Block Set, Min: 1) The list of ThousandEyes agents to use. (see [below for nested schema](#nestedblock--agents))
- `interval` (Number) The interval to run the test on, in seconds.
- `target_agent_id` (Number) The target agent's unique ID. Pulled from the /agents endpoint. Both the 'agents': [] and the targetAgentID cannot be Cloud Agents. Can be Enterprise Agent -> Cloud, Cloud -> Enterprise Agent, or Enterprise Agent -> Enterprise Agent.
- `test_name` (String) The name of the test.

### Optional

- `alert_rules` (Block Set) Gets the ruleId from the /alert-rules endpoint. If alertsEnabled is set to 'true' and alertRules is not included in a creation/update query, the applicable defaults will be used. (see [below for nested schema](#nestedblock--alert_rules))
- `alerts_enabled` (Boolean) Set to 'true' to enable alerts, or 'false' to disable alerts. The default value is 'true'.
- `bgp_measurements` (Boolean) Enable BGP measurements. Set to true for enabled, false for disabled.
- `bgp_monitors` (Block List) The array of BGP monitor object IDs. The monitorIDs can be sourced from the /bgp-monitors endpoint. (see [below for nested schema](#nestedblock--bgp_monitors))
- `codec_id` (Number) The unique ID of the codec to use.
- `description` (String) A description of the alert rule. Defaults to an empty string.
- `dscp_id` (Number) The DSCP ID.
- `duration` (Number) The duration of the test, in seconds (5 to 30).
- `enabled` (Boolean) Enables or disables the test.
- `jitter_buffer` (Number) The de-jitter buffer size, in seconds (0 to 150).
- `mtu_measurements` (Boolean) Measure MTU sizes on the network from agents to the target.
- `num_path_traces` (Number) The number of path traces.
- `shared_with_accounts` (Block List) [“serverName”: “fqdn of server”] The array of DNS Server objects. (see [below for nested schema](#nestedblock--shared_with_accounts))
- `use_public_bgp` (Boolean) Enable to automatically add all available Public BGP Monitors to the test.

### Read-Only

- `api_links` (List of Object) Self links to the endpoint to pull test metadata, and data links to the endpoint for test data. Read-only, and shows rel and href elements. (see [below for nested schema](#nestedatt--api_links))
- `codec` (String) The label of the codec.
- `created_by` (String) Created by user.
- `created_date` (String) The date of creation.
- `dscp` (String) The Differentiated Services Code Point (DSCP) label.
- `groups` (Set of Object) The array of label objects. (see [below for nested schema](#nestedatt--groups))
- `id` (String) The ID of this resource.
- `live_share` (Boolean) Set to 'true' for a test shared with your account group, or to 'false' for a normal test.
- `modified_by` (String) Last modified by this user.
- `modified_date` (String) The date the test was last modified. Shown in UTC.
- `saved_event` (Boolean) Set to 'true' for a saved event, or to 'false' for a normal test.
- `test_id` (Number) The unique ID of the test.
- `type` (String) The type of test.

<a id="nestedblock--agents"></a>
### Nested Schema for `agents`

Required:

- `agent_id` (Number) The unique ID for the ThousandEyes agent.


<a id="nestedblock--alert_rules"></a>
### Nested Schema for `alert_rules`

Optional:

- `rule_id` (Number) The unique ID of the alert rule.


<a id="nestedblock--bgp_monitors"></a>
### Nested Schema for `bgp_monitors`

Required:

- `monitor_id` (Number) The unique ID of the BGP monitor.

Optional:

- `ip_address` (String) The IP address of the BGP monitor.
- `monitor_name` (String) The name of the BGP monitor.
- `monitor_type` (String) [Public or Private] Shows the type of BGP monitor.
- `network` (String) The name of the autonomous system in which the BGP monitor is found.


<a id="nestedblock--shared_with_accounts"></a>
### Nested Schema for `shared_with_accounts`

Required:

- `aid` (Number) The account group ID.

Read-Only:

- `name` (String) Account name.


<a id="nestedatt--api_links"></a>
### Nested Schema for `api_links`

Read-Only:

- `href` (String)
- `rel` (String)


<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

Read-Only:

- `agents` (List of Object) (see [below for nested schema](#nestedobjatt--groups--agents))
- `builtin` (Boolean)
- `group_id` (Number)
- `name` (String)
- `tests` (List of Object) (see [below for nested schema](#nestedobjatt--groups--tests))
- `type` (String)

<a id="nestedobjatt--groups--agents"></a>
### Nested Schema for `groups.agents`

Read-Only:

- `agent_id` (Number)


<a id="nestedobjatt--groups--tests"></a>
### Nested Schema for `groups.tests`

Read-Only:

- `test_id` (Number)


