package schemas

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func LegacyTestSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"agents": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"alert_rules": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"rule_id": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"bgp_monitors": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"shared_with_accounts": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"aid": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"dns_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_name": {
							Type: schema.TypeString,
						},
						"server_id": {
							Type: schema.TypeInt,
						},
					},
				},
			},
			"custom_headers": {
				Type:        schema.TypeMap,
				Description: "The custom headers.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func LegacyTestStateUpgrade(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {
	// This conditional is required because the schema version was introduced on `v3.0.2`.
	// That means the provider will try to upgrade the state for all versions before that, `v3.0.0|v3.0.1` included.
	// However, all v3 schemas comply with v7 API, so there is no need to upgrade the state for those versions.
	// The `link` field was introduced in `v3.0.0` and is present in all test schemas.
	if _, ok := rawState["link"].(string); !ok {

		if agents, ok := rawState["agents"].([]interface{}); ok {
			for i, v := range agents {
				agent := v.(map[string]interface{})
				agents[i] = agent["agent_id"]
			}
		}

		if alertRules, ok := rawState["alert_rules"].([]interface{}); ok {
			for i, v := range alertRules {
				alertRule := v.(map[string]interface{})
				alertRules[i] = alertRule["rule_id"]
			}
		}

		if bgpMonitors, ok := rawState["bgp_monitors"].([]interface{}); ok {
			for i, v := range bgpMonitors {
				monitor := v.(map[string]interface{})
				bgpMonitors[i] = monitor["monitor_id"]
			}
			rawState["monitors"] = bgpMonitors
			rawState["bgp_monitors"] = nil
		}

		// Only to maintain the backward compatibility
		if groups, ok := rawState["groups"].([]interface{}); ok {
			rawState["labels"] = make([]interface{}, len(groups))
			for i, v := range groups {
				group := v.(map[string]interface{})
				groups[i] = group["group_id"]
			}
			rawState["labels"] = groups
			rawState["groups"] = nil
		}

		if sharedWithAccounts, ok := rawState["shared_with_accounts"].([]interface{}); ok {
			for i, v := range sharedWithAccounts {
				account := v.(map[string]interface{})
				sharedWithAccounts[i] = account["aid"]
			}
		}

		if dnsSevers, ok := rawState["dns_servers"].([]interface{}); ok {
			for i, v := range dnsSevers {
				dnsServer := v.(map[string]interface{})
				dnsSevers[i] = dnsServer["server_name"]
			}
		}

		if _, ok := rawState["custom_headers"].(map[string]interface{}); ok {
			rawState["custom_headers"] = nil
		}
	}

	return rawState, nil
}
