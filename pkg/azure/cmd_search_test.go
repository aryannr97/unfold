package azure

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aryannr97/unfold/pkg/helpers"
)

func Test_commandSearchConfig_Execute(t *testing.T) {
	prepareTestEnvironment()

	tests := []struct {
		name          string
		args          []string
		transport     map[string]*http.Response
		httpCallError error
		want          string
	}{
		{
			name: "search by id",
			args: []string{"-id", "12345678-1234-1234-1234-123456789abc", "-o", "offer-1"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/resource-tree/product/12345678-1234-1234-1234-123456789abc": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"resources": [{"privateAudiences": [{"id": "12345678-1234-1234-1234-123456789abc", "type": "subscription"}]}]}`)),
				},
			},
			httpCallError: nil,
			want:          fmt.Sprintf("found %s in private audience with type %s", helpers.GreenValue("12345678-1234-1234-1234-123456789abc"), helpers.GreenValue("subscription")),
		},
		{
			name: "search by id not found",
			args: []string{"-id", "12345678-1234-1234-1234-123456789abc", "-o", "offer-1"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/resource-tree/product/12345678-1234-1234-1234-123456789abc": {
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(`{"resources": [{"privateAudiences": [{"id": "12345678-1234-1234-1234-123456789abd", "type": "subscription"}]}]}`)),
				},
			},
			httpCallError: nil,
			want:          fmt.Sprintf("given id %s in private audience", helpers.RedValue("not found")),
		},
		{
			name: "search by id http call failure",
			args: []string{"-id", "12345678-1234-1234-1234-123456789abc", "-o", "offer-1"},
			transport: map[string]*http.Response{
				"https://graph.microsoft.com/rp/product-ingestion/resource-tree/product/12345678-1234-1234-1234-123456789abc": {
					StatusCode: http.StatusNotFound,
					Body:       io.NopCloser(bytes.NewBufferString(`{"error": "not found"}`)),
				},
			},
			httpCallError: nil,
			want:          "marketplace returned 404 with response",
		},
		{
			name:          "search by id http call error",
			args:          []string{"-id", "12345678-1234-1234-1234-123456789abc", "-o", "offer-1"},
			transport:     nil,
			httpCallError: errors.New("http call error"),
			want:          "http call error",
		},
		{
			name:          "prodductDurableID is empty",
			args:          []string{"-id", "12345678-1234-1234-1234-123456789abc"},
			transport:     nil,
			httpCallError: nil,
			want:          "offer cannot be empty",
		},
		{
			name:          "id is empty",
			args:          []string{"-o", "offer-1"},
			transport:     nil,
			httpCallError: nil,
			want:          "id cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCommandModule().CommandSearchConfig
			c.GetFlagSet().Parse(tt.args)
			instances[graphResourceIndex].httpClient = &http.Client{
				Transport: &MockHTTPRoundTripper{
					Transport: tt.transport,
					Error:     tt.httpCallError,
				},
			}
			if got := c.Execute(); !strings.Contains(got, tt.want) {
				t.Errorf("commandSearchConfig.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
