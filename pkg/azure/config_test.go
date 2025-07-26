package azure

import (
	"os"
	"testing"
)

func TestStartService(t *testing.T) {
	tests := []struct {
		name         string
		modifyConfig func()
		wantErr      bool
	}{
		{
			name:    "start service success",
			wantErr: false,
		},
		{
			name: "start service error no offers file",
			modifyConfig: func() {
				os.Setenv("AZURE_OFFERS_FILE", "offers_test.yml")
			},
			wantErr: true,
		},
		{
			name: "start service error invalid offers file",
			modifyConfig: func() {
				os.Setenv("AZURE_OFFERS_FILE", "testdata/invalid_offers_test.yml")
			},
			wantErr: true,
		},
		{
			name: "start service error cert file not found",
			modifyConfig: func() {
				os.Setenv("AZURE_CERT_FILE", "testdata/ca_cert.txt")
			},
			wantErr: true,
		},
		{
			name: "start service error invalid cert file",
			modifyConfig: func() {
				os.Setenv("AZURE_CERT_FILE", "testdata/invalid_ca_cert.txt")
			},
			wantErr: true,
		},
		{
			name: "start service error resource index out of bounds",
			modifyConfig: func() {
				// add a new resource index in config, as iota cannot be modified
				instances[len(config.Resources)+1] = nil
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prepareTestEnvironment()
			if tt.modifyConfig != nil {
				tt.modifyConfig()
			}
			if err := StartService(); (err != nil) != tt.wantErr {
				t.Errorf("StartService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
