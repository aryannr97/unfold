package azure

import (
	"log"
	"os"
)

func prepareTestEnvironment() {
	os.Setenv("AZURE_OFFERS_FILE", "testdata/offers_test.yml")
	os.Setenv("AZURE_CERT_FILE", "testdata/test_ca_cert.txt")
	config.IdentityCAFile = os.Getenv("AZURE_CERT_FILE")
	err := StartService()
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
}
