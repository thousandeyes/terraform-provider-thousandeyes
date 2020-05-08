package thousandeyes

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var schemas = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		Description: "name of the test",
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
	"include_covered_prefixes": {
		Type:         schema.TypeInt,
		Description:  "set to 1 to include queries for subprefixes detected under this prefix",
		Optional:     true,
		ValidateFunc: validation.IntBetween(0, 1),
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
