package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestAccThousandEyesWebhookOperation(t *testing.T) {
	var webhookOperationResourceName = "thousandeyes_webhook_operation.test"
	testCases := []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:               "test_webhook_operation_create_update_delete",
			createResourceFile: "acceptance_resources/webhook_operation/basic.tf",
			updateResourceFile: "acceptance_resources/webhook_operation/update.tf",
			checkDestroyFunction: func(state *terraform.State) error {
				resourceList := []ResourceType{
					{
						Name:         "Webhook Operation",
						ResourceName: "thousandeyes_webhook_operation",
						GetResource: func(id string) (interface{}, error) {
							return getWebhookOperation(id)
						}},
				}
				return testAccCheckResourceDestroy(resourceList, state)
			},
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(webhookOperationResourceName, "name", "Test Webhook Operation"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "enabled", "false"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "category", "alerts"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "status", "pending"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "path", "/custom/alerts"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "type", "webhook"),
				resource.TestCheckResourceAttrSet(webhookOperationResourceName, "id"),
				resource.TestCheckResourceAttrSet(webhookOperationResourceName, "link"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(webhookOperationResourceName, "name", "Test Webhook Operation Updated"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "enabled", "true"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "category", "alerts"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "status", "connected"),
				resource.TestCheckResourceAttr(webhookOperationResourceName, "path", "/custom/alerts/v2"),
				resource.TestCheckResourceAttrSet(webhookOperationResourceName, "id"),
				resource.TestCheckResourceAttrSet(webhookOperationResourceName, "link"),
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
						Config: testAccThousandEyesWebhookOperationConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesWebhookOperationConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccThousandEyesWebhookOperationConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getWebhookOperation(id string) (interface{}, error) {
	api := (*connectors.WebhookOperationsAPIService)(&testClient.Common)
	req := api.GetWebhookOperation(id)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
