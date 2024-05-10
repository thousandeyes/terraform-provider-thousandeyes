package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceAlertRule() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.AlertRule{}, schemas, nil),
		Create: resourceAlertRuleCreate,
		Read:   resourceAlertRuleRead,
		Update: resourceAlertRuleUpdate,
		Delete: resourceAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create alert rules for ThousandEyes alerts. Alert rules define what alerts are sent, when, and to whom. For more information, see [Alert Rules](https://docs.thousandeyes.com/product-documentation/alerts#rule-configuration).",
	}
	resource.Schema["direction"] = schemas["direction-alert_rule"]
	return &resource
}

func resourceAlertRuleRead(d *schema.ResourceData, m interface{}) error {
	return GetResource(d, m, func(client *thousandeyes.Client, id int64) (interface{}, error) {
		var alertRule, err = client.GetAlertRule(id)
		if err != nil {
			return nil, err
		}
		alertRule.TestIds = testIds(*alertRule.Tests)
		alertRule.Tests = nil
		return alertRule, nil
	})
}

func resourceAlertRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Alert Rule %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, Alert Rules actually require the full list of
	// fields. Terraform schema validation should guarantee their existence.
	local := buildAlertRuleStruct(d)
	_, err := client.UpdateAlertRule(id, *local)
	if err != nil {
		return err
	}
	return resourceAlertRuleRead(d, m)
}

func resourceAlertRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteAlertRule(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceAlertRuleCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildAlertRuleStruct(d)
	remote, err := client.CreateAlertRule(*local)
	if err != nil {
		return err
	}
	id := remote.RuleID
	d.SetId(strconv.FormatInt(*id, 10))
	return resourceAlertRuleRead(d, m)
}

func buildAlertRuleStruct(d *schema.ResourceData) *thousandeyes.AlertRule {
	return ResourceBuildStruct(d, &thousandeyes.AlertRule{}).(*thousandeyes.AlertRule)
}

func testIds(tests []thousandeyes.GenericTest) *[]int64 {
	var testIds []int64
	for _, test := range tests {
		testIds = append(testIds, *test.TestID)
	}
	return &testIds
}
