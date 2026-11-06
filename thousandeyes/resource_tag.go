package thousandeyes

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tags"
)

func resourceTag() *schema.Resource {
	resource := schema.Resource{
		Schema:        ResourceSchemaBuild(tags.Tag{}, schemas.TagSchema, nil),
		Create:        resourceTagCreate,
		Read:          resourceTagRead,
		Update:        resourceTagUpdate,
		Delete:        resourceTagDelete,
		CustomizeDiff: resourceTagCustomizeDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description: "This resource manages key/value tags for ThousandEyes assets, including static tags for tests and dashboards and dynamic tags that automatically select Endpoint Agents using filter criteria.",
	}
	return &resource
}

func resourceTagRead(d *schema.ResourceData, m interface{}) error {
	ctx := context.WithValue(context.Background(), tagsKey, struct{}{})

	return GetResource(ctx, d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*tags.TagsAPIService)(&apiClient.Common)

		req := api.GetTag(id)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		resp, _, err := req.Execute()
		if err != nil {
			return resp, err
		}

		// The API may omit the default scope and absent matching logic.
		// Normalize these before generic mapping so refresh and import agree.
		for i := range resp.Filters {
			if resp.Filters[i].Scope == nil {
				resp.Filters[i].SetScope(tags.TAGFILTERSCOPE_DEFAULT)
			}
		}
		if resp.MatchType == nil {
			if err := d.Set("match_type", nil); err != nil {
				return resp, err
			}
		}

		// set nullable fields
		if resp.Icon.IsSet() {
			if err := d.Set("icon", resp.Icon.Get()); err != nil {
				return resp, err
			}
		}
		if resp.Description.IsSet() {
			if err := d.Set("description", resp.Description.Get()); err != nil {
				return resp, err
			}
		}
		if resp.LegacyId.IsSet() {
			if err := d.Set("legacy_id", resp.LegacyId.Get()); err != nil {
				return resp, err
			}
		}
		if resp.ModifiedDate.IsSet() {
			if modifiedAt := resp.ModifiedDate.Get(); modifiedAt != nil {
				if err := d.Set("modified_date", modifiedAt.UTC().Format(time.RFC3339)); err != nil {
					return resp, err
				}
			} else {
				if err := d.Set("modified_date", nil); err != nil {
					return resp, err
				}
			}
		}

		return resp, err
	})
}

func resourceTagUpdate(d *schema.ResourceData, m interface{}) error {
	if err := validateTagConfiguration(d, func(string) bool { return true }); err != nil {
		return err
	}
	apiClient := m.(*client.APIClient)
	api := (*tags.TagsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Tag %s", d.Id())
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

	log.Printf("[INFO] Deleting ThousandEyes Tag %s", d.Id())

	req := api.DeleteTag(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTagCreate(d *schema.ResourceData, m interface{}) error {
	if err := validateTagConfiguration(d, func(string) bool { return true }); err != nil {
		return err
	}
	apiClient := m.(*client.APIClient)
	api := (*tags.TagsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Tag %s", d.Id())
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
	// Empty matching logic is absent on static tags, never an empty API enum.
	if d.Get("match_type").(string) == "" {
		tag.MatchType = nil
	}
	if len(tag.Filters) == 0 {
		tag.Filters = nil
	}
	// set nullable fields
	if v := d.Get("icon").(string); len(v) != 0 {
		tag.Icon.Set(
			getPointer(v),
		)
	}
	if v := d.Get("description").(string); len(v) != 0 {
		tag.Description.Set(
			getPointer(v),
		)
	}
	return tag
}

func resourceTagCustomizeDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	known := func(key string) bool {
		config := d.GetRawConfig()
		if !config.IsNull() && config.IsKnown() && config.Type().HasAttribute(key) {
			value := config.GetAttr(key)
			if !value.IsKnown() {
				return false
			}
			// Distinguish an omitted optional/computed type from an unknown
			// configured expression in a new plan.
			if key == "type" && value.IsNull() {
				return true
			}
		}
		if key == "filters" {
			return d.NewValueKnown("filters.#")
		}
		return d.NewValueKnown(key)
	}
	return validateTagConfiguration(d, known)
}

func validateTagConfiguration(d interface{ Get(string) interface{} }, known func(string) bool) error {
	objectType := d.Get("object_type").(string)
	tagType := d.Get("type").(string)
	matchType := d.Get("match_type").(string)
	filterCount := d.Get("filters").(*schema.Set).Len()

	if known("object_type") && objectType == "endpoint-agent" && known("type") && tagType != "dynamic" {
		return fmt.Errorf("endpoint-agent tags require type = dynamic")
	}
	if known("type") && tagType == "dynamic" {
		if known("object_type") && objectType != "endpoint-agent" {
			return fmt.Errorf("dynamic tags require object_type = endpoint-agent")
		}
		if known("match_type") && matchType == "" {
			return fmt.Errorf("dynamic tags require match_type")
		}
		if known("filters") && filterCount == 0 {
			return fmt.Errorf("dynamic tags require at least one filter")
		}
	}
	if (known("object_type") && objectType != "endpoint-agent") || (known("type") && tagType != "dynamic") {
		if (known("match_type") && matchType != "") || (known("filters") && filterCount > 0) {
			return fmt.Errorf("match_type and filters are only supported for dynamic endpoint-agent tags")
		}
	}
	return nil
}
