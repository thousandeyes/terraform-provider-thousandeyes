package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/william20111/go-thousandeyes"
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
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetAlertRule(id)
	if err != nil {
		return err
	}
	ResourceRead(d, remote)
	return nil
}

func resourceAlertRuleUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.AlertRule{}).(*thousandeyes.AlertRule)
	_, err := client.UpdateAlertRule(id, *update)
	if err != nil {
		return err
	}
	return resourceAlertRuleRead(d, m)
}

func resourceAlertRuleDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
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
	d.SetId(strconv.Itoa(id))
	return resourceAlertRuleRead(d, m)
}

func buildAlertRuleStruct(d *schema.ResourceData) *thousandeyes.AlertRule {
	return ResourceBuildStruct(d, &thousandeyes.AlertRule{}).(*thousandeyes.AlertRule)
}
