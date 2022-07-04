package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

func resourceFTPServer() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(thousandeyes.FTPServer{}, schemas),
		Create: resourceFTPServerCreate,
		Read:   resourceFTPServerRead,
		Update: resourceFTPServerUpdate,
		Delete: resourceFTPServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
	resource.Schema["password"] = schemas["password-ftp"]
	resource.Schema["username"] = schemas["username-ftp"]
	return &resource
}

func resourceFTPServerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	remote, err := client.GetFTPServer(id)
	if err != nil {
		return err
	}
	err = ResourceRead(d, remote)
	if err != nil {
		return err
	}
	return nil
}

func resourceFTPServerUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	update := ResourceUpdate(d, &thousandeyes.FTPServer{}).(*thousandeyes.FTPServer)
	_, err := client.UpdateFTPServer(id, *update)
	if err != nil {
		return err
	}
	return resourceFTPServerRead(d, m)
}

func resourceFTPServerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteFTPServer(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceFTPServerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	local := buildFTPServerStruct(d)
	remote, err := client.CreateFTPServer(*local)
	if err != nil {
		return err
	}
	id := *remote.TestID
	d.SetId(strconv.FormatInt(id, 10))
	return resourceFTPServerRead(d, m)
}

func buildFTPServerStruct(d *schema.ResourceData) *thousandeyes.FTPServer {
	return ResourceBuildStruct(d, &thousandeyes.FTPServer{}).(*thousandeyes.FTPServer)
}
