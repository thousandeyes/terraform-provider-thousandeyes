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

func TestAccThousandEyesConnector_authBasic(t *testing.T) {
	resourceName := "thousandeyes_connector.test_auth"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckConnectorDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Basic"
  target = "%s/auth-basic"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Bearer"
  target = "%s/auth-bearer"

  authentication {
    type  = "bearer-token"
    token = "test-bearer-token-value"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Other Token"
  target = "%s/auth-other"

  authentication {
    type  = "other-token"
    token = "test-other-token-value"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth OAuth Client Creds"
  target = "%s/auth-oauth-cc"

  authentication {
    type                = "oauth-client-credentials"
    oauth_client_id     = "test-client-id"
    oauth_client_secret = "test-client-secret"
    oauth_token_url     = "https://auth.example.com/oauth/token"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth OAuth Code"
  target = "%s/auth-oauth-code"

  authentication {
    type                = "oauth-auth-code"
    oauth_client_id     = "test-client-id"
    oauth_client_secret = "test-client-secret"
    oauth_token_url     = "https://auth.example.com/oauth/token"
    oauth_auth_url      = "https://auth.example.com/oauth/authorize"
    code                = "test-auth-code"
    redirect_uri        = "https://app.example.com/callback"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Switch"
  target = "%s/auth-switch"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
`, connectorTargetBase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Auth Switch"
  target = "%s/auth-switch"

  authentication {
    type  = "bearer-token"
    token = "new-bearer-token"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_full" {
  name   = "UAT - Connector Full"
  target = "%s/full"

  headers {
    name  = "X-Correlation-ID"
    value = "abc-123"
  }

  authentication {
    type     = "basic"
    username = "admin"
    password = "secret"
  }
}
`, connectorTargetBase),
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
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Remove Auth"
  target = "%s/remove-auth"

  authentication {
    type     = "basic"
    username = "testuser"
    password = "testpass"
  }
}
`, connectorTargetBase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authentication.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "authentication.0.type", "basic"),
				),
			},
			{
				Config: fmt.Sprintf(`
resource "thousandeyes_connector" "test_auth" {
  name   = "UAT - Connector Remove Auth"
  target = "%s/remove-auth"
}
`, connectorTargetBase),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "authentication.#", "0"),
				),
			},
		},
	})
}
