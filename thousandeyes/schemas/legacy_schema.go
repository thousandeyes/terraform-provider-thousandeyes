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
		},
	}
}

func LegacyTestStateUpgrade(ctx context.Context, rawState map[string]any, meta any) (map[string]any, error) {

	for i, v := range rawState["agents"].([]interface{}) {
		agent := v.(map[string]interface{})
		rawState["agents"].([]interface{})[i] = agent["agent_id"]
	}

	if _, ok := rawState["alert_rules"].([]interface{}); ok {
		for i, v := range rawState["alert_rules"].([]interface{}) {
			alertRule := v.(map[string]interface{})
			rawState["alert_rules"].([]interface{})[i] = alertRule["rule_id"]
		}
	}

	if bgpMonitors, ok := rawState["bgp_monitors"].([]interface{}); ok {
		for i, v := range rawState["bgp_monitors"].([]interface{}) {
			monitor := v.(map[string]interface{})
			bgpMonitors[i] = monitor["monitor_id"]
		}
		rawState["monitors"] = bgpMonitors
		rawState["bgp_monitors"] = nil
	}

	// Only to maintain the backward compatibility
	if groups, ok := rawState["groups"].([]interface{}); ok {
		rawState["labels"] = make([]interface{}, len(groups))
		for i, v := range rawState["groups"].([]interface{}) {
			group := v.(map[string]interface{})
			groups[i] = group["group_id"]
		}
		rawState["labels"] = groups
		rawState["groups"] = nil
	}

	if sharedWithAccounts, ok := rawState["shared_with_accounts"].([]interface{}); ok {
		for i, v := range rawState["shared_with_accounts"].([]interface{}) {
			account := v.(map[string]interface{})
			sharedWithAccounts[i] = account["aid"]
		}
	}

	if dnsSevers, ok := rawState["dns_servers"].([]interface{}); ok {
		for i, v := range rawState["dns_servers"].([]interface{}) {
			dnsServer := v.(map[string]interface{})
			dnsSevers[i] = dnsServer["server_name"]
		}
	}

	return rawState, nil
}
