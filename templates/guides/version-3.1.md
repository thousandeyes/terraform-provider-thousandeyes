---
page_title: "ThousandEyes Terraform Provider v3.1 Upgrade Guide"
subcategory: ""
---

# ThousandEyes Terraform Provider Version 3.1 Upgrade Guide

Version 3.1.0 of the ThousandEyes Terraform provider removes fields that were unintentionally introduced in v3.0.x. This guide is specifically for users upgrading from v3.0.x to v3.1.0 or later.

> **Note:** If you are upgrading directly from v2.x to v3.1.0+, you will not encounter this issue as these fields never existed in v2.x. Please refer to the [Version 3.0 Upgrade Guide](version-3.md) for changes related to the v2.x to v3.x upgrade.

## Alert Rule webhook notification changes (Regression Fix)

The `integration_name` and `target` fields have been removed from webhook and custom webhook notifications in alert rules. These fields were unintentionally introduced in v3.0.x and were read-only, returned by the ThousandEyes API based on the integration configuration. They could not be set through alert rules and caused configuration drift.

**Breaking change in v3.1.0:** If you are using v3.0.x and your alert rule configurations include these fields, you must remove them before upgrading to v3.1.0.

### Before (v3.0.x):

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

### After (v3.1.0+):

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

## Additional resources

- [ThousandEyes API Documentation](https://developer.thousandeyes.com/v7/)
- [Version 3.0 Upgrade Guide](version-3.md)

If you encounter issues during the upgrade, open an issue on GitHub or reach out to us directly.

