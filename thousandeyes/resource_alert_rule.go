package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceAlertRule() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.AlertRule{}, schemas),
		Create: resourceAlertRuleCreate,
		Read:   resourceAlertRuleRead,
		Update: resourceAlertRuleUpdate,
		Delete: resourceAlertRuleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	resource.Schema["direction"] = schemas["direction-alert_rule"]
	return &resource
}

func resourceAlertRuleRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	remote, err := client.GetAlertRule(id)
	if err != nil {
		return err
	}
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceAlertRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.AlertRule{}).(*thousandeyes.AlertRule)
	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, Alert Rules actually require a few fields
	// to be retained on update.
	// Terraform schema validation should guarantee their existence.
	update.AlertType = thousandeyes.String(d.Get("alert_type").(string))
	update.Expression = thousandeyes.String(d.Get("expression").(string))

	// MinimumSources and MinimumSourcesPct are mutually exclusive
	minimumSources := d.Get("minimum_sources").(int)
	if minimumSources > 0 {
		update.MinimumSources = thousandeyes.Int(minimumSources)
	} else {
		update.MinimumSourcesPct = thousandeyes.Int(d.Get("minimum_sources_pct").(int))
	}

	update.RoundsViolatingRequired = thousandeyes.Int(d.Get("rounds_violating_required").(int))
	update.RoundsViolatingOutOf = thousandeyes.Int(d.Get("rounds_violating_out_of").(int))
	update.RuleName = thousandeyes.String(d.Get("rule_name").(string))

	_, err := client.UpdateAlertRule(id, *update)
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
