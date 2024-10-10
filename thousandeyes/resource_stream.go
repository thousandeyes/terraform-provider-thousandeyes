package thousandeyes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/thousandeyes/thousandeyes-sdk-go/v2"
)

type Stream struct {
	ID                string            `json:"id,omitempty"`
	Enabled           bool              `json:"enabled,omitempty"`
	Type              string            `json:"type,omitempty"`
	EndpointType      string            `json:"endpointType,omitempty"`
	StreamEndpointUrl string            `json:"streamEndpointUrl,omitempty"`
	DataModelVersion  string            `json:"dataModelVersion,omitempty"`
	TestMatch         []StreamTestMatch `json:"testMatch,omitempty"`
	TagMatch          []StreamTagMatch  `json:"tagMatch,omitempty"`
}

type StreamTestMatch struct {
	ID     string `json:"id,omitempty"`
	Domain string `json:"domain,omitempty"`
}

type StreamTagMatch struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// StreamClient extends thousandeyes.Client such that we can use the (at the
// time of writing) not-yet-implemented ThousandEyes v7 API.
type StreamClient struct {
	c *thousandeyes.Client
}

// NewStreamClientFrom copies an existing v6 client, modifies the URL in its
// APIEndpoint field to point to v7, and returns a StreamClient with the
// modified copy.
func NewStreamClientFrom(v6 *thousandeyes.Client) *StreamClient {
	v7 := new(thousandeyes.Client)
	*v7 = *v6
	v7.APIEndpoint = strings.ReplaceAll(v6.APIEndpoint, "v6", "v7")
	return &StreamClient{
		c: v7,
	}
}

// Exposing a simpler version of do() because thousandeyes.Client won't.
func (sc *StreamClient) do(method string, path string, payload interface{}) (*http.Response, error) {
	if sc.c.Limiter != nil {
		sc.c.Limiter.Wait()
	}
	endpoint := sc.c.APIEndpoint + path + ".json"
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(data))
	if sc.c.AccountGroupID != "" {
		q := req.URL.Query()
		q.Add("aid", sc.c.AccountGroupID)
		req.URL.RawQuery = q.Encode()
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", sc.c.AuthToken))
	req.Header.Set("content-type", "application/json")
	req.Header.Set("user-agent", sc.c.UserAgent)

	resp, err := sc.c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	if 199 >= resp.StatusCode || 300 <= resp.StatusCode {
		return nil, fmt.Errorf("Failed call API endpoint. HTTP response code: %v.", resp.StatusCode)
	}
	return resp, nil
}

// We need to expose decodeJSON() because thousandeyes.Client won't.
func (sc *StreamClient) decodeJSON(resp *http.Response, payload interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(payload)
}

func (sc *StreamClient) CreateStream(s Stream) (*Stream, error) {
	resp, err := sc.do("POST", "/stream", s)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to create stream, response code %d", resp.StatusCode)
	}
	var target Stream
	if dErr := sc.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return &target, nil
}

func (sc *StreamClient) GetStream(id string) (*Stream, error) {
	resp, err := sc.do("GET", fmt.Sprintf("/stream/%s", id), nil)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get stream, response code %d", resp.StatusCode)
	}
	var target Stream
	if dErr := sc.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return &target, nil
}

func (sc *StreamClient) UpdateStream(id string, s Stream) (*Stream, error) {
	resp, err := sc.do("PUT", fmt.Sprintf("/stream/%d", id), s)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to update stream, response code %d", resp.StatusCode)
	}
	var target Stream
	if dErr := sc.decodeJSON(resp, &target); dErr != nil {
		return nil, fmt.Errorf("Could not decode JSON response: %v", dErr)
	}
	return &target, nil
}

func (sc *StreamClient) DeleteStream(id string) error {
	resp, err := sc.do("DELETE", fmt.Sprintf("/stream/%s", id), nil)
	if err != nil {
		return err
	}
	if resp.StatusCode != 204 {
		return fmt.Errorf("failed to delete stream, response code %d", resp.StatusCode)
	}
	return nil
}

func resourceStream() *schema.Resource {
	resource := schema.Resource{
		Schema: ResourceSchemaBuild(Stream{}, schemas, nil),
		Create: resourceStreamCreate,
		Read:   resourceStreamRead,
		Update: resourceStreamUpdate,
		Delete: resourceStreamDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Description: "This resource allows you to create an OpenTelemetry data stream. For more information, see [Streams](https://developer.cisco.com/docs/thousandeyes/list-data-streams/).",
	}
	return &resource
}

func resourceStreamCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	streamClient := NewStreamClientFrom(client)
	log.Printf("[INFO] Creating ThousandEyes Stream %s", d.Id())
	local := buildStreamStruct(d)

	remote, err := streamClient.CreateStream(*local)
	if err != nil {
		return err
	}
	d.SetId(remote.ID)
	return resourceStreamRead(d, m)
}

func resourceStreamUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	streamClient := NewStreamClientFrom(client)

	log.Printf("[INFO] Updating ThousandEyes Stream %s", d.Id())
	update := ResourceUpdate(d, &Stream{}).(*Stream)
	_, err := streamClient.UpdateStream(d.Id(), *update)
	if err != nil {
		return err
	}
	return resourceStreamRead(d, m)
}

func resourceStreamRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	streamClient := NewStreamClientFrom(client)

	log.Printf("[INFO] Reading Thousandeyes Stream %s", d.Id())
	remote, err := streamClient.GetStream(d.Id())
	if err != nil {
		d.SetId("") // Set ID to empty to mark the resource as non-existent
		return err
	}
	return ResourceRead(d, remote)
}

func resourceStreamDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*thousandeyes.Client)
	streamClient := NewStreamClientFrom(client)

	log.Printf("[INFO] Deleting ThousandEyes Stream %s", d.Id())
	if err := streamClient.DeleteStream(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func buildStreamStruct(d *schema.ResourceData) *Stream {
	return ResourceBuildStruct(d, &Stream{}).(*Stream)
}
