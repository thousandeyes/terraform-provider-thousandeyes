package schemas

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

// Structs used for mapping
var _ = alerts.Link{}
var _ = tests.TestSelfLink{}

var link = &schema.Schema{
	Type:        schema.TypeSet,
	Description: "A hyperlink from the containing resource to a URI.",
	Optional:    true,
	Elem: &schema.Resource{
		Schema: map[string]*schema.Schema{
			"href": {
				Type:        schema.TypeString,
				Description: "Its value is either a URI [RFC3986] or a URI template [RFC6570].",
				Required:    true,
			},
			"templated": {
				Type:        schema.TypeString,
				Description: "Should be true when the link object's \"href\" property is a URI template.",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Used as a hint to indicate the media type expected when dereferencing the target resource.",
				Required:    true,
			},
			"deprecation": {
				Type:        schema.TypeString,
				Description: "Its presence indicates that the link is to be deprecated at a future date. Its value is a URL that should provide further information about the deprecation.",
				Required:    true,
			},
			"name": {
				Type:        schema.TypeString,
				Description: "Its value may be used as a secondary key for selecting link objects that share the same relation type.",
				Required:    true,
			},
			"profile": {
				Type:        schema.TypeString,
				Description: "A URI that hints about the profile of the target resource.",
				Required:    true,
			},
			"title": {
				Type:        schema.TypeString,
				Description: "Intended for labelling the link with a human-readable identifier",
				Required:    true,
			},
			"hreflang": {
				Type:        schema.TypeString,
				Description: "Indicates the language of the target resource",
				Required:    true,
			},
		},
	},
}
