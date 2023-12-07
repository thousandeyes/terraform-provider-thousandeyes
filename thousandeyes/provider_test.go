package thousandeyes

import (
	"context"
	"fmt"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testClient *thousandeyes.Client

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
	if err != nil {
		fmt.Printf("Error creating provider: %v", err)
		os.Exit(1)
	}

	ctx := context.TODO()
	resourceData := schema.TestResourceDataRaw(t, provider.Schema, map[string]interface{}{})
	testClientRaw, diags := provider.ConfigureContextFunc(ctx, resourceData)
	if diags != nil && diags.HasError() {
		fmt.Printf("Error configuring client: %v", diags)
		os.Exit(1)
	}

	testClient = testClientRaw.(*thousandeyes.Client)
	if testClient == nil {
		fmt.Printf("Error converting client: unexpected type")
		os.Exit(1)
	}
}
