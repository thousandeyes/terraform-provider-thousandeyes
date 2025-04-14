package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func resourceTagAssignment() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tags.BulkTagAssignment{}, schemas.TagAssignmentSchema, nil),
		Create: resourceTagAssignmentCreate,
		Read:   resourceTagAssignmentRead,
		Delete: resourceTagAssignmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource provides a tagging system with key/value pairs. It allows you to tag assets within the ThousandEyes platform (such as agents, tests, or alert rules) with meaningful metadata.",
	}
	return &resource
}

func resourceTagAssignmentRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tags.TagsAPIService)(&apiClient.Common)

		req := api.GetTag(id).Expand(tags.AllowedExpandTagsOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		if err != nil {
			return resp, err
		}

		return mapTagToBulkTagAssignment(resp), err
	})
}

func resourceTagAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagAssignmentAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, local := buildTagAssignmentStruct(d)

	req := api.UnassignTag(*id).TagAssignment(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTagAssignmentCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagAssignmentAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	id, local := buildTagAssignmentStruct(d)

	req := api.AssignTag(*id).TagAssignment(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	d.SetId(*resp.TagId)
	return resourceTagAssignmentRead(d, m)
}

func buildTagAssignmentStruct(d *schema.ResourceData) (*string, *tags.TagAssignment) {
	bulkTagAssignment := ResourceBuildStruct(d, &tags.BulkTagAssignment{})
	return bulkTagAssignment.TagId, &tags.TagAssignment{Assignments: bulkTagAssignment.Assignments}
}

func mapTagToBulkTagAssignment(in *tags.Tag) *tags.BulkTagAssignment {
	return &tags.BulkTagAssignment{
		TagId:       in.Id,
		Assignments: in.Assignments,
	}
}
