package thousandeyes

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func TestAccThousandEyesTag(t *testing.T) {
	var httpTestResourceName = "thousandeyes_http_server.test"
	var agentToServerTestResourceName = "thousandeyes_agent_to_server.test"
	var tag1ResourceName = "thousandeyes_tag.tag1"
	var tag2ResourceName = "thousandeyes_tag.tag2"
	var tag3ResourceName = "thousandeyes_tag.tag3"
	var assign1ResourceName = "thousandeyes_tag_assignment.assign1"
	var assign2ResourceName = "thousandeyes_tag_assignment.assign2"
	var assign3ResourceName = "thousandeyes_tag_assignment.assign3"
	var httpTestId string
	var agentToServerTestId string
	var tag1Id string
	var tag2Id string
	var tag3Id string

	var testCases = []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:                 "create_update_delete_tags_and_tags_assignments_test",
			createResourceFile:   "acceptance_resources/tags/basic.tf",
			updateResourceFile:   "acceptance_resources/tags/update.tf",
			checkDestroyFunction: testAccCheckDefaultTagsDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				testAccCheckResourceExistsAndStoreID(httpTestResourceName, &httpTestId),
				testAccCheckResourceExistsAndStoreID(agentToServerTestResourceName, &agentToServerTestId),
				testAccCheckResourceExistsAndStoreID(tag1ResourceName, &tag1Id),
				resource.TestCheckResourceAttr(tag1ResourceName, "key", "A UAT Tag1 Key"),
				resource.TestCheckResourceAttr(tag1ResourceName, "value", "A UAT Tag1 Value"),
				resource.TestCheckResourceAttr(tag1ResourceName, "color", "#b3de69"),
				resource.TestCheckResourceAttr(tag1ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag1ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag1ResourceName, "icon", "LABEL"),
				testAccCheckResourceExistsAndStoreID(tag2ResourceName, &tag2Id),
				resource.TestCheckResourceAttr(tag2ResourceName, "key", "A UAT Tag2 Key"),
				resource.TestCheckResourceAttr(tag2ResourceName, "value", "A UAT Tag2 Value"),
				resource.TestCheckResourceAttr(tag2ResourceName, "color", "#fdb462"),
				resource.TestCheckResourceAttr(tag2ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag2ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag2ResourceName, "icon", "LABEL"),
				testAccCheckResourceExistsAndStoreID(tag3ResourceName, &tag3Id),
				resource.TestCheckResourceAttr(tag3ResourceName, "key", "A UAT Tag3 Key"),
				resource.TestCheckResourceAttr(tag3ResourceName, "value", "A UAT Tag3 Value"),
				resource.TestCheckResourceAttr(tag3ResourceName, "color", "#8dd3c7"),
				resource.TestCheckResourceAttr(tag3ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag3ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag3ResourceName, "icon", "LABEL"),
				testAccCheckDependentResourceHasCorrectID(assign1ResourceName, "tag_id", &tag1Id),
				testAccCheckDependentResourceHasCorrectID(assign1ResourceName, "assignments.0.id", &httpTestId),
				resource.TestCheckResourceAttr(assign1ResourceName, "assignments.0.type", "test"),
				testAccCheckAsignmentsCountInTagsResponse(&tag1Id, 1),
				testAccCheckDependentResourceHasCorrectID(assign2ResourceName, "tag_id", &tag2Id),
				testAccCheckDependentResourceHasCorrectID(assign2ResourceName, "assignments.0.id", &agentToServerTestId),
				resource.TestCheckResourceAttr(assign2ResourceName, "assignments.0.type", "test"),
				testAccCheckAsignmentsCountInTagsResponse(&tag2Id, 1),
				testAccCheckDependentResourceHasCorrectID(assign3ResourceName, "tag_id", &tag3Id),
				resource.TestCheckResourceAttr(assign3ResourceName, "assignments.0.type", "test"),
				resource.TestCheckResourceAttr(assign3ResourceName, "assignments.1.type", "test"),
				testAccCheckAsignmentsCountInTagsResponse(&tag3Id, 2),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				testAccCheckResourceExistsAndStoreID(httpTestResourceName, &httpTestId),
				testAccCheckResourceExistsAndStoreID(agentToServerTestResourceName, &agentToServerTestId),
				testAccCheckResourceExistsAndStoreID(tag1ResourceName, &tag1Id),
				resource.TestCheckResourceAttr(tag1ResourceName, "key", "A UAT Tag1 Key (Updated)"),
				resource.TestCheckResourceAttr(tag1ResourceName, "value", "A UAT Tag1 Value (Updated)"),
				resource.TestCheckResourceAttr(tag1ResourceName, "color", "#fdb462"),
				resource.TestCheckResourceAttr(tag1ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag1ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag1ResourceName, "icon", "LABEL"),
				testAccCheckResourceExistsAndStoreID(tag2ResourceName, &tag2Id),
				resource.TestCheckResourceAttr(tag2ResourceName, "key", "A UAT Tag2 Key (Updated)"),
				resource.TestCheckResourceAttr(tag2ResourceName, "value", "A UAT Tag2 Value (Updated)"),
				resource.TestCheckResourceAttr(tag2ResourceName, "color", "#8dd3c7"),
				resource.TestCheckResourceAttr(tag2ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag2ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag2ResourceName, "icon", "LABEL"),
				testAccCheckResourceExistsAndStoreID(tag3ResourceName, &tag3Id),
				resource.TestCheckResourceAttr(tag3ResourceName, "key", "A UAT Tag3 Key (Updated)"),
				resource.TestCheckResourceAttr(tag3ResourceName, "value", "A UAT Tag3 Value (Updated)"),
				resource.TestCheckResourceAttr(tag3ResourceName, "color", "#b3de69"),
				resource.TestCheckResourceAttr(tag3ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag3ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag3ResourceName, "icon", "LABEL"),
				testAccCheckDependentResourceHasCorrectID(assign1ResourceName, "tag_id", &tag2Id),
				testAccCheckDependentResourceHasCorrectID(assign1ResourceName, "assignments.0.id", &httpTestId),
				resource.TestCheckResourceAttr(assign1ResourceName, "assignments.0.type", "test"),
				testAccCheckAsignmentsCountInTagsResponse(&tag2Id, 1),
				testAccCheckDependentResourceHasCorrectID(assign2ResourceName, "tag_id", &tag1Id),
				testAccCheckDependentResourceHasCorrectID(assign2ResourceName, "assignments.0.id", &agentToServerTestId),
				resource.TestCheckResourceAttr(assign2ResourceName, "assignments.0.type", "test"),
				testAccCheckAsignmentsCountInTagsResponse(&tag1Id, 1),
				testAccCheckAsignmentsCountInTagsResponse(&tag3Id, 0),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      tc.checkDestroyFunction,
				Steps: []resource.TestStep{
					{
						Config: testAccThousandEyesTagConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesTagConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultTagsDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_tag.tag1",
			GetResource: func(id string) (interface{}, error) {
				return getTag(id)
			}},
		{
			ResourceName: "thousandeyes_tag.tag2",
			GetResource: func(id string) (interface{}, error) {
				return getTag(id)
			}},
		{
			ResourceName: "thousandeyes_tag.tag3",
			GetResource: func(id string) (interface{}, error) {
				return getTag(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesTagConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getTag(id string) (*tags.Tag, error) {
	api := (*tags.TagsAPIService)(&testClient.Common)
	req := api.GetTag(id).Expand(tags.AllowedExpandTagsOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}

func testAccCheckResourceExistsAndStoreID(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID is not set")
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckDependentResourceHasCorrectID(resourceName string, attributeName string, expectedID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Resource not found: %s", resourceName)
		}
		actualParentID := rs.Primary.Attributes[attributeName]
		if actualParentID != *expectedID {
			return fmt.Errorf("Expected %s is %s, got %s", attributeName, *expectedID, actualParentID)
		}
		return nil
	}
}

func testAccCheckAsignmentsCountInTagsResponse(expectedTagID *string, count int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		tag, err := getTag(*expectedTagID)
		if err != nil {
			return fmt.Errorf("Unable to get Tag with Id %s", *expectedTagID)
		}
		actualCount := len(tag.Assignments)
		if actualCount != count {
			return fmt.Errorf("Expected asignments count is %d, got %d", count, actualCount)
		}
		return nil
	}
}
