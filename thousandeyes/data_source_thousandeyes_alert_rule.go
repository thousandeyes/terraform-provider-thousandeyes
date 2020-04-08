package thousandeyes

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
)

func dataSourceThousandeyesAlertRule() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceThousandeyesAlertRuleRead,

		Schema: map[string]*schema.Schema{
			"rule_id": {
				Type:          schema.TypeInt,
				Optional:      true,
				Description:   "Rule ID of alert rule",
				ConflictsWith: []string{"rule_name"},
			},
			"rule_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Name of the alert rule",
				ConflictsWith: []string{"rule_id"},
			},
			"expression": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "string expression of alert rule",
			},
			"direction": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "optional field with one of the following values: TO_TARGET, FROM_TARGET, BIDIRECTIONAL, for applicable alert types (eg. path trace, End-to-End (Agent) etc.)",
			},
			"notify_on_clear": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "1 to send notification when alert clears",
			},
			"default": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Alert rules allow up to 1 alert rule to be selected as a default for each type. By checking the default option, this alert rule will be automatically included on subsequently created tests that test a metric used in alerting here",
			},
			"alert_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of alert rule, as determined by metric selection",
			},
			"minimum_sources": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum number of agents or monitors that must meet the specified criteria in order to trigger the alert",
			},
			"minimum_sources_pct": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the minimum percentage of all assigned agents or monitors that must meet the specified criteria in order to trigger the alert",
			},
			"rounds_violating_out_of": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "applies to only v6 and higher, specifies the divisor (y value) for the “X of Y times” condition.",
			},
			"rounds_violating_required": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "applies to only v6 and higher, specifies the numerator (x value) for the X of Y times” condition",
			},
			"rounds_violating_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "allows user to enforce a condition on alerts that will required the same agent to meet the condition(s) multiple rounds in a row, rather than any agents meeting the condition(s) ",
			},
		},
	}
}

func dataSourceThousandeyesAlertRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*thousandeyes.Client)

	var found *thousandeyes.AlertRule

	searchName := d.Get("rule_name").(string)
	searchRuleID := d.Get("rule_id").(int)

	alertRules, err := client.GetAlertRules()
	if err != nil {
		return err
	}

	if searchName != "" {
		log.Printf("[INFO] ###### Reading Thousandeyes alert rules by name [%s]", searchName)

		for _, ar := range *alertRules {
			if ar.RuleName == searchName {
				found = &ar
				break
			}

		}
	} else if searchRuleID != 0 {
		for _, ar := range *alertRules {
			if ar.RuleID == searchRuleID {
				found = &ar
				break
			}

		}
	} else {
		return fmt.Errorf("must define name or rule id")
	}

	d.SetId(strconv.Itoa(found.RuleID))
	d.Set("rule_name", found.RuleName)
	d.Set("rule_id", found.RuleID)

	return nil
}
