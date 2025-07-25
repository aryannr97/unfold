package google

import (
	"log"
	"os"
)

func prepareTestEnvironment() {
	os.Setenv("GOOGLE_KEYFILE", "testdata/google-keyfile.json")
	Config.ServiceAccountKeyFile = os.Getenv("GOOGLE_KEYFILE")
	err := StartService()
	if err != nil {
		log.Fatalf("failed to start service: %v", err)
	}
}
