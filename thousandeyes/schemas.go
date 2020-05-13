package thousandeyes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var schemas = map[string]*schema.Schema{
	"agents": {
		Type:        schema.TypeList,
		Description: "agents to use ",
		Required:    true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"agent_id": {
					Type:        schema.TypeInt,
					Description: "agent id",
					Optional:    true,
				},
			},
		},
	},
	"bgp_monitors": {
		Type:        schema.TypeList,
		Optional:    true,
		Description: "array of BGP Monitor objects",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"monitor_id": {
					Type:        schema.TypeInt,
					Description: "monitor id",
					Optional:    true,
				},
			},
		},
	},
	"domain": {
		Type:        schema.TypeString,
		Description: "target record for test, followed by record type (ie, www.thousandeyes.com A)",
		Required:    true,
	},
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "name of the test",
	},

	"include_covered_prefixes": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to include queries for subprefixes detected under this prefix",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
	"interval": {
		Type:         schema.TypeInt,
		Required:     true,
		Description:  "interval to run test on, in seconds",
		ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1800, 3600}),
	},
	"prefix": {
		Type:        schema.TypeString,
		Description: "BGP network address prefix",
		Required:    true,
		// a.b.c.d is a network address, with the prefix length defined as e.
		// Prefixes can be any length from 8 to 24
		// Can only use private BGP monitors for a local prefix.
	},
	"use_public_bgp": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to automatically add all available Public BGP Monitors",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
	},
}
