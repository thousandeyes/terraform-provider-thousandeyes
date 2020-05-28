package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/william20111/go-thousandeyes"
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
	cmpStruct := thousandeyes.BGP{
		Prefix: prefix,
	}
	newStruct := thousandeyes.BGP{}
	attrs := map[string]string{
		"prefix": "8.19.2.2/19",
	}
	d := getReferenceData(schemas, attrs)
	ResourceBuildStruct(d, &newStruct)

	if newStruct.Prefix != cmpStruct.Prefix {
		t.Error("Building resource did not assign struct field correctly.")
	}

}

func TestResourceRead(t *testing.T) {
	prefix := "8.19.2.2/19"
	attrs := map[string]string{}
	d := getReferenceData(schemas, attrs)
	remoteResource := thousandeyes.BGP{
		Prefix: prefix,
	}
	err := ResourceRead(d, &remoteResource)
	if err != nil {
		t.Errorf("Setting resource data returned error: %+v", err.Error())
	}
	if d.Get("prefix") != remoteResource.Prefix {
		t.Errorf("Reading resource did not assign resource data correctly.\nStruct is %+v\nResource is %+v", remoteResource, d.State().Attributes)
	}
}

func TestResourceUpdate(t *testing.T) {

}

func TestResourceSchemaBuild(t *testing.T) {
	type refStruct struct {
		FieldName string `json:"fieldName"`
	}

	refSchema := map[string]*schema.Schema{
		"field_name": {
			Type: schema.TypeString,
		},
	}

	schm := ResourceSchemaBuild(refStruct{}, refSchema)

	for k, v := range refSchema {
		if _, ok := schm[k]; !ok {
			t.Errorf("Key %s missing from generated schema", k)
		}
		if reflect.DeepEqual(*v, *schm[k]) != true {
			t.Errorf("Schemas not equal: Reference schema is %+v\nNew Schema is %+v", refSchema, schm)
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
