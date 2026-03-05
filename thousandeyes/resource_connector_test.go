package thousandeyes

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/connectors"
)

const (
	connectorTargetBase = "https://example.com/webhooks/thousandeyes"
)

// --- Helpers ---

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
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}

func testAccCheckConnectorDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_connector",
			GetResource: func(id string) (interface{}, error) {
				return getConnector(id)
			},
		},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccCheckConnectorExists(resourceName string, connector **connectors.GenericConnector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}
		result, err := getConnector(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching connector %s: %w", rs.Primary.ID, err)
		}
		*connector = result
		return nil
	}
}

func testAccCheckConnectorServerValues(connector **connectors.GenericConnector, expectedName, expectedTarget string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := *connector
		if c == nil {
			return fmt.Errorf("connector is nil")
		}
		if c.Name != expectedName {
			return fmt.Errorf("server name = %q, want %q", c.Name, expectedName)
		}
		if c.Target != expectedTarget {
			return fmt.Errorf("server target = %q, want %q", c.Target, expectedTarget)
		}
		if c.Type != connectors.CONNECTORTYPE_GENERIC {
			return fmt.Errorf("server type = %q, want %q", c.Type, connectors.CONNECTORTYPE_GENERIC)
		}
		return nil
	}
}

func testAccCheckConnectorHasHeaders(connector **connectors.GenericConnector, expectedCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := *connector
		if c == nil {
			return fmt.Errorf("connector is nil")
		}
		if len(c.Headers) != expectedCount {
			return fmt.Errorf("server headers count = %d, want %d", len(c.Headers), expectedCount)
		}
		return nil
	}
}

func testAccCheckConnectorHasComputedFields(connector **connectors.GenericConnector) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		c := *connector
		if c == nil {
			return fmt.Errorf("connector is nil")
		}
		if c.Id == nil || *c.Id == "" {
			return fmt.Errorf("server id is empty")
		}
		if c.LastModifiedDate == nil {
			return fmt.Errorf("server lastModifiedDate is nil")
		}
		return nil
	}
}

// --- Tests ---

func TestAccThousandEyesConnector_basic(t *testing.T) {
	resourceName := "thousandeyes_connector.test"
	var connector *connectors.GenericConnector
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
					testAccCheckConnectorExists(resourceName, &connector),
					resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Basic"),
					resource.TestCheckResourceAttr(resourceName, "target", connectorTargetBase),
					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_date"),
					testAccCheckConnectorServerValues(&connector, "UAT - Connector Basic", connectorTargetBase),
					testAccCheckConnectorHasComputedFields(&connector),
				),
			},
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/update.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					testAccCheckConnectorExists(resourceName, &connector),
					resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Basic (Updated)"),
					resource.TestCheckResourceAttr(resourceName, "target", connectorTargetBase+"/updated"),
					resource.TestCheckResourceAttr(resourceName, "type", "generic"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_date"),
					testAccCheckConnectorServerValues(&connector, "UAT - Connector Basic (Updated)", connectorTargetBase+"/updated"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_import(t *testing.T) {
	resourceName := "thousandeyes_connector.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/basic.tf"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Basic"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication", "headers"},
			},
		},
	})
}

func TestAccThousandEyesConnector_withHeaders(t *testing.T) {
	resourceName := "thousandeyes_connector.test_headers"
	var connector *connectors.GenericConnector
	var connectorId string

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/with_headers.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					testAccCheckConnectorExists(resourceName, &connector),
					resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector With Headers"),
					resource.TestCheckResourceAttr(resourceName, "headers.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "headers.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr(resourceName, "headers.1.name", "X-Another-Header"),
					testAccCheckConnectorHasHeaders(&connector, 2),
				),
			},
			{
				Config: testAccConnectorConfig("acceptance_resources/connector/with_headers_update.tf"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceExistsAndStoreID(resourceName, &connectorId),
					testAccCheckConnectorExists(resourceName, &connector),
					resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector With Headers (Updated)"),
					resource.TestCheckResourceAttr(resourceName, "headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "headers.0.name", "X-Updated-Header"),
					testAccCheckConnectorHasHeaders(&connector, 1),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_withAuthentication(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/with_authentication.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Auth Basic"),
						resource.TestCheckResourceAttr(resourceName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.username", "testuser"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_authBearerToken(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_bearer_token.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Auth Bearer"),
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "bearer-token"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_authOtherToken(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_other_token.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Auth Other Token"),
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "other-token"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_authOAuthClientCredentials(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_oauth_client_credentials.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Auth OAuth Client Creds"),
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "oauth-client-credentials"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.oauth_client_id", "test-client-id"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.oauth_token_url", "https://auth.example.com/oauth/token"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_authOAuthCode(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_oauth_code.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Auth OAuth Code"),
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "oauth-auth-code"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.oauth_client_id", "test-client-id"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.oauth_token_url", "https://auth.example.com/oauth/token"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.oauth_auth_url", "https://auth.example.com/oauth/authorize"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.redirect_uri", "https://app.example.com/callback"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func TestAccThousandEyesConnector_authSwitch(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_switch_basic.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
					),
				},
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/auth_switch_bearer.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "bearer-token"),
					),
			},
		},
	})
}

func TestAccThousandEyesConnector_headersAndAuth(t *testing.T) {
	resourceName := "thousandeyes_connector.test_full"
	var connector *connectors.GenericConnector

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/headers_and_auth.tf"),
					Check: resource.ComposeTestCheckFunc(
						testAccCheckConnectorExists(resourceName, &connector),
						resource.TestCheckResourceAttr(resourceName, "name", "UAT - Connector Full"),
					resource.TestCheckResourceAttr(resourceName, "target", connectorTargetBase+"/full"),
					resource.TestCheckResourceAttr(resourceName, "headers.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "headers.0.name", "X-Correlation-ID"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.username", "admin"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "last_modified_date"),
					testAccCheckConnectorHasHeaders(&connector, 1),
					testAccCheckConnectorHasComputedFields(&connector),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authentication", "headers"},
			},
		},
	})
}

func TestAccThousandEyesConnector_removeAuth(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
			Steps: []resource.TestStep{
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/remove_auth_with_auth.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "authentication.#", "1"),
						resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
					),
				},
				{
					Config: testAccConnectorConfig("acceptance_resources/connector/remove_auth_without_auth.tf"),
					Check: resource.ComposeTestCheckFunc(
						resource.TestCheckResourceAttr(resourceName, "authentication.#", "0"),
					),
			},
		},
	})
}
