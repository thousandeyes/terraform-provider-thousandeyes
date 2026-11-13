---
page_title: "Endpoint Agents and Tests"
subcategory: ""
---

# Endpoint Agents and Tests

Use this guide to import existing Endpoint Agents, look them up by name, and configure Endpoint Tests with individual agents or dynamic tags.

The examples cover these resources:

| Resource | Purpose |
| --- | --- |
| [`thousandeyes_endpoint_agent`](../resources/endpoint_agent.md) | Import and manage an existing Endpoint Agent. |
| [`thousandeyes_endpoint_agent_to_server`](../resources/endpoint_agent_to_server.md) | Schedule a test against a fixed server. |
| [`thousandeyes_endpoint_agent_to_server_dynamic`](../resources/endpoint_agent_to_server_dynamic.md) | Monitor the servers used by a supported application. |
| [`thousandeyes_endpoint_http_server`](../resources/endpoint_http_server.md) | Schedule an HTTP Server test. |
| [`thousandeyes_tag`](../resources/tag.md) | Define a dynamic tag that selects Endpoint Agents. |

For product concepts, see the ThousandEyes documentation for [Endpoint Agents](https://docs.thousandeyes.com/product-documentation/endpoint-experience/endpoint-agents).

## Before you begin

Use a provider release that includes the resources in this guide. The `import` block example requires Terraform 1.5 or later; the `terraform import` command is also supported.

Install and register an Endpoint Agent outside Terraform before importing or looking it up. Set `TE_TOKEN` to your ThousandEyes API token and, when needed, set `TE_AID` to the account group containing the agent. Use the same account group for the agents, tags, and tests in these examples.

```terraform
terraform {
  required_providers {
    thousandeyes = {
      source = "thousandeyes/thousandeyes"
    }
  }
}

provider "thousandeyes" {}
```

Run `terraform init` to install the provider. Replace the example agent ID and exact agent name with values from your account. The examples below share this provider configuration; choose the resource examples you need.

## Import an existing Endpoint Agent

The API has no Endpoint Agent create operation. Declaring a resource and applying it without an import returns an error; Terraform does not install or register the agent.

Start with a resource block that leaves the remote name, license, and status unchanged:

```terraform
resource "thousandeyes_endpoint_agent" "laptop" {
  deletion_protection = true
}

import {
  to = thousandeyes_endpoint_agent.laptop
  id = "123456"
}
```

Run `terraform plan`, review the imported agent and the deletion protection setting, then run `terraform apply`. The ID is the existing Endpoint Agent ID, not its name or a test ID.

Alternatively, omit the `import` block and run:

```shell
terraform import thousandeyes_endpoint_agent.laptop 123456
terraform plan
terraform apply
```

With either import method, apply `deletion_protection = true` before removing the resource or destroying it. Importing alone does not apply this setting.

### Manage an imported agent

After importing, you can replace the resource block above with the following configuration to manage its name, supported license type, and status:

```terraform
resource "thousandeyes_endpoint_agent" "laptop" {
  name                = "London laptop"
  license_type        = "advantage"
  status              = "enabled"
  deletion_protection = true
}
```

Set `status = "disabled"` to disable the agent, or `status = "enabled"` to enable it. Review the plan before applying name, license, or status changes. Omitted mutable properties retain the values read from the API.

Agents with the `embed-wireless-advantage` or `embed-wireless-essentials` license types can be imported and read, but those license types cannot be assigned. Leave `license_type` unset to retain an existing wireless license.

### Remove an agent from Terraform or delete it remotely

`deletion_protection` defaults to `false`. Destroying the resource, including removing its block and applying, deletes the remote Endpoint Agent when protection is disabled.

To stop managing the agent while keeping it in ThousandEyes:

1. Set `deletion_protection = true` in the resource block and apply the change.
2. Remove the resource block and any `import` block targeting it, then apply again.

With protection enabled in state, the provider removes the agent from Terraform state without calling the remote delete operation. This setting does not block Terraform destroy. To delete the remote agent instead, apply `deletion_protection = false` before removing or destroying the resource.

## Look up an Endpoint Agent by exact name

Use the [`thousandeyes_endpoint_agent` data source](../data-sources/endpoint_agent.md) when a test needs an existing agent's ID without managing that agent's lifecycle:

```terraform
data "thousandeyes_endpoint_agent" "laptop" {
  agent_name = "London laptop"
}

resource "thousandeyes_endpoint_agent_to_server_dynamic" "webex" {
  test_name           = "Endpoint Webex"
  application         = "webex"
  interval            = 300
  agent_selector_type = "specific-agents"
  agents              = [data.thousandeyes_endpoint_agent.laptop.agent_id]
  max_machines        = 25
}
```

The data source matches the exact agent name, not a substring, and reports an error when no agent matches. Endpoint Agent IDs are strings. An imported agent's ID is also available as `thousandeyes_endpoint_agent.laptop.id` and can be used in `agents` instead.

Dynamic Agent-to-Server tests support `application = "webex"`, `"zoom"`, or `"microsoft-teams"`. The application determines the destination servers; these tests do not take a fixed `server` or `url`.

## Choose an agent selector

All three Endpoint Test resources support the following selector modes. Set `agent_selector_type` explicitly in each example.

| `agent_selector_type` | `agents` | `tag_ids` |
| --- | --- | --- |
| `all-agents` | Omit. | Omit. |
| `specific-agents` | Set a nonempty set of Endpoint Agent IDs. | Omit. |
| `agent-tags` | Omit. | Set a nonempty set of dynamic Endpoint Agent tag IDs. |

`max_machines` is optional for every selector and defaults to 25. It limits how many selected agents can execute the test; `all-agents` does not remove that limit. Set a positive value when you need a different maximum.

Endpoint Labels, `agent-labels`, and the deprecated API input `endpointAgentLabels` are unsupported. Use dynamic Endpoint Agent tags for selection by matching agent properties. Static Endpoint Agent tag assignment is also unsupported.

### Scheduled Agent-to-Server test with all agents

```terraform
resource "thousandeyes_endpoint_agent_to_server" "connectivity" {
  test_name           = "Endpoint connectivity"
  server              = "example.com"
  port                = 443
  protocol            = "tcp"
  interval            = 300
  agent_selector_type = "all-agents"
}
```

This example omits `agents` and `tag_ids` and uses the default `max_machines` of 25. To target individual agents instead, change the selector to `specific-agents` and add `agents = [data.thousandeyes_endpoint_agent.laptop.agent_id]` using the lookup above.

### Scheduled HTTP Server test with a dynamic tag

Define a tag whose filters match the Endpoint Agents you want to use, then reference its ID in the test:

```terraform
resource "thousandeyes_tag" "windows_laptops" {
  key         = "platform"
  value       = "windows-laptops"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = "and"

  filters {
    key    = "platform"
    values = ["Windows"]
    mode   = "in"
    scope  = "default"
  }
}

resource "thousandeyes_endpoint_http_server" "website" {
  test_name           = "Endpoint website"
  url                 = "https://example.com"
  interval            = 300
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.windows_laptops.id]
  max_machines        = 50
}
```

Dynamic tags select agents automatically according to their filter criteria. Use `match_type = "and"` to require every filter to match or `"or"` to require any filter to match. At least one `filters` block is required. Its `scope` defaults to `default` for built-in agent properties; use `custom` for metadata supplied when an agent checks in. The [tag resource reference](../resources/tag.md) describes filter modes and value semantics.

Terraform uses the tag reference to create the tag before the test. Choose filter criteria that match agents in your account. Changing the tag's filter criteria changes which agents are eligible for tests that reference the tag.

### Use the same selectors with other Endpoint Tests

Selector configuration is shared across the three test resource types. For example, the following scheduled and dynamic Agent-to-Server tests use the tag defined above:

```terraform
resource "thousandeyes_endpoint_agent_to_server" "tagged_connectivity" {
  test_name           = "Tagged Endpoint connectivity"
  server              = "example.com"
  port                = 443
  protocol            = "tcp"
  interval            = 300
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.windows_laptops.id]
}

resource "thousandeyes_endpoint_agent_to_server_dynamic" "tagged_zoom" {
  test_name           = "Tagged Endpoint Zoom"
  application         = "zoom"
  interval            = 300
  agent_selector_type = "agent-tags"
  tag_ids             = [thousandeyes_tag.windows_laptops.id]
}
```

For any of the three resources, use `specific-agents` with `agents` to select individual agents, or `all-agents` with neither list. Keep the resource's target fields: `server` for scheduled Agent-to-Server, `application` for dynamic Agent-to-Server, and `url` for scheduled HTTP Server.

## Review changes before applying

Run `terraform plan` before applying changes. Changing `agent_selector_type`, `agents`, `tag_ids`, or `max_machines` replaces the Endpoint Test. Updating the criteria of a referenced dynamic tag can change agent selection while retaining the same test and tag IDs.

Other update and replacement behavior depends on the resource. Review its schema and lifecycle notes, especially when changing HTTP authentication or URLs. Tests configured with `is_enabled = false` may briefly be enabled during creation. Test enablement uses `is_enabled`; the Endpoint Agent resource uses `status`.

After applying, run `terraform plan` again to check that configuration and state agree. Deleting an Endpoint Test deletes the test, while deleting a tag deletes the tag. Neither operation deletes the Endpoint Agents selected by that test or tag.
