package azure

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test_commandConfigureConfig_Execute(t *testing.T) {
	prepareTestEnvironment()
	tests := []struct {
		name          string
		args          []string
		transport     map[string]*http.Response
		errorOnIndex  int
		httpCallError error
		want          string
	}{
		{
			name: "add subscription success",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobId": "12345678-1234-1234-1234-123456789def", "jobStatus": "pending", "jobResult": "pending", "errors": []}`)),
				},
			},
			httpCallError: nil,
			want:          "configure response",
		},
		{
			name: "remove subscription success",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2", "-r"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobId": "12345678-1234-1234-1234-123456789def", "jobStatus": "pending", "jobResult": "pending", "errors": []}`)),
				},
			},
			httpCallError: nil,
			want:          "configure response",
		},
		{
			name: "add tenant success",
			args: []string{"-tid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobId": "12345678-1234-1234-1234-123456789def", "jobStatus": "pending", "jobResult": "pending", "errors": []}`)),
				},
			},
			httpCallError: nil,
			want:          "configure response",
		},
		{
			name: "remove tenant success",
			args: []string{"-tid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2", "-r"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"jobId": "12345678-1234-1234-1234-123456789def", "jobStatus": "pending", "jobResult": "pending", "errors": []}`)),
				},
			},
			httpCallError: nil,
			want:          "configure response",
		},
		{
			name: "add resource failed to get plans",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": {"code": "InternalServerError", "message": "Internal server error"}}`)),
				},
			},
			httpCallError: nil,
			want:          "marketplace returned 500 for getPlans with response",
		},
		{
			name: "add resource get plans invalid response",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{`)),
				},
			},
			httpCallError: nil,
			want:          "json decode",
		},
		{
			name:          "add resource http call error",
			args:          []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport:     nil,
			errorOnIndex:  0,
			httpCallError: errors.New("http call error"),
			want:          "http call error",
		},
		{
			name: "add subscription failed to configure",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": {"code": "InternalServerError", "message": "Internal server error"}}`)),
				},
			},
			httpCallError: nil,
			want:          "marketplace returned 500 for configurePrivateAudienceAPI with response",
		},
		{
			name: "add subscription failed to configure invalid response",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
				"https://graph.microsoft.com/rp/product-ingestion/configure?$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{`)),
				},
			},
			httpCallError: nil,
			want:          "json decode",
		},
		{
			name: "add subscription configure http call error",
			args: []string{"-sid", "12345678-1234-1234-1234-123456789abc", "-o", "offer-2"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/plan?product=product/87654321-4321-4321-4321-210987654321&$version=2022-03-01-preview2": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"value": [{"$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#", "id": "12345678-1234-1234-1234-12345678plan-1"}]}`)),
				},
			},
			errorOnIndex:  1,
			httpCallError: errors.New("http call error"),
			want:          "http call error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandConfigureConfig
			c.GetFlagSet().Parse(tt.args)
			instances[graphResourceIndex].httpClient = &http.Client{
				Transport: &MockHTTPRoundTripper{
					Transport:    tt.transport,
					Error:        tt.httpCallError,
					ErrorOnIndex: tt.errorOnIndex,
				},
			}
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandConfigureConfig.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
