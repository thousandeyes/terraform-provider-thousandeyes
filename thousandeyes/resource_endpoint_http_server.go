package thousandeyes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointtests"
)

func resourceEndpointHTTPServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointHTTPServerCreate,
		Read:   resourceEndpointHTTPServerRead,
		Update: resourceEndpointHTTPServerUpdate,
		Delete: resourceEndpointHTTPServerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Description:   "Manages a scheduled Endpoint HTTP Server test. Agent selectors and other create-only settings require replacement when changed.",
		CustomizeDiff: resourceEndpointHTTPServerCustomizeDiff,
		Schema:        schemas.EndpointHTTPServerSchema(),
	}
}

func resourceEndpointHTTPServerCustomizeDiff(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	if err := validateEndpointTestSelector(ctx, d, meta); err != nil {
		return err
	}
	if d.NewValueKnown("ip_version") && d.NewValueKnown("protocol") {
		ipVersion := d.Get("ip_version").(string)
		if ipVersion != "" && ipVersion != "V4_ONLY" && d.Get("protocol").(string) != "icmp" {
			return fmt.Errorf("ip_version %q requires protocol = icmp; other protocols require V4_ONLY", ipVersion)
		}
	}
	// The API rejects URL changes on tests configured with HTTP authentication.
	if d.Id() != "" && d.HasChange("url") && d.Get("auth_type").(string) != "none" {
		oldURL, newURL := d.GetChange("url")
		if schemas.NormalizeEndpointHTTPServerURL(oldURL.(string)) != schemas.NormalizeEndpointHTTPServerURL(newURL.(string)) {
			return d.ForceNew("url")
		}
	}
	return nil
}

func buildEndpointHTTPServerCreate(d *schema.ResourceData) *endpointtests.EndpointHttpServerTestRequest {
	req := ResourceBuildStruct(d, &endpointtests.EndpointHttpServerTestRequest{})
	req.Url = schemas.NormalizeEndpointHTTPServerURL(req.Url)
	return req
}

func buildEndpointHTTPServerUpdate(d *schema.ResourceData) *endpointtests.EndpointHttpTestUpdate {
	req := ResourceBuildStruct(d, &endpointtests.EndpointHttpTestUpdate{})
	if req.Url != nil {
		url := schemas.NormalizeEndpointHTTPServerURL(*req.Url)
		req.Url = &url
	}
	return req
}

func endpointHTTPServerAPI(m interface{}) *endpointtests.HTTPServerEndpointScheduledTestsAPIService {
	return (*endpointtests.HTTPServerEndpointScheduledTestsAPIService)(&m.(*client.APIClient).Common)
}

func resourceEndpointHTTPServerCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	req := endpointHTTPServerAPI(m).CreateHttpServerEndpointScheduledTest().
		EndpointHttpServerTestRequest(*buildEndpointHTTPServerCreate(d))
	resp, _, err := SetAidFromContext(apiClient.GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	if resp.GetTestId() == "" {
		return fmt.Errorf("created endpoint HTTP server test returned no ID")
	}
	d.SetId(resp.GetTestId())
	// The create API always enables the test and does not accept is_enabled.
	// Disable it through the update API when the configuration requests it.
	if !d.Get("is_enabled").(bool) {
		return resourceEndpointHTTPServerUpdate(d, m)
	}
	return resourceEndpointHTTPServerRead(d, m)
}

func resourceEndpointHTTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).UpdateHttpServerEndpointScheduledTest(d.Id()).
		EndpointHttpTestUpdate(*buildEndpointHTTPServerUpdate(d))
	_, _, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil {
		return err
	}
	return resourceEndpointHTTPServerRead(d, m)
}

func resourceEndpointHTTPServerDelete(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).DeleteHttpServerEndpointScheduledTest(d.Id())
	httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		return err
	}
	d.SetId("")
	return nil
}

func resourceEndpointHTTPServerRead(d *schema.ResourceData, m interface{}) error {
	req := endpointHTTPServerAPI(m).GetHttpServerEndpointScheduledTest(d.Id())
	resp, httpResp, err := SetAidFromContext(m.(*client.APIClient).GetConfig().Context, req).Execute()
	if httpResp != nil && httpResp.StatusCode == 404 {
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}
	return flattenEndpointHTTPServer(d, resp)
}

func flattenEndpointHTTPServer(d *schema.ResourceData, resp *endpointtests.EndpointHttpServerTest) error {
	if err := ResourceRead(context.Background(), d, resp); err != nil {
		return err
	}
	// The create API calls this targetResponseTime, while reads return httpTargetTime.
	if resp.HttpTargetTime != nil {
		if err := d.Set("target_response_time", *resp.HttpTargetTime); err != nil {
			return err
		}
	}
	return flattenEndpointTestSelector(d, resp.AgentSelectorConfig)
}
