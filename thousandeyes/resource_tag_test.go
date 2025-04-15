package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func TestAccThousandEyesTag(t *testing.T) {
	var tag1ResourceName = "thousandeyes_tag.tag1"
	var tag2ResourceName = "thousandeyes_tag.tag2"
	var tag3ResourceName = "thousandeyes_tag.tag3"
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
				resource.TestCheckResourceAttr(tag1ResourceName, "key", "UAT Tag1 Key"),
				resource.TestCheckResourceAttr(tag1ResourceName, "value", "UAT Tag1 Value"),
				resource.TestCheckResourceAttr(tag1ResourceName, "color", "#b3de69"),
				resource.TestCheckResourceAttr(tag1ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag1ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag1ResourceName, "icon", "LABEL"),
				resource.TestCheckResourceAttr(tag2ResourceName, "key", "UAT Tag2 Key"),
				resource.TestCheckResourceAttr(tag2ResourceName, "value", "UAT Tag2 Value"),
				resource.TestCheckResourceAttr(tag2ResourceName, "color", "#fdb462"),
				resource.TestCheckResourceAttr(tag2ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag2ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag2ResourceName, "icon", "LABEL"),
				resource.TestCheckResourceAttr(tag3ResourceName, "key", "UAT Tag3 Key"),
				resource.TestCheckResourceAttr(tag3ResourceName, "value", "UAT Tag3 Value"),
				resource.TestCheckResourceAttr(tag3ResourceName, "color", "#8dd3c7"),
				resource.TestCheckResourceAttr(tag3ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag3ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag3ResourceName, "icon", "LABEL"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(tag1ResourceName, "key", "UAT Tag1 Key (Updated)"),
				resource.TestCheckResourceAttr(tag1ResourceName, "value", "UAT Tag1 Value (Updated)"),
				resource.TestCheckResourceAttr(tag1ResourceName, "color", "#fdb462"),
				resource.TestCheckResourceAttr(tag1ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag1ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag1ResourceName, "icon", "LABEL"),
				resource.TestCheckResourceAttr(tag2ResourceName, "key", "UAT Tag2 Key (Updated)"),
				resource.TestCheckResourceAttr(tag2ResourceName, "value", "UAT Tag2 Value (Updated)"),
				resource.TestCheckResourceAttr(tag2ResourceName, "color", "#8dd3c7"),
				resource.TestCheckResourceAttr(tag2ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag2ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag2ResourceName, "icon", "LABEL"),
				resource.TestCheckResourceAttr(tag3ResourceName, "key", "UAT Tag3 Key (Updated)"),
				resource.TestCheckResourceAttr(tag3ResourceName, "value", "UAT Tag3 Value (Updated)"),
				resource.TestCheckResourceAttr(tag3ResourceName, "color", "#b3de69"),
				resource.TestCheckResourceAttr(tag3ResourceName, "object_type", "test"),
				resource.TestCheckResourceAttr(tag3ResourceName, "access_type", "all"),
				resource.TestCheckResourceAttr(tag3ResourceName, "icon", "LABEL"),
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

func getTag(id string) (interface{}, error) {
	api := (*tags.TagsAPIService)(&testClient.Common)
	req := api.GetTag(id).Expand(tags.AllowedExpandTagsOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
