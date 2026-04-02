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

	testCases := []struct {
		name               string
		updateResourceFile string
		createCheckFunc    []resource.TestCheckFunc
		updateCheckFunc    []resource.TestCheckFunc
	}{
		{
			name:               "update_tag_assignment_tag_id",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_tag_id.tf",
			createCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag1ID,
				&httpTestID,
				1,
				0,
			),
			updateCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag2ID,
				&httpTestID,
				0,
				1,
			),
		},
		{
			name:               "update_tag_assignment_assignments",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_assignments.tf",
			createCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag1ID,
				&httpTestID,
				1,
				0,
			),
			updateCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag1ID,
				&agentToServerTestID,
				1,
				0,
			),
		},
		{
			name:               "update_tag_assignment_tag_id_and_assignments",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_both.tf",
			createCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag1ID,
				&httpTestID,
				1,
				0,
			),
			updateCheckFunc: testAccCheckTagAssignmentUpdateState(
				&httpTestID,
				&agentToServerTestID,
				&tag1ID,
				&tag2ID,
				&tag2ID,
				&agentToServerTestID,
				0,
				1,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckTagAssignmentUpdateDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccThousandEyesTagConfig("acceptance_resources/tag_assignment_update/basic.tf"),
						Check:  resource.ComposeTestCheckFunc(tc.createCheckFunc...),
					},
					{
						Config: testAccThousandEyesTagConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.updateCheckFunc...),
					},
				},
			})
		})
	}
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
