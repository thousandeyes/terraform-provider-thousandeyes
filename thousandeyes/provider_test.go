package thousandeyes

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/require"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

type ResourceType struct {
	Name         string
	ResourceName string
	GetResource  func(id string) (interface{}, error)
}

var testClient *client.APIClient

var providerFactories = map[string]func() (*schema.Provider, error){
	"thousandeyes": func() (*schema.Provider, error) {
		return New()(), nil
	},
}

func init() {
}

func testAccPreCheck(t *testing.T) {
	providerFunc, _ := providerFactories["thousandeyes"]
	provider, err := providerFunc()
	require.Nil(t, err, "Error creating provider: %v", err)

	ctx := context.TODO()
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{})
	testClientRaw, diags := provider.ConfigureContextFunc(ctx, resourceData)

	require.False(t, diags != nil && diags.HasError(), "Error configuring client: %v", diags)

	testClient = testClientRaw.(*client.APIClient)
	require.NotNil(t, testClient, "Error converting client: unexpected type")
}

func testAccCheckResourceDestroy(resources []ResourceType, s *terraform.State) error {
	for _, resource := range resources {
		for _, rs := range s.RootModule().Resources {
			if rs.Type == resource.ResourceName {
				var err error
				id := rs.Primary.ID
				_, err = resource.GetResource(id)
				if err == nil {
					return fmt.Errorf("%s with id %s still exists", resource.ResourceName, rs.Primary.ID)
				}
			}
		}
	}
	return nil
}
