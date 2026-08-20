package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestNormalizeAPIResponseRequests(t *testing.T) {
	testCases := []struct {
		name          string
		reader        rawConfigReader
		requests      []tests.ApiRequest
		wantHeaders   [][]string
		wantVariables [][]string
	}{
		{
			name:   "import removes only exact generated artifacts",
			reader: &ctyRawConfigReader{values: map[string]cty.Value{}},
			requests: []tests.ApiRequest{
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_1}}"),
					apiHeader("authorization", "{{OAuth2_Token_Step_1}}"),
					apiHeader("Content-Type", "application/json"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_1", ""),
					apiVariable("OAuth2_Token_Step_1", "$.token"),
					apiVariable("OAuth2_Token_Step_10", ""),
				}),
			},
			wantHeaders: [][]string{{
				"authorization={{OAuth2_Token_Step_1}}",
				"Content-Type=application/json",
			}},
			wantVariables: [][]string{{
				"OAuth2_Token_Step_1=$.token",
				"OAuth2_Token_Step_10=",
			}},
		},
		{
			name:   "uses execution step position across oauth2 requests",
			reader: &ctyRawConfigReader{values: map[string]cty.Value{}},
			requests: []tests.ApiRequest{
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_1}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_1", ""),
				}),
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_2}}"),
					apiHeader("Authorization", "{{OAuth2_Token_Step_3}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_2", ""),
					apiVariable("OAuth2_Token_Step_3", ""),
				}),
			},
			wantHeaders:   [][]string{nil, {"Authorization={{OAuth2_Token_Step_2}}"}},
			wantVariables: [][]string{nil, {"OAuth2_Token_Step_2="}},
		},
		{
			name:   "counts non-oauth2 requests as one execution step",
			reader: &ctyRawConfigReader{values: map[string]cty.Value{}},
			requests: []tests.ApiRequest{
				apiRequest("none", nil, nil),
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_2}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_2", ""),
				}),
			},
			wantHeaders:   [][]string{nil, nil},
			wantVariables: [][]string{nil, nil},
		},
		{
			name: "preserves explicitly configured exact artifacts",
			reader: apiRequestsRawConfigReader(
				apiConfiguredRequest(
					[]cty.Value{apiConfiguredObject("key", "Authorization", "value", "{{OAuth2_Token_Step_1}}")},
					[]cty.Value{apiConfiguredObject("name", "OAuth2_Token_Step_1", "value", "")},
				),
			),
			requests: []tests.ApiRequest{
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_1}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_1", ""),
				}),
			},
			wantHeaders:   [][]string{{"Authorization={{OAuth2_Token_Step_1}}"}},
			wantVariables: [][]string{{"OAuth2_Token_Step_1="}},
		},
		{
			name:   "does not normalize non-oauth2 requests",
			reader: &ctyRawConfigReader{values: map[string]cty.Value{}},
			requests: []tests.ApiRequest{
				apiRequest("basic", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_1}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_1", ""),
				}),
			},
			wantHeaders:   [][]string{{"Authorization={{OAuth2_Token_Step_1}}"}},
			wantVariables: [][]string{{"OAuth2_Token_Step_1="}},
		},
		{
			name: "preserves artifacts when configured values are unknown",
			reader: apiRequestsRawConfigReader(cty.ObjectVal(map[string]cty.Value{
				"headers":   cty.DynamicVal,
				"variables": cty.DynamicVal,
			})),
			requests: []tests.ApiRequest{
				apiRequest("oauth2", []tests.ApiRequestHeader{
					apiHeader("Authorization", "{{OAuth2_Token_Step_1}}"),
				}, []tests.ApiRequestVariable{
					apiVariable("OAuth2_Token_Step_1", ""),
				}),
			},
			wantHeaders:   [][]string{{"Authorization={{OAuth2_Token_Step_1}}"}},
			wantVariables: [][]string{{"OAuth2_Token_Step_1="}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := &tests.ApiTestResponse{Requests: tc.requests}
			normalizeAPIResponseRequests(tc.reader, response)

			gotHeaders, gotVariables := apiRequestValues(response.Requests)
			if !reflect.DeepEqual(gotHeaders, tc.wantHeaders) {
				t.Fatalf("unexpected headers: got %#v want %#v", gotHeaders, tc.wantHeaders)
			}
			if !reflect.DeepEqual(gotVariables, tc.wantVariables) {
				t.Fatalf("unexpected variables: got %#v want %#v", gotVariables, tc.wantVariables)
			}
		})
	}
}

func TestConfiguredAPIRequestObjectMatchesKnownMismatchTakesPrecedence(t *testing.T) {
	item := cty.ObjectVal(map[string]cty.Value{
		"key":   cty.StringVal("X-Other-Header"),
		"value": cty.DynamicVal,
	})

	matches, known := configuredAPIRequestObjectMatches(item, map[string]string{
		"key":   "Authorization",
		"value": "{{OAuth2_Token_Step_1}}",
	})
	if matches || !known {
		t.Fatalf("expected a known non-match, got matches=%t known=%t", matches, known)
	}
}

func apiRequest(authType string, headers []tests.ApiRequestHeader, variables []tests.ApiRequestVariable) tests.ApiRequest {
	typedAuthType := tests.ApiRequestAuthType(authType)
	return tests.ApiRequest{
		AuthType:  &typedAuthType,
		Headers:   headers,
		Variables: variables,
	}
}

func apiHeader(key, value string) tests.ApiRequestHeader {
	return tests.ApiRequestHeader{Key: &key, Value: &value}
}

func apiVariable(name, value string) tests.ApiRequestVariable {
	return tests.ApiRequestVariable{Name: &name, Value: &value}
}

func apiRequestsRawConfigReader(requests ...cty.Value) rawConfigReader {
	return &ctyRawConfigReader{values: map[string]cty.Value{
		"requests": cty.TupleVal(requests),
	}}
}

func apiConfiguredRequest(headers, variables []cty.Value) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		"headers":   apiConfiguredCollection(headers),
		"variables": apiConfiguredCollection(variables),
	})
}

func apiConfiguredCollection(values []cty.Value) cty.Value {
	if len(values) == 0 {
		return cty.EmptyTupleVal
	}
	return cty.TupleVal(values)
}

func apiConfiguredObject(firstKey, firstValue, secondKey, secondValue string) cty.Value {
	return cty.ObjectVal(map[string]cty.Value{
		firstKey:  cty.StringVal(firstValue),
		secondKey: cty.StringVal(secondValue),
	})
}

func apiRequestValues(requests []tests.ApiRequest) ([][]string, [][]string) {
	headers := make([][]string, len(requests))
	variables := make([][]string, len(requests))
	for i, request := range requests {
		for _, header := range request.Headers {
			headers[i] = append(headers[i], *header.Key+"="+*header.Value)
		}
		for _, variable := range request.Variables {
			variables[i] = append(variables[i], *variable.Name+"="+*variable.Value)
		}
	}
	return headers, variables
}
