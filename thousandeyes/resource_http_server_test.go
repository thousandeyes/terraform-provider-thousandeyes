package thousandeyes

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestBuildHTTPServerStructOmitsEmptyPostBodyFromState(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID: "test-id",
		Attributes: map[string]string{
			"test_name": "http server",
			"url":       "https://www.thousandeyes.com",
			"interval":  "120",
			"post_body": "",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody != nil {
		t.Fatalf("expected empty post_body state to be omitted, got PostBody=%q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructOmitsConfiguredEmptyPostBodyWithoutRequestMethod(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.ObjectVal(map[string]cty.Value{"post_body": cty.StringVal("")}),
		Attributes: map[string]string{
			"test_name": "http server",
			"url":       "https://www.thousandeyes.com",
			"interval":  "120",
			"post_body": "",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody != nil {
		t.Fatalf("expected configured empty post_body without request_method to be omitted, got PostBody=%q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructOmitsNullPostBodyWithoutRawRequestMethod(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.ObjectVal(map[string]cty.Value{"post_body": cty.NullVal(cty.String)}),
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "post",
			"post_body":      "stale payload",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody != nil {
		t.Fatalf("expected null post_body without raw request_method to omit stale PostBody, got %q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructRequestMethodGetOmitsPostBody(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.ObjectVal(map[string]cty.Value{"request_method": cty.StringVal("get")}),
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "get",
			"post_body":      "",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody != nil {
		t.Fatalf("expected get request_method to omit PostBody, got %q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructRequestMethodPostDefaultsToEmptyPostBody(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.ObjectVal(map[string]cty.Value{"request_method": cty.StringVal("post")}),
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "post",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody == nil {
		t.Fatal("expected post request_method to send empty PostBody when no body is configured")
	}
	if *req.PostBody != "" {
		t.Fatalf("expected empty PostBody, got %q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructRequestMethodPostOmitsStalePostBody(t *testing.T) {
	for _, tc := range []struct {
		name      string
		rawConfig cty.Value
	}{
		{
			name:      "omitted post_body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{"request_method": cty.StringVal("post")}),
		},
		{
			name: "null post_body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"request_method": cty.StringVal("post"),
				"post_body":      cty.NullVal(cty.String),
			}),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			d := resourceHTTPServer().Data(&terraform.InstanceState{
				ID:        "test-id",
				RawConfig: tc.rawConfig,
				Attributes: map[string]string{
					"test_name":      "http server",
					"url":            "https://www.thousandeyes.com",
					"interval":       "120",
					"request_method": "post",
					"post_body":      "payload",
				},
			})

			req := buildHTTPServerStruct(d)

			if req.PostBody == nil {
				t.Fatal("expected post request_method to send empty PostBody")
			}
			if *req.PostBody != "" {
				t.Fatalf("expected stale PostBody to be cleared, got %q", *req.PostBody)
			}
		})
	}
}

func TestBuildHTTPServerStructOmitsComputedPostWithEmptyPostBody(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID: "test-id",
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "post",
			"post_body":      "",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody != nil {
		t.Fatalf("expected computed post request_method with empty post_body to omit PostBody, got %q", *req.PostBody)
	}
}

func TestBuildHTTPServerStructPreservesConfiguredNonEmptyPostBodyOverComputedGet(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.ObjectVal(map[string]cty.Value{"post_body": cty.StringVal("payload")}),
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "get",
			"post_body":      "payload",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody == nil {
		t.Fatal("expected configured non-empty post_body to be preserved over computed get request_method")
	}
	if *req.PostBody != "payload" {
		t.Fatalf("expected PostBody %q, got %q", "payload", *req.PostBody)
	}
}

func TestHTTPServerRequestMethodGetRejectsConfiguredPostBody(t *testing.T) {
	rawConfig := cty.ObjectVal(map[string]cty.Value{
		"request_method": cty.StringVal(httpServerRequestMethodGET),
		"post_body":      cty.StringVal(""),
	})
	conf := terraform.NewResourceConfigRaw(map[string]interface{}{
		"test_name":      "http server",
		"url":            "https://www.thousandeyes.com",
		"interval":       120,
		"agents":         []interface{}{"1"},
		"request_method": httpServerRequestMethodGET,
		"post_body":      "",
	})

	_, err := resourceHTTPServer().Diff(context.Background(), &terraform.InstanceState{RawConfig: rawConfig}, conf, nil)
	if err == nil {
		t.Fatal("expected get request_method with configured post_body to fail")
	}
	if !strings.Contains(err.Error(), "post_body can only be set when request_method is post") {
		t.Fatalf("expected post_body incompatibility error, got %v", err)
	}
}

func TestResourceHTTPServerReadRequestMethodFallsBackToPostBody(t *testing.T) {
	for _, tc := range []struct {
		name               string
		apiPostBody        string
		expectedHTTPMethod string
	}{
		{
			name:               "nil post body is get",
			expectedHTTPMethod: "get",
		},
		{
			name:               "non-nil post body is post",
			apiPostBody:        `"payload"`,
			expectedHTTPMethod: "post",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			postBody := ""
			if tc.apiPostBody != "" {
				postBody = `,"postBody":` + tc.apiPostBody
			}
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{
  "interval": 120,
  "url": "https://www.thousandeyes.com",
  "testId": "test-id",
  "testName": "http server"` + postBody + `
}`))
			}))
			defer server.Close()

			d := resourceHTTPServer().Data(&terraform.InstanceState{
				ID: "test-id",
			})
			apiClient := client.NewAPIClient(&client.Configuration{
				AuthToken:  "test-token",
				ServerURL:  server.URL,
				HTTPClient: server.Client(),
				Context:    context.Background(),
			})

			if err := resourceHTTPServerRead(d, apiClient); err != nil {
				t.Fatalf("resourceHTTPServerRead returned error: %v", err)
			}

			if got := d.Get(httpServerRequestMethodField); got != tc.expectedHTTPMethod {
				t.Fatalf("expected request_method %s, got %v", tc.expectedHTTPMethod, got)
			}
		})
	}
}

func TestResourceHTTPServerReadRequestMethodStoresAPIValue(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
  "interval": 120,
  "url": "https://www.thousandeyes.com",
  "testId": "test-id",
  "testName": "http server",
  "requestMethod": "post",
  "postBody": ""
}`))
	}))
	defer server.Close()

	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID: "test-id",
	})
	apiClient := client.NewAPIClient(&client.Configuration{
		AuthToken:  "test-token",
		ServerURL:  server.URL,
		HTTPClient: server.Client(),
		Context:    context.Background(),
	})

	if err := resourceHTTPServerRead(d, apiClient); err != nil {
		t.Fatalf("resourceHTTPServerRead returned error: %v", err)
	}

	if got := d.Get(httpServerRequestMethodField); got != "post" {
		t.Fatalf("expected request_method post, got %v", got)
	}
}

func TestBuildHTTPServerStructPreservesNonEmptyPostBody(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID: "test-id",
		Attributes: map[string]string{
			"test_name": "http server",
			"url":       "https://www.thousandeyes.com",
			"interval":  "120",
			"post_body": "payload",
		},
	})

	req := buildHTTPServerStruct(d)

	if req.PostBody == nil {
		t.Fatal("expected non-empty post_body to be preserved")
	}
	if *req.PostBody != "payload" {
		t.Fatalf("expected PostBody %q, got %q", "payload", *req.PostBody)
	}
}

func TestAccThousandEyesHTTPServer(t *testing.T) {
	var httpResourceName = "thousandeyes_http_server.test"
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
			name:                 "create_update_delete_http_server_test",
			createResourceFile:   "acceptance_resources/http_server/basic.tf",
			updateResourceFile:   "acceptance_resources/http_server/update.tf",
			resourceName:         httpResourceName,
			checkDestroyFunction: testAccCheckDefaultHTTPResourceDestroy,
			checkCreateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - HTTP"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "120"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
			},
			checkUpdateFunc: []resource.TestCheckFunc{
				resource.TestCheckResourceAttr(httpResourceName, "url", "https://www.thousandeyes.com"),
				resource.TestCheckResourceAttr(httpResourceName, "test_name", "User Acceptance Test - HTTP (Updated)"),
				resource.TestCheckResourceAttr(httpResourceName, "interval", "300"),
				resource.TestCheckResourceAttr(httpResourceName, "alerts_enabled", "true"),
				resource.TestCheckResourceAttr(httpResourceName, "alert_rules.#", "2"),
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
						Config: testAccThousandEyesHTTPServerConfig(tc.createResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkCreateFunc...),
					},
					{
						Config: testAccThousandEyesHTTPServerConfig(tc.updateResourceFile),
						Check:  resource.ComposeTestCheckFunc(tc.checkUpdateFunc...),
					},
				},
			})
		})
	}
}

func testAccCheckDefaultHTTPResourceDestroy(s *terraform.State) error {
	resourceList := []ResourceType{
		{
			ResourceName: "thousandeyes_http_server",
			GetResource: func(id string) (interface{}, error) {
				return getHttpServer(id)
			}},
	}
	return testAccCheckResourceDestroy(resourceList, s)
}

func testAccThousandEyesHTTPServerConfig(testResource string) string {
	content, err := os.ReadFile(testResource)
	if err != nil {
		panic(err)
	}
	return string(content)
}

func getHttpServer(id string) (interface{}, error) {
	api := (*tests.HTTPServerTestsAPIService)(&testClient.Common)
	req := api.GetHttpServerTest(id).Expand(knownExpandTestOptions())
	req = SetAidFromContext(testClient.GetConfig().Context, req)
	resp, _, err := req.Execute()
	return resp, err
}
