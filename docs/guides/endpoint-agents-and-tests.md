---
page_title: "Endpoint Agents and Tests"
subcategory: ""
---

# Endpoint Agents and Tests

Use this guide to import an Endpoint Agent and configure Endpoint Tests. Agents must already be installed and registered because Terraform cannot create them.

## Set up the provider

Set `TE_TOKEN` to your API token. If needed, set `TE_AID` to the account group containing your agents. Use that group for the tags and tests too.

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

Run `terraform init`, then replace the sample agent ID and name below with values from your account. The examples share this provider configuration.

## Import an agent

The [Endpoint Agent resource](../resources/endpoint_agent.md) manages existing agents. Applying it without an import returns an error because the API has no create operation.

With Terraform 1.5 or later, use an `import` block:

```terraform
resource "thousandeyes_endpoint_agent" "laptop" {
  deletion_protection = true
}

import {
  to = thousandeyes_endpoint_agent.laptop
  id = "123456"
}
```

Run `terraform plan` and review the agent, then run `terraform apply`. You can also omit the `import` block and use this command before planning and applying:

```shell
terraform import thousandeyes_endpoint_agent.laptop 123456
```

### Manage the agent

After importing, replace the resource block to manage its name, license, and status:

```terraform
resource "thousandeyes_endpoint_agent" "laptop" {
  name                = "London laptop"
  license_type        = "advantage"
  status              = "enabled"
  deletion_protection = true
}
```

Use `status = "disabled"` to disable the agent. Omitting `name`, `license_type`, or `status` keeps its remote value. The `embed-wireless-advantage` and `embed-wireless-essentials` license types can be read but cannot be assigned, so leave `license_type` unset to retain either license.

### Remove or delete the agent

By default, `deletion_protection` is `false`, so destroying the resource deletes the remote agent. To keep the agent, first apply `deletion_protection = true`. Then remove the resource block and any `import` block and apply again. Terraform removes the agent from state while keeping it in ThousandEyes.

Importing alone does not apply deletion protection. To delete a protected agent remotely, first apply `deletion_protection = false` and then destroy the resource.

## Look up an agent and configure a dynamic test

Use the [Endpoint Agent data source](../data-sources/endpoint_agent.md) to find an agent by exact name. Then pass its `agent_id` to a [dynamic Agent to Server test](../resources/endpoint_agent_to_server_dynamic.md):

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

The lookup returns an error if no agent matches. If you imported the agent instead, use `thousandeyes_endpoint_agent.laptop.id`. Dynamic tests follow the servers used by `webex`, `zoom`, or `microsoft-teams`.

## Choose an agent selector

All three test resources support the same selectors. Use `all-agents` and omit both `agents` and `tag_ids`. For individual agents, use `specific-agents` with a nonempty `agents` list and omit `tag_ids`. For dynamic tags, use `agent-tags` with a nonempty `tag_ids` list and omit `agents`.

`max_machines` is optional and defaults to 25 for every selector. Set a positive value to change the maximum number of agents that can run the test.

Endpoint Labels, `agent-labels`, and the deprecated `endpointAgentLabels` input are unsupported. Use dynamic tags for selection by agent properties.

## Configure a scheduled Agent to Server test

This [scheduled Agent to Server test](../resources/endpoint_agent_to_server.md) uses `all-agents`. It omits `max_machines`, so the maximum is 25.

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

## Configure an HTTP Server test with a tag

Create a [dynamic Endpoint Agent tag](../resources/tag.md), then pass its ID to a [scheduled HTTP Server test](../resources/endpoint_http_server.md):

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

The tag selects Windows agents. If you add more filters, `match_type = "and"` requires all of them to match. Use `or` if any filter may match. The same `tag_ids` setting works with both Agent to Server test resources.

## Apply changes

Run `terraform plan` before applying. Changing `agent_selector_type`, `agents`, `tag_ids`, or `max_machines` replaces the test. Updating a tag's filters changes eligible agents while keeping the test. After applying, run `terraform plan` again to check for drift.
