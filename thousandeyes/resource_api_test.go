package thousandeyes

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestAccThousandEyesAPI(t *testing.T) {
	var resourceName = "thousandeyes_api.test"
	var testCases = []struct {
		name                 string
		createResourceFile   string
		updateResourceFile   string
		resourceName         string
		checkDestroyFunction func(*terraform.State) error
		checkCreateFunc      []resource.TestCheckFunc
		checkUpdateFunc      []resource.TestCheckFunc
	}{
		{
			name:                 "create_update_delete_api_test",
			createResourceFile:   "acceptance_resources/api/basic.tf",
			updateResourceFile:   "acceptance_resources/api/update.tf",
			resourceName:         resourceName,
			checkDestroyFunction: testAccCheckDefaultAPIResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - API"),
				resource.TestCheckResourceAttr(resourceName, "interval", "120"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.name", "Step 1 - GET Request"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.url", "https://api.stg.thousandeyes.com/v6/status.json"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.client_authentication", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.0.username", "test_username"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.password", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.0.method", "get"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.auth_type", "basic"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.0.key", "Accept"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.0.value", "application/json"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.name", "response-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.operator", "includes"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.value", "timestamp"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.name", "status-code"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.operator", "is"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.value", "200"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.name", "Step 2 - POST OAuth2 request"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.url", "https://api.stg.thousandeyes.com/v6/credentials/new.json"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_authentication", "in-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_id", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_secret", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.1.scope", "test_scope"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.token_url", "https://api.stg.thousandeyes.com/v6"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.method", "post"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.auth_type", "oauth2"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.body", "{\"firstName\":\"John\",\"lastName\":\"Doe\"}"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.0.key", "Content-Type"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.0.value", "application/json"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.name", "response-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.operator", "includes"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.value", "error"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.name", "status-code"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.operator", "is"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.value", "401"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(resourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(resourceName, "test_name", "User Acceptance Test - API (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "interval", "300"),
				resource.TestCheckResourceAttr(resourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(resourceName, "alert_rules.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.name", "Step 1 - GET Request (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.url", "https://api.stg.thousandeyes.com/v6/status.json"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.client_authentication", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.0.username", "new_test_username"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.password", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.0.method", "get"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.auth_type", "basic"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.0.key", "Accept"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.headers.0.value", "application/json"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.name", "response-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.operator", "includes"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.0.value", "timestamp"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.name", "status-code"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.operator", "is"),
				resource.TestCheckResourceAttr(resourceName, "requests.0.assertions.1.value", "200"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.name", "Step 2 - POST OAuth2 request (Updated)"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.url", "https://api.stg.thousandeyes.com/v6/credentials/new.json"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_authentication", "in-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_id", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.1.client_secret", ""),
				resource.TestCheckResourceAttr(resourceName, "requests.1.scope", "test_scope"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.token_url", "https://api.stg.thousandeyes.com/v6"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.method", "post"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.auth_type", "oauth2"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.body", "{\"firstName\":\"Jack\",\"lastName\":\"Doe\"}"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.0.key", "Content-Type"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.headers.0.value", "application/json"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.#", "2"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.name", "response-body"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.operator", "includes"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.0.value", "error"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.name", "status-code"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.operator", "is"),
				resource.TestCheckResourceAttr(resourceName, "requests.1.assertions.1.value", "401"),
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
						Config:             testAccThousandEyesAPIConfig(tc.createResourceFile),
						Check:              resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
						ExpectNonEmptyPlan: true,
					},
					{
						Config:             testAccThousandEyesAPIConfig(tc.updateResourceFile),
						Check:              resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
						ExpectNonEmptyPlan: true,
					},
				},
			})
		})
	}
}

func testAccCheckDefaultAPIResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_api",
			GetResource: func(id string) (interface{}, error) {
				return getAPI(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesAPIConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getAPI(id string) (interface{}, error) {
	api := (*tests.APITestsAPIService)(&testClient.Common)
	req := api.GetApiTest(id).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
