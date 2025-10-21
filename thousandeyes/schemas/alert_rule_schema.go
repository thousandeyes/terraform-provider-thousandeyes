package schemas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
)

// Structs used for mapping
var _ = alerts.RuleDetailUpdate{}
var _ = alerts.RuleDetail{}
var _ = alerts.AlertNotification{}
var _ = alerts.Link{}

var AlertRuleSchema = map[string]*schema.Schema{
	// ruleId
	"rule_id": {
		Type:        schema.TypeString,
		Description: "The unique ID of the alert rule.",
		Computed:    true,
	},
	// ruleName
	"rule_name": {
		Type:        schema.TypeString,
		Description: "The name of the alert rule.",
		Required:    true,
	},
	// expression
	"expression": {
		Type:         schema.TypeString,
		Description:  "The alert rule evaluation expression.",
		Required:     true,
		ValidateFunc: validation.StringIsNotEmpty,
	},
	// direction
	"direction": {
		Type:         schema.TypeString,
		Description:  "[to-target, from-target, bidirectional] The direction of the test (affects how results are shown).",
		Optional:     true,
		ValidateFunc: validation.StringInSlice([]string{"to-target", "from-target", "bidirectional"}, false),
	},
	// notifyOnClear
	"notify_on_clear": {
		Type:        schema.TypeBool,
		Description: "Set to 'true' to trigger the notification when the alert clears.",
		Optional:    true,
		Default:     true,
	},
	// isDefault
	"is_default": {
		Type:        schema.TypeBool,
		Description: "If set to `true`, this alert rule becomes the default for its test type and is automatically applied to newly created tests with relevant metrics. Only one default alert rule is allowed per test type.",
		Optional:    true,
	},
	// alertType
	"alert_type": {
		Description: "The type of alert rule.",
		Type:        schema.TypeString,
		Required:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"page-load",
			"http-server",
			"end-to-end-server",
			"end-to-end-agent",
			"voice",
			"dns-server",
			"dns-trace",
			"dnssec",
			"bgp",
			"path-trace",
			"ftp",
			"sip-server",
			"transactions",
			"web-transactions",
			"agent",
			"network-outage",
			"application-outage",
			"device-device",
			"device-interface",
			"endpoint-network-server",
			"endpoint-http-server",
			"endpoint-path-trace",
			"endpoint-browser-sessions-agent",
			"endpoint-browser-sessions-application",
			"api",
			"web-transaction",
			"unknown",
		}, false),
	},
	// alertGroupType
	"alert_group_type": {
		Description: "The type of alert group.",
		Type:        schema.TypeString,
		Optional:    true,
		Computed:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"bgp",
			"browser-session",
			"cloud-enterprise",
			"endpoint",
		}, false),
	},
	// minimumSources
	"minimum_sources": {
		Type:         schema.TypeInt,
		Description:  "The minimum number of agents or monitors that must meet the specified criteria in order to trigger an alert. This option is mutually exclusive with 'minimum_sources_pct'.",
		Optional:     true,
		ValidateFunc: validation.IntAtLeast(1),
	},
	// minimumSourcesPct
	"minimum_sources_pct": {
		Type:         schema.TypeInt,
		Description:  "The minimum percentage of agents or monitors that must meet the specified criteria in order to trigger an alert. This option is mutually exclusive with 'minimum_sources'.",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 100),
	},
	// roundsViolatingMode
	"rounds_violating_mode": {
		Type:        schema.TypeString,
		Description: "[any, auto or exact] Defines whether the same agent(s) must meet the 'exact' same threshold in consecutive rounds or not. The default value is 'any'.",
		Default:     "any",
		Optional:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"exact",
			"any",
			"auto",
		}, false),
	},
	// roundsViolatingOutOf
	"rounds_violating_out_of": {
		Type:         schema.TypeInt,
		Description:  "Specifies the divisor (Y value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10.",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	// roundsViolatingRequired
	"rounds_violating_required": {
		Type:         schema.TypeInt,
		Description:  "Specifies the numerator (X value) of the “X of Y times” condition in an alert rule.  Minimum value is 1, maximum value is 10. Must be less than or equal to 'roundsViolatingOutOf'.",
		Required:     true,
		ValidateFunc: validation.IntBetween(1, 10),
	},
	// includeCoveredPrefixes
	"include_covered_prefixes": {
		Type:        schema.TypeBool,
		Description: "Include queries for subprefixes detected under this prefix.",
		Optional:    true,
	},
	// sensitivityLevel
	"sensitivity_level": {
		Description: "[high, medium or low] Defines sensitivity level.",
		Type:        schema.TypeString,
		Required:    false,
		Optional:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"high",
			"medium",
			"low",
		}, false),
	},
	// severity
	"severity": {
		Type:        schema.TypeString,
		Description: "[info, minor, major, critical or unknown] The severity level of the alert rule. The default value is 'info'.",
		Default:     "info",
		Optional:    true,
		ValidateFunc: validation.StringInSlice([]string{
			"info",
			"major",
			"minor",
			"critical",
			"unknown",
		}, false),
	},
	// endpointAgentIds
	"endpoint_agent_ids": {
		Type:        schema.TypeSet,
		Description: "An array of endpoint agent IDs associated with the rule (get `id` from `/endpoint/agents` API). This is applicable when `alertGroupType` is `browser-session`.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// endpointLabelIds
	"endpoint_label_ids": {
		Type:        schema.TypeSet,
		Description: "An array of label IDs used to assign specific Endpoint Agents to the test (get `id` from `/endpoint/labels`). This is applicable when `alertGroupType` is `browser-session`.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// visitedSitesFilter
	"visited_sites_filter": {
		Type:        schema.TypeSet,
		Description: "A list of website domains visited during the session. This is applicable when `alertGroupType` is `browser-session`.",
		Optional:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// notifications
	"notifications": {
		Type:        schema.TypeSet,
		Description: "The list of notifications for the alert rule.",
		Optional:    true,
		Set: func(v interface{}) int {
			// Custom hash function that only considers meaningful fields for comparison
			// This prevents drift when read-only fields like integration_name or target differ
			var buf strings.Builder
			m := v.(map[string]interface{})

			// Hash email notifications
			if email, ok := m["email"].(*schema.Set); ok && email != nil {
				for _, e := range email.List() {
					buf.WriteString(fmt.Sprintf("email:%v|", e))
				}
			}

			// Hash third_party notifications (only integration_id and integration_type)
			if tp, ok := m["third_party"].(*schema.Set); ok && tp != nil {
				for _, t := range tp.List() {
					if tpMap, ok := t.(map[string]interface{}); ok {
						buf.WriteString(fmt.Sprintf("tp:%s:%s|", tpMap["integration_id"], tpMap["integration_type"]))
					}
				}
			}

			// Hash webhook notifications (only integration_id and integration_type)
			if webhook, ok := m["webhook"].(*schema.Set); ok && webhook != nil {
				for _, w := range webhook.List() {
					if wMap, ok := w.(map[string]interface{}); ok {
						buf.WriteString(fmt.Sprintf("webhook:%s:%s|", wMap["integration_id"], wMap["integration_type"]))
					}
				}
			}

			// Hash custom_webhook notifications (only integration_id and integration_type)
			if cw, ok := m["custom_webhook"].(*schema.Set); ok && cw != nil {
				for _, w := range cw.List() {
					if wMap, ok := w.(map[string]interface{}); ok {
						buf.WriteString(fmt.Sprintf("custom_webhook:%s:%s|", wMap["integration_id"], wMap["integration_type"]))
					}
				}
			}

			return schema.HashString(buf.String())
		},
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
							"recipients": {
								Type:        schema.TypeSet,
								Description: "The email addresses to send the notification to.",
								Optional:    true,
								Elem: &schema.Schema{
									Type: schema.TypeString,
								},
							},
						},
					},
				},
				"third_party": {
					Type:        schema.TypeSet,
					Description: "Third party notification.",
					Optional:    true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"integration_id": {
								Type:        schema.TypeString,
								Description: "The integration ID, as a string.",
								Required:    true,
							},
							"integration_type": {
								Type:        schema.TypeString,
								Description: "The integration type, as a string.",
								Required:    true,
							},
						},
					},
				},
				"webhook": {
					Type:        schema.TypeSet,
					Description: "Webhook notification.",
					Optional:    true,
					Set: func(v interface{}) int {
						// Custom hash function that only considers integration_id and integration_type
						// This prevents drift when integration_name or target differ between config and state
						m := v.(map[string]interface{})
						integrationID := ""
						integrationType := ""
						if id, ok := m["integration_id"]; ok && id != nil {
							integrationID = id.(string)
						}
						if typ, ok := m["integration_type"]; ok && typ != nil {
							integrationType = typ.(string)
						}
						return schema.HashString(integrationID + "|" + integrationType)
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"integration_id": {
								Type:        schema.TypeString,
								Description: "The integration ID, as a string.",
								Required:    true,
							},
							"integration_type": {
								Type:        schema.TypeString,
								Description: "The integration type, as a string.",
								Required:    true,
							},
							"integration_name": {
								Type:        schema.TypeString,
								Description: "Name of the integration.",
								Optional:    true,
								Computed:    true,
								Deprecated:  "This field is informational only and cannot be set through the alert rule. It reflects the integration name from the integration configuration. Please remove this field from your configuration to avoid perpetual drift.",
							},
							"target": {
								Type:        schema.TypeString,
								Description: "Webhook target URL.",
								Optional:    true,
								Computed:    true,
								Deprecated:  "This field is informational only and cannot be set through the alert rule. It reflects the target URL from the integration configuration. Please remove this field from your configuration to avoid perpetual drift.",
							},
						},
					},
				},
				"custom_webhook": {
					Type:        schema.TypeSet,
					Description: "Webhook notification.",
					Optional:    true,
					Set: func(v interface{}) int {
						// Custom hash function that only considers integration_id and integration_type
						// This prevents drift when integration_name or target differ between config and state
						m := v.(map[string]interface{})
						integrationID := ""
						integrationType := ""
						if id, ok := m["integration_id"]; ok && id != nil {
							integrationID = id.(string)
						}
						if typ, ok := m["integration_type"]; ok && typ != nil {
							integrationType = typ.(string)
						}
						return schema.HashString(integrationID + "|" + integrationType)
					},
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"integration_id": {
								Type:        schema.TypeString,
								Description: "The integration ID, as a string.",
								Required:    true,
							},
							"integration_type": {
								Type:        schema.TypeString,
								Description: "The integration type, as a string.",
								Required:    true,
							},
							"integration_name": {
								Type:        schema.TypeString,
								Description: "Name of the integration.",
								Optional:    true,
								Computed:    true,
								Deprecated:  "This field is informational only and cannot be set through the alert rule. It reflects the integration name from the integration configuration. Please remove this field from your configuration to avoid perpetual drift.",
							},
							"target": {
								Type:        schema.TypeString,
								Description: "Webhook target URL.",
								Optional:    true,
								Computed:    true,
								Deprecated:  "This field is informational only and cannot be set through the alert rule. It reflects the target URL from the integration configuration. Please remove this field from your configuration to avoid perpetual drift.",
							},
						},
					},
				},
			},
		},
	},
	// testIds (or ids from "tests")
	"test_ids": {
		Type:        schema.TypeSet,
		Description: "The valid test IDs.",
		Computed:    true,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
	},
	// link (_links.self.href)
	"link": {
		Type:        schema.TypeString,
		Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
		Computed:    true,
	},
}
