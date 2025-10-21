---
page_title: "ThousandEyes Terraform Provider v3 Upgrade Guide"
subcategory: ""
---

# ThousandEyes Terraform Provider Version 3 Upgrade Guide

Version 3.x of the ThousandEyes Terraform provider introduces significant changes. This guide covers breaking changes in v3.0.0 and subsequent v3.x releases. Review the relevant sections based on your upgrade path to ensure a smooth migration.

Upgrade topics:

- [Provider Version Configuration](#provider-version-configuration)
- [Provider Arguments](#provider-arguments)
- [Labels to Tags migration](#labels-to-tags-migration)
- [Integrations Data source](#integrations-data-source)
- [Test schema changes](#test-schema-changes)
- [Alert Rule webhook notification changes](#alert-rule-webhook-notification-changes)
- [Other notable changes](#other-notable-changes)
- [Additional resources](#additional-resources)

## Provider version configuration

Update your provider version constraint and run [`terraform init -upgrade`](https://developer.hashicorp.com/terraform/cli/commands/init) to download the new version.

**Example:**

```terraform
terraform {
  required_providers {
    thousandeyes = {
      source  = "thousandeyes/thousandeyes"
      version = "~> 3.0"
    }
  }
}

provider "thousandeyes" {
  # Configuration options
}
```

## Provider Arguments

The `api_sdk_logs_enabled` argument was added. The default value is `false`.

## Labels to Tags migration

Both Groups and Labels have been deprecated in ThousandEyes API v7 and are no longer supported in the current version of the provider. Instead, [Tags](https://developer.cisco.com/docs/thousandeyes/tags-api-overview/) are now used as a replacement for Labels.
Tags maintain backwards compatibility with any existing Groups and Labels. However, Tags can no longer be assigned directly to tests using its own resource; instead, the `thousandeyes_tag_assignment` resource should be used for this purpose.

**Before (using `thousandeyes_label`):**

```hcl
resource "thousandeyes_label" "example_label" {
  # ... other configuration ...
  
  tests {
    test_id = 123456
  }
}
```

**After (using `thousandeyes_tag` and `thousandeyes_tag_assignment`):**

```hcl
resource "thousandeyes_tag" "example_tag" {
  key         = "Example Tag Key"
  value       = "Example Tag Value"
  object_type = "test"
  color       = "#b3de69"
  access_type = "all"
  icon        = "LABEL"
}

resource "thousandeyes_tag_assignment" "example_assignment" {
  tag_id = data.thousandeyes_tag.example_tag.id

  assignments {
    id   = "123456" #Id of existing entity (Test, Dashboard, etc.)
    type = "test"
  }
}
```

## Integrations Data source

With the retirement of the Integrations API, the `thousandeyes_integration` data source has been removed. You can still reference existing integrations by their ID.

**Before (using `thousandeyes_integration`):**

```hcl
data "thousandeyes_integration" "example_integration" {
  integration_name = "test-pd-service"
}

resource "thousandeyes_alert_rule" "example_alert_rule" {
  # ... other configuration ...

  notifications {
    third_party {
      integration_id = data.thousandeyes_integration.example_integration.id
      integration_type = data.thousandeyes_integration.example_integration.integration_type
    }
  }
}
```

**After:**

```hcl
resource "thousandeyes_alert_rule" "example_alert_rule" {
  # ... other configuration ...

  notifications {
    third_party {
      integration_id = "integration-id" # Replace with the actual integration ID
      integration_type = "integration_type" # Replace with the actual integration type
    }
  }
}
```

## Test schema changes

The following attributes in test resources have been updated for consistency and alignment with the ThousandEyes v7 API. Review and update your configurations accordingly:

- **`alert_rules`**: Now a set of strings representing alert rule IDs. Previously, this was a block set with objects containing a `rule_id` field.
- **`agents`**: Now a set of strings representing agent IDs. Previously, this was a block set with objects containing an `agent_id` field.
- **`bgp_monitors` → `monitors`**: The attribute has been renamed to `monitors` and is now a set of strings representing BGP monitor IDs, instead of a block set with objects containing a `monitor_id` field.
- **`shared_with_accounts`**: Now a set of strings representing account IDs. Previously, this was a block list with objects containing an `aid` field.
- **`dns_servers`**: Now a set of strings representing DNS server names. Previously, this was a block set with objects containing a `server_name` and optionally a `server_id`.
- **`use_public_bgp`**: The default value for this attribute is no longer set in the resource schema. Instead, the provider now relies on the default behavior defined by the ThousandEyes API.
- **`api_links` → `link`**: The attribute has been renamed to `link`.
- **`target_agent_id`**: Now a string value instead of a numeric value.
- **`dscp_id`**: Now a string value instead of a numeric value.
- **`test_id`**: Now a string value instead of a numeric value.

**Before:**

```hcl
resource "thousandeyes_dns_server" "dns_server_test" {
  test_name      = "Example DNS server test"
  interval       = 120
  domain         = "www.thousandeyes.com ANY"
  dns_servers {
    server_name = "ns1.google.com"
  }
  agents {
    agent_id = 123
  }
  alert_rules {
    rule_id = 123
  }
  shared_with_accounts  {
    aid = 123
  }
  bgp_monitors {
    monitor_id = 123
  }
}
```

**After:**

```hcl
resource "thousandeyes_dns_server" "dns_server_test" {
  test_name      = "Example DNS server test"
  interval       = 120
  domain         = "www.thousandeyes.com ANY"
  dns_servers    = ["ns1.google.com"]
  agents         = ["123"]
  alert_rules    = ["123"]
  shared_with_accounts = ["123"]
  monitors       = ["123"]
}
```

## Alert Rule webhook notification changes

> **Note:** This change only affects users upgrading from v3.0.x to v3.1.0 or later. These fields were introduced as a regression in v3.0.x and are being removed in v3.1.0. If you are upgrading directly from v2.x to v3.1.0+, you will not encounter this issue as these fields never existed in v2.x.

The `integration_name` and `target` fields have been removed from webhook and custom webhook notifications in alert rules. These fields were unintentionally introduced in v3.0.x and were read-only, returned by the ThousandEyes API based on the integration configuration. They could not be set through alert rules and caused configuration drift.

**Breaking change in v3.1.0 (regression fix):** If you are using v3.0.x and your alert rule configurations include these fields, you must remove them before upgrading to v3.1.0.

**Before:**

```hcl
resource "thousandeyes_alert_rule" "example_alert_rule" {
  rule_name                 = "Example Alert Rule"
  alert_type                = "http-server"
  expression                = "((errorType != \"None\"))"
  minimum_sources           = 2
  rounds_violating_required = 4
  rounds_violating_out_of   = 4

  notifications {
    webhook {
      integration_id   = "wb-12345"
      integration_type = "webhook"
      integration_name = "My Webhook"  # Remove this field
      target           = "https://example.com/webhook"  # Remove this field
    }

    custom_webhook {
      integration_id   = "cw-67890"
      integration_type = "custom-webhook"
      integration_name = "My Custom Webhook"  # Remove this field
      target           = "https://example.com/custom"  # Remove this field
    }
  }
}
```

**After:**

```hcl
resource "thousandeyes_alert_rule" "example_alert_rule" {
  rule_name                 = "Example Alert Rule"
  alert_type                = "http-server"
  expression                = "((errorType != \"None\"))"
  minimum_sources           = 2
  rounds_violating_required = 4
  rounds_violating_out_of   = 4

  notifications {
    webhook {
      integration_id   = "wb-12345"
      integration_type = "webhook"
    }

    custom_webhook {
      integration_id   = "cw-67890"
      integration_type = "custom-webhook"
    }
  }
}
```

The integration name and target URL are managed through the integration configuration in ThousandEyes and will be automatically applied when the alert rule is triggered.

## Other notable changes

- Version 3 of this provider is only compatible with the [ThousandEyes API v7](https://developer.cisco.com/docs/thousandeyes/introduction/). The provider exclusively uses the v7 API, and some fields or behaviors may differ from earlier API versions.
- Error messages and diagnostics have been enhanced and may be stricter than in previous versions.

## Additional resources

- [ThousandEyes API Documentation](https://developer.thousandeyes.com/v7/)
- [ThousandEyes API Migration Guide](https://developer.cisco.com/docs/thousandeyes/migration-guide-overview/)

If you encounter issues during the upgrade, open an issue on GitHub or reach out to us directly.
