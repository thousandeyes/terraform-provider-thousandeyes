package thousandeyes

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/william20111/go-thousandeyes"
)

func resourceWebTransaction() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_type": {
				Type:         schema.TypeString,
				Description:  "auth type",
				Optional:     true,
				Default:      "NONE",
				ValidateFunc: validation.StringInSlice([]string{"NONE", "BASIC", "NTLM", "KERBEROS"}, false),
			},
			"bandwidth_measurements": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to measure bandwidth; defaults to 0. Only applies to Enterprise Agents assigned to the test, and requires that networkMeasurements is set.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"content_regex": {
				Type: schema.TypeString,
				Description: "regular Expressions	Verify content using a regular expression. This field does not require escaping",
				Optional: true,
				Default:  "NONE",
			},
			"desired_status_code": {
				Type: schema.TypeString,
				Description: "A valid HTTP response code	Set to the value you’re interested in retrieving",
				Optional: true,
			},
			"http_target_time": {
				Type:         schema.TypeInt,
				Description:  "target time for HTTP server completion; specified in milliseconds",
				Optional:     true,
				Default:      1000,
				ValidateFunc: validation.IntBetween(100, 5000),
			},
			"http_time_limit": {
				Type:        schema.TypeInt,
				Description: "target time for HTTP server limits; specified in seconds",
				Default:     5,
				Optional:    true,
			},
			"include_headers": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to capture response headers for objects loaded by the test.Default is 1.",
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"mtu_measurements": {
				Type:         schema.TypeInt,
				Description:  "set to 1 to measure MTU sizes on network from agents to the target.",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"network_measurements": {
				Type:         schema.TypeInt,
				Description:  "choose 1 to enable network measurements, 0 to disable; defaults to 1.",
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(0, 1),
			},
			"num_path_traces": {
				Type:         schema.TypeInt,
				Description:  "number of path traces. default 3.",
				Default:      3,
				Optional:     true,
				ValidateFunc: validation.IntBetween(3, 10),
			},
			"password": {
				Type:        schema.TypeString,
				Description: "password to be used for Basic authentication",
				Optional:    true,
			},
			"probe_mode": {
				Type:        schema.TypeString,
				Description: "probe mode used by End-to-end Network Test; only valid if protocol is set to TCP; defaults to AUTO",
				Optional:    true,
				Default:     "AUTO",
			},
			"protocol": {
				Type:         schema.TypeString,
				Description:  "protocol used by dependent Network tests (End-to-end, Path Trace, PMTUD); defaults to TCP",
				Optional:     true,
				Default:      "TCP",
				ValidateFunc: validation.StringInSlice([]string{"TCP", "ICMP"}, false),
			},
			"ssl_version_id": {
				Type:         schema.TypeInt,
				Description:  "0 for auto, 3 for SSLv3, 4 for TLS v1.0, 5 for TLS v1.1, 6 for TLS v1.2",
				Optional:     true,
				Default:      0,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"sub_interval": {
				Type:         schema.TypeInt,
				Description:  "subinterval for round-robin testing (in seconds), must be less than or equal to interval",
				Optional:     true,
				ValidateFunc: validation.IntInSlice([]int{60, 120, 300, 600, 900, 1200, 1800, 3600}),
			},
			"target_time": {
				Type:         schema.TypeInt,
				Description:  "target time for completion, defaults to 50% of time limit; specified in seconds",
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"time_limit": {
				Type:         schema.TypeInt,
				Description:  "time limit for transaction; defaults to 30s",
				Optional:     true,
				Default:      30,
				ValidateFunc: validation.IntBetween(1, 60),
			},
			"use_ntlm": {
				Type:        schema.TypeInt,
				Description: "choose 0 to use Basic Authentication, or omit field.Requires username/password to be set",
				Optional:    true,
			},
			"user_agent": {
				Type:        schema.TypeString,
				Description: "user-agent string to be provided during the test",
				Optional:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "username to be used for Basic authentication",
				Optional:    true,
			},
			"verify_certificate": {
				Type:        schema.TypeInt,
				Description: "set to 0 to ignore certificate errors (defaults to 1)",
				Optional:    true,
				Default:     1,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the test",
			},
			"interval": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "interval to run http server test on",
			},
			"transaction_script": {
				Type:        schema.TypeString,
				Description: "selenium transaction script",
				Required:    true,
			},
			"http_version": {
				Type:         schema.TypeInt,
				Description:  "2 for default (prefer HTTP/2), 1 for HTTP/1.1 only",
				Default:      2,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1, 2),
			},
			"url": {
				Type:        schema.TypeString,
				Description: "target for the test",
				Required:    true,
			},
			"agents": {
				Type:        schema.TypeList,
				Description: "agents to use ",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": {
							Type:        schema.TypeInt,
							Description: "agent id",
							Optional:    true,
						},
					},
				},
			},
		},
		Create: resourceWebTransactionCreate,
		Read:   resourceWebTransactionRead,
		Update: resourceWebTransactionUpdate,
		Delete: resourceWebTransactionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceWebTransactionRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Reading Thousandeyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	test, err := client.GetWebTransaction(id)
	if err != nil {
		return err
	}

	d.Set("name", test.TestName)
	d.Set("auth_type", test.AuthType)
	d.Set("transaction_script", test.TransactionScript)
	d.Set("interval", test.Interval)
	d.Set("http_version", test.HTTPVersion)
	d.Set("url", test.URL)
	d.Set("agents", test.Agents)
	d.Set("bandwidth_measurements", test.BandwidthMeasurements)
	d.Set("content_regex", test.ContentRegex)
	d.Set("desired_status_code", test.DesiredStatusCode)
	d.Set("http_target_time", test.HTTPTargetTime)
	d.Set("http_time_limit", test.HTTPTimeLimit)
	d.Set("include_headers", test.IncludeHeaders)
	d.Set("mtu_measurements", test.MtuMeasurements)
	d.Set("network_measurements", test.NetworkMeasurements)
	d.Set("num_path_traces", test.NumPathTraces)
	d.Set("password", test.Password)
	d.Set("probe_mode", test.ProbeMode)
	d.Set("protocol", test.Protocol)
	d.Set("ssl_version_id", test.SslVersionID)
	d.Set("sub_interval", test.Subinterval)
	d.Set("target_time", test.TargetTime)
	d.Set("time_limit", test.TimeLimit)
	d.Set("use_ntlm", test.UseNtlm)
	d.Set("user_agent", test.UserAgent)
	d.Set("username", test.Username)
	d.Set("verify_certificate", test.VerifyCertificate)
	return nil
}

func resourceWebTransactionUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Updating ThousandEyes Test %s", d.Id())
	d.Partial(true)
	id, _ := strconv.Atoi(d.Id())
	var update thousandeyes.WebTransaction
	if d.HasChange("agents") {
		update.Agents = expandAgents(d.Get("agents").([]interface{}))
	}
	if d.HasChange("auth_type") {
		update.AuthType = d.Get("auth_type").(string)
	}
	if d.HasChange("name") {
		update.TestName = d.Get("name").(string)
	}
	if d.HasChange("transaction_script") {
		update.TransactionScript = d.Get("transaction_script").(string)
	}
	if d.HasChange("interval") {
		update.Interval = d.Get("interval").(int)
	}
	if d.HasChange("http_version") {
		update.HTTPVersion = d.Get("http_version").(int)
	}
	if d.HasChange("url") {
		update.URL = d.Get("url").(string)
	}
	if d.HasChange("bandwidth_measurements") {
		update.BandwidthMeasurements = d.Get("bandwidth_measurements").(int)
	}
	if d.HasChange("content_regex") {
		update.ContentRegex = d.Get("content_regex").(string)
	}
	if d.HasChange("desired_status_code") {
		update.DesiredStatusCode = d.Get("desired_status_code").(string)
	}
	if d.HasChange("http_target_time") {
		update.HTTPTargetTime = d.Get("http_target_time").(int)
	}
	if d.HasChange("http_time_limit") {
		update.HTTPTimeLimit = d.Get("http_time_limit").(int)
	}
	if d.HasChange("include_headers") {
		update.IncludeHeaders = d.Get("include_headers").(int)
	}
	if d.HasChange("mtu_measurements") {
		update.MtuMeasurements = d.Get("mtu_measurements").(int)
	}
	if d.HasChange("network_measurements") {
		update.NetworkMeasurements = d.Get("network_measurements").(int)
	}
	if d.HasChange("num_path_traces") {
		update.NumPathTraces = d.Get("num_path_traces").(int)
	}
	if d.HasChange("password") {
		update.Password = d.Get("password").(string)
	}
	if d.HasChange("probe_mode") {
		update.ProbeMode = d.Get("probe_mode").(string)
	}
	if d.HasChange("protocol") {
		update.Protocol = d.Get("protocol").(string)
	}
	if d.HasChange("ssl_version_id") {
		update.SslVersionID = d.Get("ssl_version_id").(int)
	}
	if d.HasChange("sub_interval") {
		update.Subinterval = d.Get("sub_interval").(int)
	}
	if d.HasChange("target_time") {
		update.TargetTime = d.Get("target_time").(int)
	}
	if d.HasChange("time_limit") {
		update.TimeLimit = d.Get("time_limit").(int)
	}
	if d.HasChange("use_ntlm") {
		update.UseNtlm = d.Get("use_ntlm").(int)
	}
	if d.HasChange("user_agent") {
		update.UserAgent = d.Get("user_agent").(string)
	}
	if d.HasChange("username") {
		update.Username = d.Get("username").(string)
	}
	if d.HasChange("verify_certificate") {
		update.VerifyCertificate = d.Get("verify_certificate").(int)
	}
	_, err := client.UpdateWebTransaction(id, update)
	if err != nil {
		return err
	}
	d.Partial(false)
	return resourceWebTransactionRead(d, m)
}

func resourceWebTransactionDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)

	log.Printf("[INFO] Deleting ThousandEyes Test %s", d.Id())
	id, _ := strconv.Atoi(d.Id())
	if err := client.DeleteWebTransaction(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceWebTransactionCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	log.Printf("[INFO] Creating ThousandEyes Test %s", d.Id())
	webTrans := buildWebTransactionStruct(d)
	httpTest, err := client.CreateWebTransaction(*webTrans)
	if err != nil {
		return err
	}
	testID := httpTest.TestID
	d.SetId(strconv.Itoa(testID))
	return resourceWebTransactionRead(d, m)
}

func buildWebTransactionStruct(d *schema.ResourceData) *thousandeyes.WebTransaction {
	transaction := thousandeyes.WebTransaction{
		TestName:            d.Get("name").(string),
		TransactionScript:   d.Get("transaction_script").(string),
		HTTPVersion:         d.Get("http_version").(int),
		URL:                 d.Get("url").(string),
		Interval:            d.Get("interval").(int),
		HTTPTargetTime:      d.Get("http_target_time").(int),
		MtuMeasurements:     d.Get("mtu_measurements").(int),
		NetworkMeasurements: d.Get("network_measurements").(int),
		Agents:              expandAgents(d.Get("agents").([]interface{})),
	}
	if attr, ok := d.GetOk("http_version"); ok {
		transaction.HTTPVersion = attr.(int)
	}
	if attr, ok := d.GetOk("bandwidth_measurements"); ok {
		transaction.BandwidthMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("content_regex"); ok {
		transaction.ContentRegex = attr.(string)
	}
	if attr, ok := d.GetOk("desired_status_code"); ok {
		transaction.DesiredStatusCode = attr.(string)
	}
	if attr, ok := d.GetOk("http_target_time"); ok {
		transaction.HTTPTargetTime = attr.(int)
	}
	if attr, ok := d.GetOk("http_time_limit"); ok {
		transaction.HTTPTimeLimit = attr.(int)
	}
	if attr, ok := d.GetOk("include_headers"); ok {
		transaction.IncludeHeaders = attr.(int)
	}
	if attr, ok := d.GetOk("mtu_measurements"); ok {
		transaction.MtuMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("network_measurements"); ok {
		transaction.NetworkMeasurements = attr.(int)
	}
	if attr, ok := d.GetOk("num_path_traces"); ok {
		transaction.NumPathTraces = attr.(int)
	}
	if attr, ok := d.GetOk("password"); ok {
		transaction.Password = attr.(string)
	}
	if attr, ok := d.GetOk("probe_mode"); ok {
		transaction.ProbeMode = attr.(string)
	}
	if attr, ok := d.GetOk("protocol"); ok {
		transaction.Protocol = attr.(string)
	}
	if attr, ok := d.GetOk("ssl_version_id"); ok {
		transaction.SslVersionID = attr.(int)
	}
	if attr, ok := d.GetOk("sub_interval"); ok {
		transaction.Subinterval = attr.(int)
	}
	if attr, ok := d.GetOk("target_time"); ok {
		transaction.TargetTime = attr.(int)
	}
	if attr, ok := d.GetOk("time_limit"); ok {
		transaction.TimeLimit = attr.(int)
	}
	if attr, ok := d.GetOk("use_ntlm"); ok {
		transaction.UseNtlm = attr.(int)
	}
	if attr, ok := d.GetOk("user_agent"); ok {
		transaction.UserAgent = attr.(string)
	}
	if attr, ok := d.GetOk("username"); ok {
		transaction.Username = attr.(string)
	}
	if attr, ok := d.GetOk("verify_certificate"); ok {
		transaction.VerifyCertificate = attr.(int)
	}
	return &transaction
}
