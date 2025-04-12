package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func resourceTag() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tags.Tag{}, schemas.TagSchema, nil),
		Create: resourceTagCreate,
		Read:   resourceTagRead,
		Update: resourceTagUpdate,
		Delete: resourceTagDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource provides a tagging system with key/value pairs. It allows you to tag assets within the ThousandEyes platform (such as agents, tests, or alert rules) with meaningful metadata.",
	}
	return &resource
}

func resourceTagRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.WithValue(context.Background(), tagsKey, struct{}{})

	return GetResource(ctx, d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tags.TagsAPIService)(&apiClient.Common)

		req := api.GetTag(id).Expand(tags.AllowedExpandTagsOptionsEnumValues)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()

		// set nullable fields
		if resp.Icon.IsSet() {
			d.Set("icon", resp.Icon.Get())
		}
		if resp.Description.IsSet() {
			d.Set("description", resp.Description.Get())
		}
		if resp.LegacyId.IsSet() {
			d.Set("legacy_id", resp.LegacyId.Get())
		}

		return resp, err
	})
}

func resourceTagUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildTagStruct(d)

	req := api.UpdateTag(d.Id()).TagInfo(*update)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceTagRead(d, m)
}

func resourceTagDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteTag(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTagCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tags.TagsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildTagStruct(d)

	req := api.CreateTag().TagInfo(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.Id
	d.SetId(id)
	return resourceTagRead(d, m)
}

func buildTagStruct(d *schema.ResourceData) *tags.TagInfo {
	tag := ResourceBuildStruct(d, &tags.TagInfo{})
	// set nullable fields
	if v, ok := d.Get("icon").(string); ok {
		tag.Icon.Set(
			getPointer(v),
		)
	}
	if v, ok := d.Get("description").(string); ok {
		tag.Description.Set(
			getPointer(v),
		)
	}
	if v, ok := d.Get("legacy_id").(int); ok {
		tag.LegacyId.Set(
			getPointer(float32(v)),
		)
	}
	return tag
}
