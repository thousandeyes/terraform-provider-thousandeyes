package thousandeyes

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestAccThousandEyesConnector(t *testing.T) {
	resourceName := "thousandeyes_connector.test"
	var connectorId string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					resource.TestCheckResourceAttr(resourceName, "name", "User Acceptance Test - Connector"),
					resource.TestCheckResourceAttr(resourceName, "target", defaultConnectorTarget),
					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
				),
			},
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/update.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					resource.TestCheckResourceAttr(resourceName, "name", "User Acceptance Test - Connector (Updated)"),
					resource.TestCheckResourceAttr(resourceName, "target", defaultConnectorTargetUpdated),
					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
				),
			},
		},
	})
}

func testAccCheckConnectorDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_connector.test",
			GetResource: func(id string) (interface{}, error) {
				return getConnector(id)
			},
		},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccConnectorConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getConnector(id string) (*connectors.GenericConnector, error) {
	api := (*connectors.GenericConnectorsAPIService)(&testClient.Common)
	req := api.GetGenericConnector(id)
	req = SetAidFromContextAny(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}

func TestAccThousandEyesConnectorImport(t *testing.T) {
	resourceName := "thousandeyes_connector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/basic.tf"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Authentication and headers are sensitive and not returned by API
				ImportStateVerifyIgnore: []string{"authentication", "headers"},
			},
		},
	})
}

func TestAccThousandEyesConnectorWithAuth(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"
	var connectorId string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy: func(s *terraform.State) error {
			resourceList := []ResourceType{
				{
					ResourceName: resourceName,
					GetResource: func(id string) (interface{}, error) {
						return getConnector(id)
					},
				},
			}
			return testAccCheckResourceDestroy(resourceList, s)
		},
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "User Acceptance Test - Connector With Auth"
  target = "%s"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
`, defaultConnectorTargetAuth),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					resource.TestCheckResourceAttr(resourceName, "name", "User Acceptance Test - Connector With Auth"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
				),
			},
		},
	})
}

const (
	defaultConnectorTarget        = "https://webhook.site/6b3c063d-d857-4bb3-92eb-b04b6fc41a85"
	defaultConnectorTargetUpdated = "https://webhook.site/6b3c063d-d857-4bb3-92eb-b04b6fc41a85/updated"
	defaultConnectorTargetAuth    = "https://webhook.site/6b3c063d-d857-4bb3-92eb-b04b6fc41a85/auth"
)
