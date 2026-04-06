package thousandeyes

import (
	"context"
	"errors"
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
		Update: resourceTagAssignmentUpdate,
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
		return mapTagToBulkTagAssignment(resp), err
	})
}

func resourceTagAssignmentDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagAssignmentAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Tag assignment %s", d.Id())
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

	log.Printf("[INFO] Creating ThousandEyes Tag assignment %s", d.Id())
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

func resourceTagAssignmentUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagAssignmentAPIService)(&apiClient.Common)

	oldTagIDRaw, newTagIDRaw := d.GetChange("tag_id")
	oldAssignmentsRaw, newAssignmentsRaw := d.GetChange("assignments")

	oldTagID, oldLocal := buildTagAssignmentStructFromValues(oldTagIDRaw.(string), oldAssignmentsRaw)
	newTagID, newLocal := buildTagAssignmentStructFromValues(newTagIDRaw.(string), newAssignmentsRaw)

	log.Printf("[INFO] Updating ThousandEyes Tag assignment %s", d.Id())

	assignmentsToAdd := newLocal.Assignments
	assignmentsToRemove := oldLocal.Assignments
	if *oldTagID == *newTagID {
		assignmentsToAdd, assignmentsToRemove = diffTagAssignments(oldLocal.Assignments, newLocal.Assignments)
	}

	if err := applyTagAssignments(apiClient, api, *newTagID, assignmentsToAdd, true); err != nil {
		return err
	}
	if err := applyTagAssignments(apiClient, api, *oldTagID, assignmentsToRemove, false); err != nil {
		if rollbackErr := applyTagAssignments(apiClient, api, *newTagID, assignmentsToAdd, false); rollbackErr != nil {
			log.Printf("[WARN] Rollback failed for ThousandEyes Tag assignment %s: %v", d.Id(), rollbackErr)
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	d.SetId(*newTagID)
	return resourceTagAssignmentRead(d, m)
}

func buildTagAssignmentStruct(d *schema.ResourceData) (*string, *tags.TagAssignment) {
	bulkTagAssignment := ResourceBuildStruct(d, &tags.BulkTagAssignment{})
	return bulkTagAssignment.TagId, &tags.TagAssignment{Assignments: bulkTagAssignment.Assignments}
}

func buildTagAssignmentStructFromValues(tagID string, assignments interface{}) (*string, *tags.TagAssignment) {
	return &tagID, &tags.TagAssignment{
		Assignments: FillValue(assignments, []tags.Assignment{}).([]tags.Assignment),
	}
}

func applyTagAssignments(apiClient *client.APIClient, api *tags.TagAssignmentAPIService, tagID string, assignments []tags.Assignment, assign bool) error {
	if len(assignments) == 0 {
		return nil
	}

	tagAssignment := tags.TagAssignment{Assignments: assignments}
	if assign {
		req := api.AssignTag(tagID).TagAssignment(tagAssignment)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)
		_, _, err := req.Execute()
		return err
	}

	req := api.UnassignTag(tagID).TagAssignment(tagAssignment)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	_, err := req.Execute()
	if err != nil && !IsNotFoundError(err) {
		return err
	}
	return nil
}

func diffTagAssignments(oldAssignments []tags.Assignment, newAssignments []tags.Assignment) ([]tags.Assignment, []tags.Assignment) {
	oldByKey := make(map[string]tags.Assignment, len(oldAssignments))
	newByKey := make(map[string]tags.Assignment, len(newAssignments))

	for _, assignment := range oldAssignments {
		oldByKey[tagAssignmentKey(assignment)] = assignment
	}
	for _, assignment := range newAssignments {
		newByKey[tagAssignmentKey(assignment)] = assignment
	}

	addedAssignments := make([]tags.Assignment, 0)
	for key, assignment := range newByKey {
		if _, exists := oldByKey[key]; !exists {
			addedAssignments = append(addedAssignments, assignment)
		}
	}

	removedAssignments := make([]tags.Assignment, 0)
	for key, assignment := range oldByKey {
		if _, exists := newByKey[key]; !exists {
			removedAssignments = append(removedAssignments, assignment)
		}
	}

	return addedAssignments, removedAssignments
}

func tagAssignmentKey(assignment tags.Assignment) string {
	var id string
	if assignment.Id != nil {
		id = *assignment.Id
	}

	var assignmentType string
	if assignment.Type != nil {
		assignmentType = string(*assignment.Type)
	}

	return assignmentType + "|" + id
}

func mapTagToBulkTagAssignment(in *tags.Tag) *tags.BulkTagAssignment {
	if in == nil {
		return nil
	}
	return &tags.BulkTagAssignment{
		TagId:       in.Id,
		Assignments: in.Assignments,
	}
}
