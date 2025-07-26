package google

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aryannr97/unfold/pkg/helpers"
	"golang.org/x/net/context"
	cloudidentity "google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/option"
)

func Test_commandGetConfig_Execute(t *testing.T) {
	prepareTestEnvironment()
	tests := []struct {
		name          string
		args          []string
		transport     map[string]*http.Response
		httpCallError error
		want          string
	}{
		{
			name: "test get group by id",
			args: []string{"-g", "test-group"},
			transport: map[string]*http.Response{
				"/v1/groups:lookup": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"name": "test-group"}`)),
				},
			},
			httpCallError: nil,
			want:          "retrieved group resource name " + helpers.GreenValue("test-group"),
		},
		{
			name:          "test get group by id flag not provided",
			args:          []string{},
			transport:     nil,
			httpCallError: nil,
			want:          "something went wrong",
		},
		{
			name:          "test get group by id with error",
			args:          []string{"-g", "test-group"},
			transport:     nil,
			httpCallError: errors.New("test error"),
			want:          "test error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandGetConfig
			c.GetFlagSet().Parse(tt.args)
			ctx := context.WithValue(context.Background(), "test", tt.name)
			instance.CloudIdentityService, _ = cloudidentity.NewService(ctx,
				option.WithHTTPClient(&http.Client{
					Transport: &MockHTTPRoundTripper{
						Transport:    tt.transport,
						Error:        tt.httpCallError,
						ErrorOnIndex: 0,
						CurrentIndex: 0,
					},
				}),
			)
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandGetConfig.Execute() = %v, want %v", got, tt.want)
			}
			instance.Groups = make(map[string]*cloudidentity.LookupGroupNameResponse)
		})
	}
}

type MockHTTPRoundTripper struct {
	Transport    map[string]*http.Response
	Error        error
	ErrorOnIndex int
	CurrentIndex int
}

func (m *MockHTTPRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.Error != nil && m.ErrorOnIndex == m.CurrentIndex {
		return nil, m.Error
	}

	m.CurrentIndex++

	// Use the full URL as the key for lookup
	url := req.URL.String()
	if resp, ok := m.Transport[url]; ok {
		return resp, nil
	}

	// Use the relative URL as the key for lookup
	url = req.URL.Path
	if resp, ok := m.Transport[url]; ok {
		return resp, nil
	}

	// Return a default response if no match found
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
	}, nil
}
