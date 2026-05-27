package google

import (
	"os"
	"testing"

	"google.golang.org/api/option"
)

func TestStartService(t *testing.T) {
	tests := []struct {
		name         string
		modifyConfig func()
		wantErr      bool
	}{
		{
			name: "test start service success",
			modifyConfig: func() {
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile.json")
			},
			wantErr: false,
		},
		{
			name: "test start service error file not found",
			modifyConfig: func() {
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile-not-found.json")
			},
			wantErr: true,
		},
		{
			name: "test start service error invalid keyfile format",
			modifyConfig: func() {
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile-invalid-format.json")
			},
			wantErr: true,
		},
		{
			name: "test start service error cloud identity service creation",
			modifyConfig: func() {
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile.json")
				Config.clientOpts = func() ([]option.ClientOption, error) {
					// this will fail because both options are not allowed together
					return []option.ClientOption{option.WithScopes("test-scope"), option.WithAudiences("test-audience")}, nil
				}
			},
			wantErr: true,
		},
		{
			name: "test start service success with https proxy",
			modifyConfig: func() {
				Config.clientOpts = prepareClientOpts
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile.json")
				os.Setenv("HTTPS_PROXY", "http://localhost:8200")
			},
			wantErr: false,
		},
		{
			name: "test start service failure with invalid https proxy",
			modifyConfig: func() {
				Config.clientOpts = prepareClientOpts
				os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile.json")
				os.Setenv("HTTPS_PROXY", "http://localhost:8200%/")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.modifyConfig()
			if err := StartService(); (err != nil) != tt.wantErr {
				t.Errorf("StartService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
