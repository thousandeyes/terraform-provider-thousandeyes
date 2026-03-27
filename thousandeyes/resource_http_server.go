package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

const httpHeaderSourceModeField = "header_source_mode"
const httpHeaderSourceModeHeaders = "headers"
const httpHeaderSourceModeCustomHeaders = "custom_headers"

func resourceHTTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema:        ResourceSchemaBuild(tests.HttpServerTestRequest{}, schemas.CommonSchema, nil),
		Create:        resourceHTTPServerCreate,
		Read:          resourceHTTPServerRead,
		Update:        resourceHTTPServerUpdate,
		Delete:        resourceHTTPServerDelete,
		CustomizeDiff: normalizeHTTPServerHeadersDiff,
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
	return &resource
}

func resourceHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	log.Printf("[INFO] Reading Thousandeyes Resource %s", d.Id())

	api := (*tests.HTTPServerTestsAPIService)(&apiClient.Common)
	req := api.GetHttpServerTest(d.Id()).Expand(tests.AllowedExpandTestOptionsEnumValues)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil && IsNotFoundError(err) {
		log.Printf("[INFO] Resource was deleted - will recreate it")
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	if err := ResourceRead(context.Background(), d, resp); err != nil {
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

	req := api.UpdateHttpServerTest(d.Id()).HttpServerTestRequest(*update).Expand(tests.AllowedExpandTestOptionsEnumValues)
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

	req := api.CreateHttpServerTest().HttpServerTestRequest(*local).Expand(tests.AllowedExpandTestOptionsEnumValues)
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
