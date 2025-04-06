package thousandeyes

import (
	"context"
	"reflect"
	"testing"

	"github.com/thousandeyes/terraform-provider-thousandeyes/thousandeyes/schemas"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/administrative"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func getReferenceData(schemaData map[string]*schema.Schema, attrs map[string]string) *schema.ResourceData {
	referenceResource := schema.Resource{
		Schema: schemaData,
	}
	referenceState := terraform.InstanceState{
		ID:         "testState",
		Attributes: attrs,
	}

	return referenceResource.Data(&referenceState)
}

func TestResourceBuildStruct(t *testing.T) {
	prefix := "8.19.2.2/19"
	cmpStruct := tests.BgpTestRequest{
		Prefix: prefix,
	}
	newStruct := tests.BgpTestRequest{}
	attrs := map[string]string{
		"prefix": "8.19.2.2/19",
	}
	d := getReferenceData(schemas.CommonSchema, attrs)
	ResourceBuildStruct(d, &newStruct)

	if newStruct.Prefix != cmpStruct.Prefix {
		t.Error("Building resource did not assign struct field correctly.")
	}

}

func TestResourceRead(t *testing.T) {
	prefix := "8.19.2.2/19"
	attrs := map[string]string{}
	d := getReferenceData(schemas.CommonSchema, attrs)
	remoteResource := tests.BgpTestResponse{
		Prefix: prefix,
	}
	err := ResourceRead(context.TODO(), d, &remoteResource)
	if err != nil {
		t.Errorf("Setting resource data returned error: %+v", err.Error())
	}
	if d.Get("prefix") != remoteResource.Prefix {
		t.Errorf("Reading resource did not assign resource data correctly.\nStruct is %+v\nResource is %+v", remoteResource, d.State().Attributes)
	}
}

func TestResourceReadValue(t *testing.T) {
	testStruct := administrative.RoleDetail{
		Name:      getPointer("TestRole"),
		RoleId:    getPointer("1"),
		IsBuiltin: getPointer(false),
		Permissions: []administrative.Permission{
			{
				IsManagementPermission: getPointer(false),
				Label:                  getPointer("foo"),
				PermissionId:           getPointer("27"),
			},
			{
				IsManagementPermission: getPointer(true),
				Label:                  getPointer("bar"),
				PermissionId:           getPointer("28"),
			},
		},
	}
	resultMap := map[string]interface{}{
		"name":       getPointer("TestRole"),
		"role_id":    getPointer("1"),
		"is_builtin": getPointer(false),
		"permissions": []map[string]interface{}{
			{
				"is_management_permission": getPointer(false),
				"label":                    getPointer("foo"),
				"permission_id":            getPointer("27"),
			},
			{
				"is_management_permission": getPointer(true),
				"label":                    getPointer("bar"),
				"permission_id":            getPointer("28"),
			},
		},
	}

	result, err := ReadValue(&testStruct)
	if err != nil {
		t.Errorf("Error running ReadValue: %+v", err)
	}

	// We need to inform Go of the type of the "permissions" field in the result,
	// in order for reflect.DeepEqual to identify them as equivalent.
	// Since it's a list of maps we can't do that in one step.
	typedResult := result.(map[string]interface{})
	permissions := typedResult["permissions"].([]interface{})
	var typedPermissions []map[string]interface{}
	for _, v := range permissions {
		typedPermissions = append(typedPermissions, v.(map[string]interface{}))
	}
	typedResult["permissions"] = typedPermissions

	if reflect.DeepEqual(typedResult, resultMap) != true {
		t.Errorf("Struct processing and test map do not match: \n\n%+v\n\n%+v\n", result, resultMap)
	}

}

