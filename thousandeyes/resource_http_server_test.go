package thousandeyes

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

const httpServerRequestMethodGET = "get"
const httpServerRequestMethodPOST = "post"

func TestBuildHTTPServerStructSetsExplicitRequestMethodAndPostBody(t *testing.T) {
	for _, tc := range []struct {
		name               string
		rawConfig          cty.Value
		stateRequestMethod string
		statePostBody      string
		wantRequestMethod  tests.RequestMethod
		wantPostBody       string
	}{
		{
			name:               "no request_method and omitted post_body defaults GET with empty body",
			rawConfig:          cty.EmptyObjectVal,
			stateRequestMethod: httpServerRequestMethodPOST,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_GET,
			wantPostBody:       "",
		},
		{
			name: "no request_method and null post_body defaults GET with empty body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"post_body": cty.NullVal(cty.String),
			}),
			stateRequestMethod: httpServerRequestMethodPOST,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_GET,
			wantPostBody:       "",
		},
		{
			name: "no request_method and empty post_body defaults GET with empty body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"post_body": cty.StringVal(""),
			}),
			stateRequestMethod: httpServerRequestMethodPOST,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_GET,
			wantPostBody:       "",
		},
		{
			name: "no request_method and non-empty post_body defaults POST with configured body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"post_body": cty.StringVal("payload"),
			}),
			stateRequestMethod: httpServerRequestMethodGET,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_POST,
			wantPostBody:       "payload",
		},
		{
			name: "GET request_method and omitted post_body sends empty body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"request_method": cty.StringVal(httpServerRequestMethodGET),
			}),
			stateRequestMethod: httpServerRequestMethodPOST,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_GET,
			wantPostBody:       "",
		},
		{
			name: "POST request_method and omitted post_body sends empty body",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"request_method": cty.StringVal(httpServerRequestMethodPOST),
			}),
			stateRequestMethod: httpServerRequestMethodGET,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_POST,
			wantPostBody:       "",
		},
		{
			name: "GET request_method and explicit post_body preserve configured values",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"request_method": cty.StringVal(httpServerRequestMethodGET),
				"post_body":      cty.StringVal("payload"),
			}),
			stateRequestMethod: httpServerRequestMethodPOST,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_GET,
			wantPostBody:       "payload",
		},
		{
			name: "POST request_method and explicit post_body preserve configured values",
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"request_method": cty.StringVal(httpServerRequestMethodPOST),
				"post_body":      cty.StringVal("payload"),
			}),
			stateRequestMethod: httpServerRequestMethodGET,
			statePostBody:      "stale payload",
			wantRequestMethod:  tests.REQUESTMETHOD_POST,
			wantPostBody:       "payload",
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
					"request_method": tc.stateRequestMethod,
					"post_body":      tc.statePostBody,
				},
			})

			req := buildHTTPServerStruct(d)

			assertHTTPServerRequestMethod(t, req, tc.wantRequestMethod)
			assertStringPointer(t, req.PostBody, getPointer(tc.wantPostBody))
		})
	}
}

