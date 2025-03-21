package thousandeyes

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/alerts"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/client"
)

func dataSourceThousandeyesAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAgentRead,

		Schema: map[string]*schema.Schema{
			"rule_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the alert rule.",
			},
			"rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique ID of the alert rule.",
			},
		},
		Description: "This data source allows you to configure alert rules. For more information, see [Creating and Editing Alert Rules](https://docs.thousandeyes.com/product-documentation/alerts/creating-and-editing-alert-rules).",
	}
}

func dataSourceThousandeyesAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.APIClient)
	api := (*alerts.AlertRulesAPIService)(&apiClient.Common)

	log.Printf("[INFO] Reading Thousandeyes alert rule")

	searchName := d.Get("rule_name").(string)

	req := api.GetAlertsRules()
	req = SetAidFromContext(apiClient.GetConfig().Context, req, req)

	resp, _, err := req.Execute()
	if err != nil {
		return err
	}

	var found *alerts.BaseRule

	for _, alertRule := range resp.GetAlertRules() {
		if alertRule.RuleName == searchName {
			found = &alertRule
			break
		}
	}

	if found == nil {
		return fmt.Errorf("unable to locate any alert rule with the name: %s", searchName)
	}
	log.Printf("[INFO] ## Found Alert Rule rule_id: %s - name: %s", *found.RuleId, found.RuleName)

	d.SetId(*found.RuleId)
	err = d.Set("rule_name", found.RuleName)
	if err != nil {
		return err
	}
	err = d.Set("rule_id", found.RuleId)
	if err != nil {
		return err
	}

	return nil
}
