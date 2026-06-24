package thousandeyes

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

const httpHeaderSourceModeField = "header_source_mode"
const httpHeaderSourceModeHeaders = "headers"
const httpHeaderSourceModeCustomHeaders = "custom_headers"
const httpServerRequestMethodField = "request_method"
const httpServerRequestMethodGET = "get"
const httpServerRequestMethodPOST = "post"

func resourceHTTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(tests.HttpServerTestRequest{}, schemas.CommonSchema, nil),
		Create: resourceHTTPServerCreate,
		Read:   resourceHTTPServerRead,
		Update: resourceHTTPServerUpdate,
		Delete: resourceHTTPServerDelete,
		CustomizeDiff: customdiff.All(
			normalizeHTTPServerHeadersDiff,
			validateHTTPServerRequestMethodDiff,
		),
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create an HTTP server test. This test type measures the availability and performance of an HTTP service. For more information, see [HTTP Server Tests](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests#http-server-test).",
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    schemas.LegacyTestSchema().CoreConfigSchema().ImpliedType(),
				Upgrade: schemas.LegacyTestStateUpgrade,
				Version: 0,
			},
		},
		SchemaVersion: 1,
	}
	resource.Schema["oauth"] = schemas.CommonSchema["oauth"]
	resource.Schema[httpHeaderSourceModeField] = &schema.Schema{
		Type:     schema.TypeString,
		Computed: true,
	}
	resource.Schema[httpServerRequestMethodField] = schemas.CommonSchema[httpServerRequestMethodField]
	return &resource
}

func resourceHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	log.Printf("[INFO] Reading Thousandeyes Resource %s", d.Id())

	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)
	req := api.GetHttpServerTest(d.Id()).Expand(knownExpandTestOptions())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil && IsNotFoundError(err) {
		log.Printf("[INFO] Resource was deleted - will recreate it")
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	existingOAuth := currentHTTPServerOAuthStateValue(d)
	if err := ResourceRead(context.Background(), d, resp); err != nil {
		return err
	}
	if err := setHTTPServerRequestMethodState(d, resp); err != nil {
		return err
	}

	if rawConfigOAuthConfigured(d) || len(existingOAuth) > 0 {
		if err := d.Set("oauth", terraformHTTPServerOAuthStateValue(d, existingOAuth, resp.OAuth)); err != nil {
			return err
		}
	} else if err := d.Set("oauth", nil); err != nil {
		return err
	}

	mode := httpHeaderSourceMode(d)
	if mode == httpHeaderSourceModeCustomHeaders {
		// Keep only custom_headers in state when it is the configured source of truth.
		if err := d.Set("headers", nil); err != nil {
			return err
		}
		if err := d.Set("custom_headers", terraformHTTPServerCustomHeadersValue(resp.CustomHeaders)); err != nil {
			return err
		}
	} else {
		apiHeaders, ok := normalizeStringInterfaceSlice(resp.Headers)
		if ok {
			if err := d.Set("headers", apiHeaders); err != nil {
				return err
			}
		} else if err := d.Set("headers", nil); err != nil {
			return err
		}
		if err := d.Set("custom_headers", []interface{}{}); err != nil {
			return err
		}
	}

	if err := d.Set(httpHeaderSourceModeField, mode); err != nil {
		return err
	}

	return nil
}

func resourceHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	update := buildHTTPServerStruct(d)

	req := api.UpdateHttpServerTest(d.Id()).HttpServerTestRequest(*update).Expand(knownExpandTestOptions())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}
	return resourceHTTPServerRead(d, m)
}

func resourceHTTPServerDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteHttpServerTest(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceHTTPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildHTTPServerStruct(d)

	req := api.CreateHttpServerTest().HttpServerTestRequest(*local).Expand(knownExpandTestOptions())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.TestId
	d.SetId(id)
	return resourceHTTPServerRead(d, m)
}

func buildHTTPServerStruct(d *schema.ResourceData) *tests.HttpServerTestRequest {
	req := ResourceBuildStruct(d, &tests.HttpServerTestRequest{})
	headers, headersConfigured := rawConfigHeaderStrings(d)
	customHeaders, customHeadersConfigured := rawConfigCustomHeaders(d)

	if requestMethod, configured := rawConfigHTTPServerRequestMethod(d); configured {
		switch requestMethod {
		case httpServerRequestMethodGET:
			req.PostBody = nil
		case httpServerRequestMethodPOST:
			if !rawConfigPostBodyConfigured(d) || req.PostBody == nil {
				empty := ""
				req.PostBody = &empty
			}
		}
	} else {
		req.PostBody = rawConfigHTTPServerPostBodyWithoutRequestMethod(d)
	}

	if headersConfigured {
		req.Headers = headers
		req.CustomHeaders = nil
	} else if customHeadersConfigured {
		req.Headers = nil
		req.CustomHeaders = customHeaders
	} else {
		req.Headers = nil
		req.CustomHeaders = nil
	}
	return req
}

func rawConfigHTTPServerPostBodyWithoutRequestMethod(d rawConfigReader) *string {
	if rawPostBody, configured := rawConfigHTTPServerPostBody(d); configured {
		if rawPostBody != "" {
			return &rawPostBody
		}
	}
	return nil
}

func rawConfigHTTPServerPostBody(d rawConfigReader) (string, bool) {
	raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: "post_body"}})
	if diags.HasError() || !raw.IsKnown() || raw.IsNull() {
		return "", false
	}
	return raw.AsString(), true
}

func rawConfigPostBodyConfigured(d rawConfigReader) bool {
	_, configured := rawConfigHTTPServerPostBody(d)
	return configured
}

func rawConfigHTTPServerRequestMethod(d rawConfigReader) (string, bool) {
	raw, diags := d.GetRawConfigAt(cty.Path{cty.GetAttrStep{Name: httpServerRequestMethodField}})
	if diags.HasError() || !raw.IsKnown() || raw.IsNull() {
		return "", false
	}
	return normalizeHTTPServerRequestMethod(raw.AsString())
}

func validateHTTPServerRequestMethodDiff(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
	method, ok := rawConfigHTTPServerRequestMethod(d)
	if ok && method == httpServerRequestMethodGET && rawConfigPostBodyConfigured(d) {
		return fmt.Errorf("post_body can only be set when request_method is post")
	}
	return nil
}

func setHTTPServerRequestMethodState(d *schema.ResourceData, resp *tests.HttpServerTestResponse) error {
	if method, ok := httpServerResponseRequestMethod(resp); ok {
		return d.Set(httpServerRequestMethodField, method)
	}
	if resp.PostBody != nil {
		return d.Set(httpServerRequestMethodField, httpServerRequestMethodPOST)
	}
	return d.Set(httpServerRequestMethodField, httpServerRequestMethodGET)
}

func httpServerResponseRequestMethod(resp *tests.HttpServerTestResponse) (string, bool) {
	v := reflect.ValueOf(resp)
	if v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return "", false
		}
		v = v.Elem()
	}
	field := v.FieldByName("RequestMethod")
	if !field.IsValid() || field.Kind() != reflect.Pointer || field.IsNil() || field.Elem().Kind() != reflect.String {
		return "", false
	}

	return normalizeHTTPServerRequestMethod(field.Elem().String())
}

func normalizeHTTPServerRequestMethod(method string) (string, bool) {
	method = strings.ToLower(method)
	if method != httpServerRequestMethodGET && method != httpServerRequestMethodPOST {
		return "", false
	}
	return method, true
}
