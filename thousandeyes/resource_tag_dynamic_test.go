package thousandeyes

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesTagDynamicEndpointAgent(t *testing.T) {
	const name = "thousandeyes_tag.endpoint_agents"
	key := acctest.RandomWithPrefix("tf-acc-endpoint")
	var tagID string
	initial := testAccDynamicTagConfig(key, "and", "Windows")
	updated := testAccDynamicTagConfig(key, "or", "Linux")
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(*terraform.State) error {
			if tagID == "" {
				return nil
			}
			_, err := getTag(tagID)
			if err == nil {
				return fmt.Errorf("dynamic tag %s still exists", tagID)
			}
			if !IsNotFoundError(err) {
				return err
			}
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: initial,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(name, &tagID),
					resource.TestCheckResourceAttr(name, "object_type", "endpoint-agent"),
					resource.TestCheckResourceAttr(name, "type", "dynamic"),
					resource.TestCheckResourceAttr(name, "match_type", "and"),
					resource.TestCheckResourceAttr(name, "filters.#", "1"),
					resource.TestCheckTypeSetElemAttr(name, "filters.*.values.*", "Windows"),
					resource.TestCheckTypeSetElemNestedAttrs(name, "filters.*", map[string]string{
						"key": "platform", "mode": "in", "scope": "default", "values.#": "1",
					}),
				),
			},
			{Config: initial, PlanOnly: true, ExpectNonEmptyPlan: false},
			{
				Config: updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDependentResourceHasCorrectID(name, "id", &tagID),
					resource.TestCheckResourceAttr(name, "match_type", "or"),
					resource.TestCheckTypeSetElemNestedAttrs(name, "filters.*", map[string]string{
						"key": "platform", "mode": "in", "scope": "default", "values.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(name, "filters.*.values.*", "Linux"),
				),
			},
			{Config: updated, PlanOnly: true, ExpectNonEmptyPlan: false},
			{ResourceName: name, ImportState: true, ImportStateVerify: true},
		},
	})
}

func testAccDynamicTagConfig(key, matchType, platform string) string {
	return fmt.Sprintf(`
resource "thousandeyes_tag" "endpoint_agents" {
  key         = %q
  value       = "selected"
  access_type = "all"
  object_type = "endpoint-agent"
  type        = "dynamic"
  match_type  = %q

  filters {
    key    = "platform"
    values = [%q]
    mode   = "in"
  }
}
`, key, matchType, platform)
}
