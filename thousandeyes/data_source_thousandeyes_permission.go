package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func dataSourceThousandeyesPermission() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesPermissionRead,

		Schema: map[string]*schema.Schema{
			"permission_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Permission name.",
			},
			"permission_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique ID representing the permission.",
			},
		},
		Description: "",
	}
}

func dataSourceThousandeyesPermissionRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*administrative.PermissionsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading Thousandeyes account group")

	searchName := d.Get("permission_name").(string)

	resp, _, err := api.GetPermissions().Execute()
	if err != nil {
		return err
	}

	var found *administrative.Permission

	for _, permission := range resp.GetPermissions() {
		if *permission.Permission == searchName {
			found = &permission
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any permission with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found Permission ID: %s - name: %s", *found.PermissionId, *found.Permission)

	d.SetId(*found.PermissionId)
	err = d.Set("permission_name", found.Permission)
	if err != nil {
		return err
	}
	err = d.Set("permission_id", found.PermissionId)
	if err != nil {
		return err
	}

	return nil
}
