package thousandeyes

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

func TestAccThousandEyesConnector(t *testing.T) {
	resourceName := "thousandeyes_connector.test"
	var connectorId string

	testCases := []struct {
		name               string
		createResourceFile string
		updateResourceFile string
		checkCreateFunc    []resource.TestCheckFunc
		checkUpdateFunc    []resource.TestCheckFunc
	}{
		{
			name:               "basic_connector",
			createResourceFile: "acceptance_resources/connector/basic.tf",
			updateResourceFile: "acceptance_resources/connector/update.tf",
			checkCreateFunc: []resource.TestCheckFunc{
				testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
				resource.TestCheckResourceAttr(resourceName, "name", "UAT Test Connector"),
				resource.TestCheckResourceAttr(resourceName, "target", "https://webhook.site/6b3c063d-d857-4bb3-92eb-b04b6fc41a85"),
				resource.TestCheckResourceAttr(resourceName, "type", "generic"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
				resource.TestCheckResourceAttr(resourceName, "name", "UAT Test Connector (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "target", "https://webhook.site/6b3c063d-d857-4bb3-92eb-b04b6fc41a85/updated"),
				resource.TestCheckResourceAttr(resourceName, "type", "generic"),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:          func() { testAccPreCheck(t) },
				ProviderFactories: providerFactories,
				CheckDestroy:      testAccCheckConnectorDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccConnectorConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccConnectorConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
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
	config := string(content)
	config = strings.ReplaceAll(config, "__CONNECTOR_TARGET__", connectorTargetEnv("TE_CONNECTOR_TARGET", defaultConnectorTarget))
	config = strings.ReplaceAll(config, "__CONNECTOR_TARGET_UPDATED__", connectorTargetEnv("TE_CONNECTOR_TARGET_UPDATED", defaultConnectorTargetUpdated))
	return config
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
	authTarget := connectorTargetEnv("TE_CONNECTOR_TARGET_AUTH", defaultConnectorTargetAuth)

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
  name   = "UAT Connector with Auth"
  target = "%s"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
`, authTarget),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					resource.TestCheckResourceAttr(resourceName, "name", "UAT Connector with Auth"),
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

func connectorTargetEnv(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
