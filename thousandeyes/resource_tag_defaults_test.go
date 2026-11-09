package thousandeyes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesTagDefaults(t *testing.T) {
	const resourceName = "thousandeyes_tag.defaults"
	config := testAccThousandEyesTagConfig("acceptance_resources/tags/defaults.tf")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			return testAccCheckResourceDestroy([]ResourceType{
				{
					ResourceName: "thousandeyes_tag",
					GetResource: func(id string) (interface{}, error) {
						return getTag(id)
					},
				},
			}, s)
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "color"),
					resource.TestCheckResourceAttrSet(resourceName, "icon"),
				),
			},
			{
				Config:             config,
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}