func TestHTTPServerTestRequestSDKMarshalJSONIncludesExplicitMethodAndPostBody(t *testing.T) {
	d := resourceHTTPServer().Data(&terraform.InstanceState{
		ID:        "test-id",
		RawConfig: cty.EmptyObjectVal,
		Attributes: map[string]string{
			"test_name":      "http server",
			"url":            "https://www.thousandeyes.com",
			"interval":       "120",
			"request_method": "post",
			"post_body":      "stale payload",
		},
	})

	req := buildHTTPServerStruct(d)
	body, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal HTTP server request: %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("failed to unmarshal HTTP server request: %v", err)
	}

	if got := payload["requestMethod"]; got != httpServerRequestMethodGET {
		t.Fatalf("expected request method %q, got %q", httpServerRequestMethodGET, got)
	}
	if got := payload["postBody"]; got != "" {
		t.Fatalf("expected postBody %q, got %q", "", got)
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

func TestHTTPServerGeneratedUpdateSendsRequestMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT request, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/tests/http-server/test%2Fid" {
			t.Fatalf("expected escaped update path, got %s", r.URL.EscapedPath())
		}
		if got := r.URL.Query().Get("aid"); got != "aid-123" {
			t.Fatalf("expected aid query param %q, got %q", "aid-123", got)
		}
		if got := r.URL.Query().Get("expand"); got != "agent,monitor" {
			t.Fatalf("expected expand query param %q, got %q", "agent,monitor", got)
		}
		if got := r.Header.Get("Content-Type"); got != "application/json" {
			t.Fatalf("expected Content-Type application/json, got %q", got)
		}
		if got := r.Header.Get("Accept"); got != "application/hal+json,application/json,application/problem+json" {
			t.Fatalf("unexpected Accept header %q", got)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("unexpected Authorization header %q", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if got := payload["requestMethod"]; got != httpServerRequestMethodPOST {
			t.Fatalf("expected requestMethod %q, got %q", httpServerRequestMethodPOST, got)
		}
		if got := payload["postBody"]; got != "payload" {
			t.Fatalf("expected postBody %q, got %q", "payload", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"testId":"test-id","testName":"http server","url":"https://www.thousandeyes.com","interval":120}`))
	}))
	defer server.Close()

	apiClient := client.NewAPIClient(&client.Configuration{
		AuthToken:  "test-token",
		ServerURL:  server.URL,
		HTTPClient: server.Client(),
		Context:    context.Background(),
	})
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)
	req := api.UpdateHttpServerTest("test/id").HttpServerTestRequest(tests.HttpServerTestRequest{
		PostBody:      getPointer("payload"),
		RequestMethod: tests.REQUESTMETHOD_POST.Ptr(),
	}).Aid("aid-123").Expand([]tests.ExpandTestOptions{
		tests.EXPANDTESTOPTIONS_AGENT,
		tests.EXPANDTESTOPTIONS_MONITOR,
	})

	resp, httpResp, err := req.Execute()
	if err != nil {
		t.Fatalf("generated update request returned error: %v", err)
	}
	if httpResp == nil || httpResp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200 response, got %#v", httpResp)
	}
	if resp == nil || resp.TestId == nil || *resp.TestId != "test-id" {
		t.Fatalf("expected decoded test response, got %#v", resp)
	}
}

func TestHTTPServerGeneratedCreateSendsRequestMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("expected POST request, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/tests/http-server" {
			t.Fatalf("expected create path, got %s", r.URL.EscapedPath())
		}
		if got := r.URL.Query().Get("aid"); got != "aid-123" {
			t.Fatalf("expected aid query param %q, got %q", "aid-123", got)
		}
		if got := r.URL.Query().Get("expand"); got != "agent,monitor" {
			t.Fatalf("expected expand query param %q, got %q", "agent,monitor", got)
		}

		var payload map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if got := payload["requestMethod"]; got != httpServerRequestMethodGET {
			t.Fatalf("expected requestMethod %q, got %q", httpServerRequestMethodGET, got)
		}
		if got := payload["postBody"]; got != "" {
			t.Fatalf("expected postBody %q, got %q", "", got)
		}

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"testId":"test-id","testName":"http server","url":"https://www.thousandeyes.com","interval":120}`))
	}))
	defer server.Close()

	apiClient := client.NewAPIClient(&client.Configuration{
		AuthToken:  "test-token",
		ServerURL:  server.URL,
		HTTPClient: server.Client(),
		Context:    context.Background(),
	})
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)
	req := api.CreateHttpServerTest().HttpServerTestRequest(tests.HttpServerTestRequest{
		PostBody:      getPointer(""),
		RequestMethod: tests.REQUESTMETHOD_GET.Ptr(),
	}).Aid("aid-123").Expand([]tests.ExpandTestOptions{
		tests.EXPANDTESTOPTIONS_AGENT,
		tests.EXPANDTESTOPTIONS_MONITOR,
	})

	resp, httpResp, err := req.Execute()
	if err != nil {
		t.Fatalf("generated create request returned error: %v", err)
	}
	if httpResp == nil || httpResp.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200 response, got %#v", httpResp)
	}
	if resp == nil || resp.TestId == nil || *resp.TestId != "test-id" {
		t.Fatalf("expected decoded test response, got %#v", resp)
	}
}

func assertHTTPServerRequestMethod(t *testing.T, req *tests.HttpServerTestRequest, want tests.RequestMethod) {
	t.Helper()

	if req.RequestMethod == nil {
		t.Fatalf("expected SDK RequestMethod %q, got nil", want)
	}
	if got := *req.RequestMethod; got != want {
		t.Fatalf("expected SDK RequestMethod %q, got %q", want, got)
	}
}

func assertStringPointer(t *testing.T, got, want *string) {
	t.Helper()

	if want == nil {
		if got != nil {
			t.Fatalf("expected nil string pointer, got %q", *got)
		}
		return
	}
	if got == nil {
		t.Fatalf("expected string pointer %q, got nil", *want)
	}
	if *got != *want {
		t.Fatalf("expected string pointer %q, got %q", *want, *got)
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
