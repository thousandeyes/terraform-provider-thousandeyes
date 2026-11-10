package thousandeyes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/endpointagents"
)

const (
	endpointAgentLicenseWirelessAdvantage  = "embed-wireless-advantage"
	endpointAgentLicenseWirelessEssentials = "embed-wireless-essentials"
)

func resourceEndpointAgent() *schema.Resource {
	return &schema.Resource{
		Create: resourceEndpointAgentCreate,
		Read:   resourceEndpointAgentRead,
		Update: resourceEndpointAgentUpdate,
		Delete: resourceEndpointAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the Endpoint Agent.",
			},
			"license_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateEndpointAgentLicenseType,
				Description:  "The license type assigned to the Endpoint Agent. The `embed-wireless-advantage` and `embed-wireless-essentials` license types can be read after import but cannot be assigned.",
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(endpointagents.STATUS_ENABLED),
					string(endpointagents.STATUS_DISABLED),
				}, false),
				Description: "The desired status of the Endpoint Agent: `enabled` or `disabled`.",
			},
		},
		Description: "Manages an existing ThousandEyes Endpoint Agent. This resource is import-only because the API does not support creating Endpoint Agents. Import an existing agent by ID before managing it.",
	}
}

func validateEndpointAgentLicenseType(value interface{}, key string) ([]string, []error) {
	license, ok := value.(string)
	if ok && (license == endpointAgentLicenseWirelessAdvantage || license == endpointAgentLicenseWirelessEssentials) {
		return nil, []error{fmt.Errorf("%s %q can be read from an imported Endpoint Agent but cannot be assigned", key, license)}
	}
	return validation.StringInSlice([]string{
		string(endpointagents.AGENTLICENSETYPE_ESSENTIALS),
		string(endpointagents.AGENTLICENSETYPE_ADVANTAGE),
		string(endpointagents.AGENTLICENSETYPE_EMBEDDED),
	}, false)(value, key)
}

func resourceEndpointAgentCreate(_ *schema.ResourceData, _ interface{}) error {
	return fmt.Errorf("endpoint agents cannot be created through Terraform; import an existing Endpoint Agent by ID")
}

func resourceEndpointAgentRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	log.Printf("[INFO] Reading ThousandEyes Endpoint Agent %s", d.Id())

	agent, response, err := getEndpointAgentResource(apiClient, d.Id())
	if err != nil {
		if response != nil && response.StatusCode == http.StatusNotFound {
			log.Printf("[INFO] Endpoint Agent %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}
	if agent == nil || agent.GetId() == "" {
		return fmt.Errorf("endpoint agent %s response did not include an ID", d.Id())
	}

	d.SetId(agent.GetId())
	if err := d.Set("name", agent.GetName()); err != nil {
		return err
	}
	if err := d.Set("license_type", string(agent.GetLicenseType())); err != nil {
		return err
	}
	if err := d.Set("status", string(agent.GetStatus())); err != nil {
		return err
	}
	return nil
}

func getEndpointAgentResource(apiClient *client.APIClient, id string) (*endpointagents.EndpointAgent, *http.Response, error) {
	api := (*endpointagents.EndpointAgentsAPIService)(&apiClient.Common)
	req := api.GetEndpointAgent(id)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	return req.Execute()
}

func resourceEndpointAgentUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointagents.EndpointAgentsAPIService)(&apiClient.Common)
	log.Printf("[INFO] Updating ThousandEyes Endpoint Agent %s", d.Id())

	if d.HasChange("license_type") {
		license := d.Get("license_type").(string)
		if license == endpointAgentLicenseWirelessAdvantage || license == endpointAgentLicenseWirelessEssentials {
			return fmt.Errorf("license_type %q cannot be assigned to an Endpoint Agent; the API only allows this license type to be read", license)
		}
	}

	update := endpointagents.EndpointAgentUpdate{}
	updateNeeded := false
	if d.HasChange("name") {
		update.SetName(d.Get("name").(string))
		updateNeeded = true
	}
	if d.HasChange("license_type") {
		update.SetLicenseType(endpointagents.AgentLicenseType(d.Get("license_type").(string)))
		updateNeeded = true
	}
	if updateNeeded {
		req := api.UpdateEndpointAgent(d.Id()).EndpointAgentUpdate(update)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)
		_, httpResp, err := req.Execute()
		if err != nil && (httpResp == nil || httpResp.StatusCode >= http.StatusMultipleChoices) {
			return err
		}
	}

	if d.HasChange("status") {
		var httpResp *http.Response
		var err error
		switch status := d.Get("status").(string); status {
		case string(endpointagents.STATUS_ENABLED):
			req := api.EnableEndpointAgent(d.Id())
			req = SetAidFromContext(apiClient.GetConfig().Context, req)
			_, httpResp, err = req.Execute()
		case string(endpointagents.STATUS_DISABLED):
			req := api.DisableEndpointAgent(d.Id())
			req = SetAidFromContext(apiClient.GetConfig().Context, req)
			_, httpResp, err = req.Execute()
		default:
			return fmt.Errorf("unsupported Endpoint Agent status %q; expected enabled or disabled", status)
		}
		if err != nil && (httpResp == nil || httpResp.StatusCode >= http.StatusMultipleChoices) {
			return err
		}
	}

	return resourceEndpointAgentRead(d, meta)
}

func resourceEndpointAgentDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*endpointagents.EndpointAgentsAPIService)(&apiClient.Common)
	log.Printf("[INFO] Deleting ThousandEyes Endpoint Agent %s", d.Id())

	req := api.DeleteEndpointAgent(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)
	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}
