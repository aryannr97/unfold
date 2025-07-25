package azure

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_commandGetConfig_Execute(t *testing.T) {
	prepareTestEnvironment()

	tests := []struct {
		name          string
		args          []string
		transport     map[string]*http.Response
		httpCallError error
		resource      string
		want          string
	}{
		{
			name: "get single tenant by subscription id",
			args: []string{"-t", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://management.azure.com/subscriptions/12345678-1234-1234-1234-123456789abc?api-version=2022-12-01": {
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"code":"InvalidAuthenticationTokenTenant","message":"The access token is from the wrong issuer. It must match the tenant 'https://sts.windows.net/0b37927b-359e-4a60-8aac-67f88409ac5a/' associated with this subscription."}}`)),
				},
			},
			httpCallError: nil,
			resource:      "management",
			want:          "retrieved tenant(s): 0b37927b-359e-4a60-8aac-67f88409ac5a",
		},
		{
			name: "get multiple tenants	 by subscription id",
			args: []string{"-t", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://management.azure.com/subscriptions/12345678-1234-1234-1234-123456789abc?api-version=2022-12-01": {
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"code":"InvalidAuthenticationTokenTenant","message":"The access token is from the wrong issuer. It must match the tenants 'https://sts.windows.net/0b37927b-359e-4a60-8aac-67f88409ac5a/,https://sts.windows.net/0b37927b-359e-4a60-8aac-67f88409ac5b/' associated with this subscription."}}`)),
				},
			},
			httpCallError: nil,
			resource:      "management",
			want:          "retrieved tenant(s): 0b37927b-359e-4a60-8aac-67f88409ac5a,0b37927b-359e-4a60-8aac-67f88409ac5b",
		},
		{
			name: "get no tenants by subscription id",
			args: []string{"-t", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://management.azure.com/subscriptions/12345678-1234-1234-1234-123456789abc?api-version=2022-12-01": {
					StatusCode: http.StatusUnauthorized,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error":{"code":"AuthorizationFailed","message":"The client '0b37927b-359e-4a60-8aac-67f88409ac5a' with object id '0b37927b-359e-4a60-8aac-67f88409ac5a' does not have authorization to perform action 'Microsoft.Authorization/roleAssignments/read' over scope '/subscriptions/12345678-1234-1234-1234-123456789abc' or the scope is invalid. If access was recently granted, please refresh your credentials."}}`)),
				},
			},
			httpCallError: nil,
			resource:      "management",
			want:          "retrieved tenant(s): ",
		},
		{
			name: "get tenant with invalid response",
			args: []string{"-t", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://management.azure.com/subscriptions/12345678-1234-1234-1234-123456789abc?api-version=2022-12-01": {
					StatusCode: http.StatusNoContent,
					Body:       nil,
				},
			},
			httpCallError: nil,
			resource:      "management",
			want:          "json decode",
		},
		{
			name: "get job status success",
			args: []string{"-s", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/configure/12345678-1234-1234-1234-123456789abc/status?$version=2022-07-01": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobResult":"success","jobId":"12345678-1234-1234-1234-123456789abc","jobStatus":"completed","errors":[]}`)),
				},
			},
			httpCallError: nil,
			resource:      "graph",
			want:          "job status response",
		},
		{
			name: "get job status failed",
			args: []string{"-s", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/configure/12345678-1234-1234-1234-123456789abc/status?$version=2022-07-01": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobResult":"failed","jobId":"12345678-1234-1234-1234-123456789abc","jobStatus":"completed","errors":[{"code":"","message":"job status error object"}]}`)),
				},
			},
			httpCallError: nil,
			resource:      "graph",
			want:          "job status error object",
		},
		{
			name: "get job status with invalid response",
			args: []string{"-s", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/configure/12345678-1234-1234-1234-123456789abc/status?$version=2022-07-01": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{`)),
				},
			},
			httpCallError: nil,
			resource:      "graph",
			want:          "json decode",
		},
		{
			name: "get job status not found",
			args: []string{"-s", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/configure/12345678-1234-1234-1234-123456789abc/status?$version=2022-07-01": {
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
				},
			},
			httpCallError: nil,
			resource:      "graph",
			want:          "marketplace returned 404 with response",
		},
		{
			name: "get job status http call error",
			args: []string{"-s", "12345678-1234-1234-1234-123456789abc"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/configure/12345678-1234-1234-1234-123456789abc/status?$version=2022-07-01": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{`)),
				},
			},
			httpCallError: errors.New("http call error"),
			resource:      "graph",
			want:          "http call error",
		},
		{
			name:          "required flags empty",
			args:          []string{},
			transport:     nil,
			httpCallError: nil,
			want:          "something went wrong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandGetConfig
			c.GetFlagSet().Parse(tt.args)
			switch tt.resource {
			case "management":
				instances[managementResourceIndex].httpClient = &http.Client{
					Transport: &MockHTTPRoundTripper{
						Transport: tt.transport,
						Error:     tt.httpCallError,
					},
				}
			case "graph":
				instances[graphResourceIndex].httpClient = &http.Client{
					Transport: &MockHTTPRoundTripper{
						Transport: tt.transport,
						Error:     tt.httpCallError,
					},
				}
			}
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandGetConfig.Execute() = %v, want %v", got, tt.want)
			}
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

	// Return a default response if no match found
	return &http.Response{
		StatusCode: http.StatusNotFound,
		Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
	}, nil
}