func TestFixReadValues(t *testing.T) {
	var output interface{}
	var err error

	// non-map, non-list
	normalInput := 4
	normalTarget := 4
	output, err = FixReadValues(context.TODO(), nil, normalInput, getPointer("normal"))
	if err != nil {
		t.Errorf("normal input returned error: %s", err.Error())
	}
	if output.(int) != 4 {
		t.Errorf("Returned wrong value for int input. Received  %#v, expected %#v", output, normalTarget)
	}

	// agents
	agentsInput := []interface{}{
		map[string]interface{}{
			"agent_name": "foo",
			"agent_id":   "1",
		},
		map[string]interface{}{
			"agent_name": "bar",
			"agent_id":   "2",
		},
	}
	agentsTarget := []interface{}{
		"1", "2",
	}
	output, err = FixReadValues(context.TODO(), nil, agentsInput, getPointer("agents"))
	if err != nil {
		t.Errorf("agents input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, agentsTarget) != true {
		t.Errorf("Values not stripped correctly from agents input: Received %#v Expected %#v", output, agentsTarget)
	}

	// dns_servers
	dnsServersInput := []interface{}{
		map[string]interface{}{
			"server_name": "foo.com",
			"server_id":   "1",
		},
		map[string]interface{}{
			"server_name": "bar.com",
			"server_id":   "2",
		},
	}
	dnsServersTarget := []interface{}{
		"foo.com", "bar.com",
	}
	output, err = FixReadValues(context.TODO(), nil, dnsServersInput, getPointer("dns_servers"))
	if err != nil {
		t.Errorf("dns_servers input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, dnsServersTarget) != true {
		t.Errorf("Values not stripped correctly from dns_servers input: Received %#v Expected %#v", output, agentsTarget)
	}

	// _links
	fieldName := getPointer("_links")
	linksInput := map[string]interface{}{
		"self": map[string]interface{}{
			"href": "foo.com",
		},
	}
	linksTarget := "foo.com"
	output, err = FixReadValues(context.TODO(), nil, linksInput, fieldName)
	if err != nil {
		t.Errorf("links input returned error: %s", err.Error())
	}
	if *fieldName != "link" {
		t.Errorf("link name returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, linksTarget) != true {
		t.Errorf("Values not stripped correctly from links input: Received %#v Expected %#v", output, agentsTarget)
	}

	// alert_rules
	alertRulesInput := []interface{}{
		map[string]interface{}{
			"rule_name": getPointer("foo"),
			"rule_id":   getPointer("1"),
			"default":   getPointer(false),
		},
		map[string]interface{}{
			"rule_name": getPointer("bar"),
			"rule_id":   getPointer("2"),
			"default":   getPointer(false),
		},
		map[string]interface{}{
			"rule_name": getPointer("bar"),
			"rule_id":   getPointer("3"),
			"default":   getPointer(true),
		},
	}
	alertRulesTarget := []interface{}{
		getPointer("1"), getPointer("2"), getPointer("3"),
	}
	output, err = FixReadValues(context.TODO(), nil, alertRulesInput, getPointer("alert_rules"))
	if err != nil {
		t.Errorf("alert_rules input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, alertRulesTarget) != true {
		t.Errorf("Values not stripped correctly from alert_rules input: Received %#v Expected %#v", output, alertRulesTarget)
	}

	// bgp_monitors
	monitorsInput := []interface{}{
		map[string]interface{}{
			"monitor_name": getPointer("foo"),
			"monitor_id":   getPointer("1"),
			"monitor_type": getPointer(tests.MonitorType("public")),
		},
		map[string]interface{}{
			"monitor_name": getPointer("bar"),
			"monitor_id":   getPointer("2"),
			"monitor_type": getPointer(tests.MonitorType("private")),
		},
	}
	monitorsTarget := []interface{}{
		getPointer("2"),
	}
	output, err = FixReadValues(context.TODO(), nil, monitorsInput, getPointer("monitors"))
	if err != nil {
		t.Errorf("bgp_monitors input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, monitorsTarget) != true {
		t.Errorf("Values not stripped correctly from bgp_monitors input: Received %#v Expected %#v", output, monitorsTarget)
	}

	//	shared_with_accounts
	accountsInput := []interface{}{
		map[string]interface{}{
			"name": getPointer("foo"),
			"aid":  getPointer("1"),
		},
		map[string]interface{}{
			"name": getPointer("bar"),
			"aid":  getPointer("2"),
		},
	}
	accountsTarget := []interface{}{
		map[string]interface{}{
			"aid": getPointer("1"),
		},
	}

	output, err = FixReadValues(context.WithValue(context.TODO(), accountGroupIdKey, "2"), nil, accountsInput, getPointer("shared_with_accounts"))
	if err != nil {
		t.Errorf("shared_with_accounts input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, accountsTarget) != true {
		t.Errorf("Values not stripped correctly from shared_with_accounts input: Received %#v Expected %#v", output, accountsTarget)
	}
	//  We should fail if account_group_id isn't set and the list of account groups is > 1
	output, err = FixReadValues(context.TODO(), nil, accountsInput, getPointer("shared_with_accounts"))
	if err == nil {
		t.Errorf("Error was not returned when shared_with_accounts length was > 1 and account_group_id  was not set")
	}
	// We should not fail if account_group_id isn't set and the list of account groups is < 2
	accountsInput = []interface{}{
		map[string]interface{}{
			"name": "bar",
			"aid":  "2",
		},
	}
	output, err = FixReadValues(context.TODO(), nil, accountsInput, getPointer("shared_with_accounts"))
	if err != nil {
		t.Errorf("shared_with_accounts input returned error when shared_with_accounts wasn't set despite list of account groups being < 2: %s", err.Error())
	}
	if output != nil {
		t.Errorf("Values not stripped correctly from shared_with_accounts input: Received %#v Expected %#v", output, accountsTarget)
	}

	// use target map
	targetMapInput := map[string]map[string]interface{}{
		"target": {
			"source_1": nil,
			"source_2": nil,
			"source_3": nil,
		},
	}
	targetMapOutput := map[string]map[string]interface{}{
		"target": {
			"source_1": 1,
			"source_2": 2,
			"source_3": 3,
		},
	}
	nameSource1 := getPointer("source_1")
	FixReadValues(context.TODO(), targetMapInput, 1, nameSource1)
	FixReadValues(context.TODO(), targetMapInput, 2, getPointer("source_2"))
	FixReadValues(context.TODO(), targetMapInput, 3, getPointer("source_3"))
	if len(*nameSource1) != 0 {
		t.Errorf("Name wasn't cleared when target map was set")
	}
	if reflect.DeepEqual(targetMapOutput, targetMapInput) != true {
		t.Errorf("Target Map didn't set correctly: Received %#v Expected %#v", targetMapInput, targetMapOutput)
	}

	//	tests
	testsInput := []interface{}{
		map[string]interface{}{
			"test_name": "foo",
			"test_id":   "1",
		},
		map[string]interface{}{
			"test_name": "bar",
			"test_id":   "2",
		},
	}
	testsTarget := []interface{}{
		"1", "2",
	}
	fieldName = getPointer("tests")
	output, err = FixReadValues(context.TODO(), nil, testsInput, fieldName)
	if err != nil {
		t.Errorf("tests input returned error: %s", err.Error())
	}
	if *fieldName != "test_ids" {
		t.Errorf("tests name returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, testsTarget) != true {
		t.Errorf("Values not stripped correctly from tests input: Received %#v Expected %#v", output, nil)
	}

	// thirdParty notifications
	thirdPartyNotificationsInput := []interface{}{
		map[string]interface{}{
			"integration_id":   "sl-0000",
			"integration_type": "SLACK",
			"integration_name": "bitconnect",
			"target":           "https://slack.com/waso",
			"channel":          "#terraform",
		},
		map[string]interface{}{
			"integration_id":   "pgd-0000",
			"integration_type": "PAGER_DUTY",
			"integration_name": "PagerDuty notification",
			"auth_method":      "Auth Token",
		},
	}
	thirdPartyNotificationsTarget := []interface{}{
		map[string]interface{}{
			"integration_id":   "sl-0000",
			"integration_type": "SLACK",
		},
		map[string]interface{}{
			"integration_id":   "pgd-0000",
			"integration_type": "PAGER_DUTY",
		},
	}

	output, err = FixReadValues(context.TODO(), nil, thirdPartyNotificationsInput, getPointer("third_party"))
	if err != nil {
		t.Errorf("third party notifications input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, thirdPartyNotificationsTarget) != true {
		t.Errorf("Values not stripped correctly from third party notifications input: Received %#v Expected %#v", output, thirdPartyNotificationsTarget)
	}

	// webhook notifications
	webhookNotificationsInput := []interface{}{
		map[string]interface{}{
			"integration_id":   "wb-0000",
			"integration_type": "WEBHOOK",
			"integration_name": "TEAMS CHANNEL",
			"target":           "https://webhook.office.com",
		},
	}
	webhookNotificationsTarget := []interface{}{
		map[string]interface{}{
			"integration_id":   "wb-0000",
			"integration_type": "WEBHOOK",
		},
	}

	output, err = FixReadValues(context.TODO(), nil, webhookNotificationsInput, getPointer("webhook"))
	if err != nil {
		t.Errorf("webhook notifications input returned error: %s", err.Error())
	}
	if reflect.DeepEqual(output, webhookNotificationsTarget) != true {
		t.Errorf("Values not stripped correctly from webhook notifications input: Received %#v Expected %#v", output, webhookNotificationsTarget)
	}

	// emulated device id
	ctx := context.WithValue(context.Background(), emulationDeviceIdKey, struct{}{})
	targetEdId := "3000"
	output, err = FixReadValues(ctx, nil, getPointer("3000"), getPointer("emulated_device_id"))
	if err != nil {
		t.Errorf("emulated device id input returned error: %s", err.Error())
	}
	if output == nil {
		t.Errorf("emulated device id was set incorrectly: Received nil Expected %s", targetEdId)
	}
	if str := output.(*string); *str != targetEdId {
		t.Errorf("emulated device id was set incorrectly: Received %s Expected %s", *str, targetEdId)
	}
	output, err = FixReadValues(context.Background(), nil, getPointer("3000"), getPointer("emulated_device_id"))
	if err != nil {
		t.Errorf("emulated device id input returned error: %s", err.Error())
	}
	if output != nil {
		t.Errorf("emulated device id was set incorrectly: Received %s Expected nil", *output.(*string))
	}
}

func TestResourceUpdate(t *testing.T) {

}

func TestResourceSchemaBuild(t *testing.T) {
	type refStruct struct {
		FieldName string `json:"fieldName"`
		Port      int    `json:"port"`
	}

	refSchema := map[string]*schema.Schema{
		"field_name": {
			Type: schema.TypeString,
		},
		"port": {
			Type:     schema.TypeInt,
			Default:  41953,
			Optional: true,
		},
	}

	// test with no schema override
	schm := ResourceSchemaBuild(refStruct{}, refSchema, nil)

	for k, v := range refSchema {
		if _, ok := schm[k]; !ok {
			t.Errorf("Key %s missing from generated schema", k)
		}
		if reflect.DeepEqual(*v, *schm[k]) != true {
			t.Errorf("Schemas not equal: Reference schema is %+v\nNew Schema is %+v", refSchema, schm)
		}
	}

	// test with a schema override
	schemaOverride := map[string]*schema.Schema{
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	expectedSchemaWithOverride := map[string]*schema.Schema{
		"field_name": {
			Type: schema.TypeString,
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
		},
	}

	ovrSchm := ResourceSchemaBuild(refStruct{}, refSchema, schemaOverride)

	for k, v := range expectedSchemaWithOverride {
		if _, ok := ovrSchm[k]; !ok {
			t.Errorf("Key %s missing from generated schema", k)
		}
		if reflect.DeepEqual(*v, *ovrSchm[k]) != true {
			t.Errorf("Schemas not equal: Expected schema is %+v\nNew Schema is %+v", expectedSchemaWithOverride, ovrSchm)
		}
	}
}

func TestFillValue(t *testing.T) {
	type refStruct struct {
		FieldName string `json:"fieldName"`
	}
	attrs := map[string]string{}
	testSchemas := map[string]*schema.Schema{
		"testInt": {
			Type: schema.TypeInt,
		},
		"testIntAsString": {
			Type: schema.TypeString,
		},
		"testString": {
			Type: schema.TypeString,
		},
		"testSlice": {
			Type: schema.TypeList,
			Elem: schema.TypeInt,
		},
		"testStruct": {
			Type: schema.TypeMap,
		},
		"testStructSlice": {
			Type: schema.TypeList,
			Elem: schema.TypeMap,
		},
	}
	d := getReferenceData(testSchemas, attrs)

	var err error

	// Integer test
	err = d.Set("testInt", 42)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testInt': %+v", err.Error())
	}
	testInt := FillValue(d.Get("testInt"), 1).(int)
	if testInt != 42 {
		t.Errorf("Expected int '42' for testInt, but received '%+v' of type %+v", reflect.ValueOf(testInt), reflect.TypeOf(testInt))
	}
	err = d.Set("testIntAsString", "27")
	if err != nil {
		t.Errorf("Error setting resourceData for 'testInt': %+v", err.Error())
	}
	testIntAsString := FillValue(d.Get("testIntAsString"), 1).(int)
	if testIntAsString != 27 {
		t.Errorf("Expected int '27' for testIntAsString, but received '%+v' of type %+v", reflect.ValueOf(testIntAsString), reflect.TypeOf(testIntAsString))
	}

	// String test - should test default return path
	err = d.Set("testString", "foo")
	if err != nil {
		t.Errorf("Error setting resourceData for 'testString': %+v", err.Error())
	}
	testString := FillValue(d.Get("testString"), "bar").(string)
	if testString != "foo" {
		t.Errorf("Expected string 'foo', but received %+v of type %+v", reflect.ValueOf(testString), reflect.TypeOf(testString))
	}

	// Slice test
	err = d.Set("testSlice", []int{0, 1})
	if err != nil {
		t.Errorf("Error setting resourceData for 'testSlice': %+v", err.Error())
	}
	testSlice := FillValue(d.Get("testSlice"), []int{}).([]int)
	if reflect.DeepEqual(testSlice, []int{0, 1}) != true {
		t.Errorf("Expected []int{0,1} for testSlice, but received '%+v' of type %+v", reflect.ValueOf(testSlice), reflect.TypeOf(testSlice))
	}

	// Struct test
	refMap := map[string]string{
		"field_name": "value foo",
	}
	cmpStruct := refStruct{
		FieldName: "value foo",
	}
	err = d.Set("testStruct", refMap)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testStruct': %+v", err.Error())
	}
	testStruct := FillValue(d.Get("testStruct"), refStruct{}).(refStruct)
	if reflect.DeepEqual(testStruct, cmpStruct) != true {
		t.Errorf("testStruct doesn't match cmpStruct\ntestStruct: %+v\ncmpStruct: %+v", testStruct, cmpStruct)
	}

	// Struct from slice test
	refSlice := []map[string]string{
		refMap,
	}
	err = d.Set("testStructSlice", refSlice)
	if err != nil {
		t.Errorf("Error setting resourceData for 'testStruct': %+v", err.Error())
	}
	testStruct = FillValue(d.Get("testStruct"), refStruct{}).(refStruct)
	if reflect.DeepEqual(testStruct, cmpStruct) != true {
		t.Errorf("testStruct doesn't match cmpStruct\ntestStruct: %+v\ncmpStruct: %+v", testStruct, cmpStruct)
	}

}

func TestUnderscoreToLowerCamelCase(t *testing.T) {
	s1 := UnderscoreToLowerCamelCase("a_test")
	if s1 != "aTest" {
		t.Errorf("String 'a_test' was not transformed to 'aTest', but instead '%s'", s1)
	}
	s2 := UnderscoreToLowerCamelCase("ateSt")
	if s2 != "atest" {
		t.Errorf("String 'aeSt' was not transformed to 'atest', but instead '%s'", s2)
	}
	s3 := UnderscoreToLowerCamelCase("_a_test")
	if s3 != "aTest" {
		t.Errorf("String '_atest' was not transformed to 'atest', but instead '%s'", s3)
	}
}

func TestCamelCaseToUnderscore(t *testing.T) {
	s1 := CamelCaseToUnderscore("CamelCase")
	if s1 != "camel_case" {
		t.Errorf("String 'CamelCase' was not transformed to 'camel_case', but instead '%s'", s1)
	}
	s2 := CamelCaseToUnderscore("camelCase")
	if s2 != "camel_case" {
		t.Errorf("String 'camelCase' was not transformed to 'camel_case', but instead '%s'", s2)
	}
	s3 := CamelCaseToUnderscore("camelCase")
	if s3 != "camel_case" {
		t.Errorf("String '_camelCase' was not transformed to 'camel_case', but instead '%s'", s3)
	}
}

func TestGetJSONKey(t *testing.T) {
	type testStruct struct {
		TestField int `json:"testField"`
	}
	ts := testStruct{}
	v := reflect.ValueOf(&ts).Elem()
	vt := reflect.TypeOf(v.Interface())
	tag := GetJSONKey(vt.Field(0))
	if tag != "testField" {
		t.Errorf("Field name should be 'testField', but received '%s'", tag)
	}
}
