package thousandeyes

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	tagAssignmentUpdateBasicConfig = "acceptance_resources/tag_assignment_update/basic.tf"
	httpTestResourceName           = "thousandeyes_http_server.test"
	agentToServerTestResourceName  = "thousandeyes_agent_to_server.test"
	tag1ResourceName               = "thousandeyes_tag.tag1"
	tag2ResourceName               = "thousandeyes_tag.tag2"
	assignResourceName             = "thousandeyes_tag_assignment.assign1"
)

type tagAssignmentState struct {
	expectedTagID        *string
	expectedAssignmentID *string
	tag1Assignments      int
	tag2Assignments      int
}

func TestAccThousandEyesTagAssignmentUpdate(t *testing.T) {
	var httpTestID string
	var agentToServerTestID string
	var tag1ID string
	var tag2ID string

	checkState := func(state tagAssignmentState) []resource.TestCheckFunc {
		return []resource.TestCheckFunc{
			testAccCheckResourceExistsAndStoreID(httpTestResourceName, &httpTestID),
			testAccCheckResourceExistsAndStoreID(agentToServerTestResourceName, &agentToServerTestID),
			testAccCheckResourceExistsAndStoreID(tag1ResourceName, &tag1ID),
			testAccCheckResourceExistsAndStoreID(tag2ResourceName, &tag2ID),
			testAccCheckDependentResourceHasCorrectID(assignResourceName, "tag_id", state.expectedTagID),
			testAccCheckDependentResourceHasCorrectID(assignResourceName, "assignments.0.id", state.expectedAssignmentID),
			resource.TestCheckResourceAttr(assignResourceName, "assignments.0.type", "test"),
			testAccCheckAsignmentsCountInTagsResponse(&tag1ID, state.tag1Assignments),
			testAccCheckAsignmentsCountInTagsResponse(&tag2ID, state.tag2Assignments),
		}
	}

	var testCases = []struct {
		name               string
		updateResourceFile string
		createState        tagAssignmentState
		updateState        tagAssignmentState
	}{
		{
			name:               "update_tag_assignment_tag_id",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_tag_id.tf",
			createState: tagAssignmentState{
				expectedTagID:        &tag1ID,
				expectedAssignmentID: &httpTestID,
				tag1Assignments:      1,
				tag2Assignments:      0,
			},
			updateState: tagAssignmentState{
				expectedTagID:        &tag2ID,
				expectedAssignmentID: &httpTestID,
				tag1Assignments:      0,
				tag2Assignments:      1,
			},
		},
		{
			name:               "update_tag_assignment_assignments",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_assignments.tf",
			createState: tagAssignmentState{
				expectedTagID:        &tag1ID,
				expectedAssignmentID: &httpTestID,
				tag1Assignments:      1,
				tag2Assignments:      0,
			},
			updateState: tagAssignmentState{
				expectedTagID:        &tag1ID,
				expectedAssignmentID: &agentToServerTestID,
				tag1Assignments:      1,
				tag2Assignments:      0,
			},
		},
		{
			name:               "update_tag_assignment_tag_id_and_assignments",
			updateResourceFile: "acceptance_resources/tag_assignment_update/update_both.tf",
			createState: tagAssignmentState{
				expectedTagID:        &tag1ID,
				expectedAssignmentID: &httpTestID,
				tag1Assignments:      1,
				tag2Assignments:      0,
			},
			updateState: tagAssignmentState{
				expectedTagID:        &tag2ID,
				expectedAssignmentID: &agentToServerTestID,
				tag1Assignments:      0,
				tag2Assignments:      1,
			},
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
						Config: testAccThousandEyesTagConfig(tagAssignmentUpdateBasicConfig),
						Check:  resource.ComposeTestCheckFunc(checkState(tc.createState)...),
					},
					{
						Config: testAccThousandEyesTagConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(checkState(tc.updateState)...),
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
