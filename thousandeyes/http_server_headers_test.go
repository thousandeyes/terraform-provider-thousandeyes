package thousandeyes

import (
	"reflect"
	"testing"

	"github.com/thousandeyes/thousandeyes-sdk-go/v3/tests"
)

func TestSyncHTTPServerHeaders_mergesHeadersAndCustomRoot(t *testing.T) {
	root := map[string]string{
		"X-TE-Test": "upgrade-full",
	}

	req := &tests.HttpServerTestRequest{
		Headers: []string{
			"Accept: application/json",
			"Content-Type: application/json",
		},
		CustomHeaders: &tests.TestCustomHeaders{
			Root: &root,
		},
	}

	syncHTTPServerHeaders(req)

	wantHeaders := []string{
		"Accept: application/json",
		"Content-Type: application/json",
		"X-TE-Test: upgrade-full",
	}

	if !reflect.DeepEqual(req.Headers, wantHeaders) {
		t.Fatalf("unexpected headers: got %#v want %#v", req.Headers, wantHeaders)
	}

	wantRoot := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
		"X-TE-Test":    "upgrade-full",
	}

	if req.CustomHeaders == nil || req.CustomHeaders.Root == nil {
		t.Fatalf("expected custom_headers.root to be populated")
	}

	if !reflect.DeepEqual(*req.CustomHeaders.Root, wantRoot) {
		t.Fatalf("unexpected custom_headers.root: got %#v want %#v", *req.CustomHeaders.Root, wantRoot)
	}
}
