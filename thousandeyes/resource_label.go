package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceGroupLabel() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.GroupLabel{}, schemas),
		Create: resourceGroupLabelCreate,
		Read:   resourceGroupLabelRead,
		Update: resourceGroupLabelUpdate,
		Delete: resourceGroupLabelDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create labels for ThousandEyes agents and test groups. Labels are used to quickly select an agent group or a test group, and apply them to a specific context (for example, applying a group of tests to an agent). For more information, see [Working with Labels](https://docs.thousandeyes.com/product-documentation/internet-and-wan-monitoring/tests/working-with-labels-for-agent-and-test-groups).",
	}
	resource.Schema["type"] = schemas["type-label"]
	resource.Schema["agents"] = schemas["agents-label"]
	resource.Schema["tests"].Elem = &schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.GenericTest{}, schemas),
	}
	return &resource
}

func resourceGroupLabelRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Label %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	remote, err := client.GetGroupLabel(id)
	if err != nil {
		return err
	}
	// In order to prevent schema conficts for test responses,  we retain
	// the stored state for tests attached to a group to just a test ID.
	testIDs := []thousandeyes.GenericTest{}
	if remote.Tests != nil {
		for _, v := range *remote.Tests {
			test := thousandeyes.GenericTest{TestID: v.TestID}
			testIDs = append(testIDs, test)
		}
	}
	remote.Tests = &testIDs
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceGroupLabelUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Label %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	update := ResourceUpdate(d, &thousandeyes.GroupLabel{}).(*thousandeyes.GroupLabel)
	// While most ThousandEyes updates only require updated fields and specifically
	// disallow some fields on update, Labels require the label name field to be
	// retained on update otherwise the call fails.
	update.Name = thousandeyes.String(d.Get("name").(string))
	_, err := client.UpdateGroupLabel(id, *update)
	if err != nil {
		return err
	}
	return resourceGroupLabelRead(d, m)
}

func resourceGroupLabelDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Label %s", d.Id())
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	if err := client.DeleteGroupLabel(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceGroupLabelCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Label %s", d.Id())
	local := buildGroupLabelStruct(d)
	remote, err := client.CreateGroupLabel(*local)
	if err != nil {
		return err
	}
	id := *remote.GroupID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceGroupLabelRead(d, m)
}

func buildGroupLabelStruct(d *schema.ResourceData) *thousandeyes.GroupLabel {
	return ResourceBuildStruct(d, &thousandeyes.GroupLabel{}).(*thousandeyes.GroupLabel)
}
