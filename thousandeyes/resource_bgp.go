package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceBGP() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the test",
			},
			"bgp_monitors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "array of BGP Monitor objects",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"monitor_id": {
							Type:        schema.TypeInt,
							Description: "monitor id",
							Optional:    true,
						},
					},
				},
			},
			"include_covered_prefixes": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to include queries for subprefixes detected under this prefix",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"prefix": {
				Type:        schema.TypeString,
				Description: "BGP network address prefix",
				Required:    true,
				// a.b.c.d is a network address, with the prefix length defined as e.
				// Prefixes can be any length from 8 to 24
				// Can only use private BGP monitors for a local prefix.
			},
			"use_public_bgp": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to automatically add all available Public BGP Monitors",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
		},
		Create: resourceBGPCreate,
		Read:   resourceBGPRead,
		Update: resourceBGPUpdate,
		Delete: resourceBGPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceBGPRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetBGP(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("bgp_monitors", test.BgpMonitors)
	d.Set("include_covered_prefixes", test.IncludeCoveredPrefixes)
	d.Set("prefix", test.Prefix)
	d.Set("use_public_bgp", test.UsePublicBGP)
	return nil
}

func resourceBGPUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.BGP
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("bgp_monitors") {
		update.BgpMonitors = expandBGPMonitors(d.Get("bgp_monitors").([]interface{}))
	}
	if d.HasChange("include_covered_prefixes") {
		update.IncludeCoveredPrefixes = d.Get("include_covered_prefixes").(int)
	}
	if d.HasChange("prefix") {
		update.Prefix = d.Get("prefix").(string)
	}
	if d.HasChange("use_public_bgp") {
		update.UsePublicBGP = d.Get("use_public_bgp").(int)
	}
	_, err := client.UpdateBGP(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceBGPRead(d, m)
}

func resourceBGPDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteBGP(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceBGPCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	bgpServer := buildBGPStruct(d)
	bgpTest, err := client.CreateBGP(*bgpServer)
	if err != nil {
		return err
	}
	testID := bgpTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceBGPRead(d, m)
}

func buildBGPStruct(d *schema.ResourceData) *thousandeyes.BGP {
	bgpServer := thousandeyes.BGP{
		TestName:               d.Get("name").(string),
		BgpMonitors:            expandBGPMonitors(d.Get("bgp_monitors").([]interface{})),
		IncludeCoveredPrefixes: d.Get("include_covered_prefixes").(int),
		Prefix:                 d.Get("prefix").(string),
		UsePublicBGP:           d.Get("use_public_bgp").(int),
	}

	return &bgpServer
}
