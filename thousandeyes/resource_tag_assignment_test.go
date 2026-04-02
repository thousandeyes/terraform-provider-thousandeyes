package thousandeyes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccThousandEyesTagAssignmentUpdate(t *testing.T) {
	var httpTestID string
	var agentToServerTestID string
	var tag1ID string
	var tag2ID string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckTagAssignmentUpdateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccThousandEyesTagConfig("acceptance_resources/tag_assignment_update/basic.tf"),
				Check: resource.ComposeTestCheckFunc(testAccCheckTagAssignmentUpdateState(
					&httpTestID,
					&agentToServerTestID,
					&tag1ID,
					&tag2ID,
					&tag1ID,
					&httpTestID,
					1,
					0,
				)...),
			},
			{
				Config: testAccThousandEyesTagConfig("acceptance_resources/tag_assignment_update/update_assignments.tf"),
				Check: resource.ComposeTestCheckFunc(testAccCheckTagAssignmentUpdateState(
					&httpTestID,
					&agentToServerTestID,
					&tag1ID,
					&tag2ID,
					&tag1ID,
					&agentToServerTestID,
					1,
					0,
				)...),
			},
		},
	})
}

func testAccCheckTagAssignmentUpdateDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_tag.tag1",
			GetResource: func(id string) (interface{}, error) {
				return getTag(id)
			},
		},
		{
			ResourceName: "thousandeyes_tag.tag2",
			GetResource: func(id string) (interface{}, error) {
				return getTag(id)
			},
		},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccCheckTagAssignmentUpdateState(httpTestID *string, agentToServerTestID *string, tag1ID *string, tag2ID *string, expectedTagID *string, expectedAssignmentID *string, tag1Assignments int, tag2Assignments int) []resource.TestCheckFunc {
	return []resource.TestCheckFunc{
		testAccCheckResourceExistsAndStoreID("thousandeyes_http_server.test", httpTestID),
		testAccCheckResourceExistsAndStoreID("thousandeyes_agent_to_server.test", agentToServerTestID),
		testAccCheckResourceExistsAndStoreID("thousandeyes_tag.tag1", tag1ID),
		testAccCheckResourceExistsAndStoreID("thousandeyes_tag.tag2", tag2ID),
		testAccCheckDependentResourceHasCorrectID("thousandeyes_tag_assignment.assign1", "tag_id", expectedTagID),
		testAccCheckDependentResourceHasCorrectID("thousandeyes_tag_assignment.assign1", "assignments.0.id", expectedAssignmentID),
		resource.TestCheckResourceAttr("thousandeyes_tag_assignment.assign1", "assignments.0.type", "test"),
		testAccCheckAsignmentsCountInTagsResponse(tag1ID, tag1Assignments),
		testAccCheckAsignmentsCountInTagsResponse(tag2ID, tag2Assignments),
	}
}
