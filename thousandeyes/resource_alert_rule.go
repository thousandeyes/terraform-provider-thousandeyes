package thousandeyes

import (
	"context"
	"log"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func resourceAlertRule() *schema.Resource {
	resource := schema.Resource{
		Schema: schemas.AlertRuleSchema,
		Create: resourceAlertRuleCreate,
		Read:   resourceAlertRuleRead,
		Update: resourceAlertRuleUpdate,
		Delete: resourceAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create alert rules for ThousandEyes alerts. Alert rules define what alerts are sent, when, and to whom. For more information, see [Alert Rules](https://docs.thousandeyes.com/product-documentation/alerts#rule-configuration).",
	}
	resource.Schema["test_ids"] = schemas.AlertRuleSchema["test_ids"]
	return &resource
}

func resourceAlertRuleRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(context.Background(), d, m, func(apiClient *client.APIClient, id string) (interface{}, error) {
		api := (*alerts.AlertRulesAPIService)(&apiClient.Common)

		req := api.GetAlertRule(id)
		req = SetAidFromContext(apiClient.GetConfig().Context, req)

		alertRule, _, err := req.Execute()
		if err != nil {
			return nil, err
		}

		return alertRule, nil
	})
}

func resourceAlertRuleUpdate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*alerts.AlertRulesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Updating ThousandEyes Alert Rule %s", d.Id())

	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, Alert Rules actually require the full list of
	// fields. Terraform schema validation should guarantee their existence.
	local := buildAlertRuleStruct(d)
	req := api.UpdateAlertRule(d.Id()).RuleDetailUpdate(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	_, _, err := req.Execute()
	if err != nil {
		return err
	}

	return resourceAlertRuleRead(d, m)
}

func resourceAlertRuleDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*alerts.AlertRulesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())

	req := api.DeleteAlertRule(d.Id())
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	if _, err := req.Execute(); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAlertRuleCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.APIClient)
	api := (*alerts.AlertRulesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())

	local := buildAlertRuleStruct(d)
	req := api.CreateAlertRule().RuleDetailUpdate(*local)
	req = SetAidFromContext(apiClient.GetConfig().Context, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	id := *resp.RuleId
	d.SetId(id)
	return resourceAlertRuleRead(d, m)
}

func buildAlertRuleStruct(d *schema.ResourceData) *alerts.RuleDetailUpdate {
	return ResourceBuildStruct(d, &alerts.RuleDetailUpdate{})
}
