package google

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	cloudidentity "google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/option"
)

// ProviderShortName is short_name for google compute engine provider in DB.
const ProviderShortName = "GCE"

// GCEConfig is used to fetch credetials and other information from service.yml file
// required for accessing google APIs
type GCEConfig struct {
	ServiceAccountKeyFile string `json:"serviceAccountKeyFile" yaml:"serviceAccountKeyFile"`
	// JwkURL is the URL where the JWK to validate the instance identity token can be found.
	JwkURL string `json:"jwkURL" yaml:"jwkURL"`
	// TokenSourceFunc is a function that returns a token source.
	// It is used to override the default token source for testing.
	tokenSource TokenSourceFunc
}

// Config contains the actual values from service.yml file
var Config = GCEConfig{
	ServiceAccountKeyFile: os.Getenv("GOOGLE_KEYFILE"),
	JwkURL:                os.Getenv("GOOGLE_JWK_URL"),
}

type TokenSourceFunc func() oauth2.TokenSource

// NewService creates a new cloudidentity.service from the GCEConfig.
func (c GCEConfig) NewService() (*Service, error) {
	Config.ServiceAccountKeyFile = os.Getenv("GOOGLE_KEYFILE")

	var err error
	var ctx = context.Background()

	// get the serviceAccountKeyFile
	jsonCredentials, err := os.ReadFile(Config.ServiceAccountKeyFile)
	if err != nil {
		return nil, err
	}

	// parse the serviceAccountKeyFile
	config, err := google.JWTConfigFromJSON(jsonCredentials,
		cloudidentity.CloudIdentityGroupsScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse service account key file to config: %v", err)
	}
	ts := config.TokenSource(ctx)

	// if tokenSourceFunc is set, use it to get the token source
	if Config.tokenSource != nil {
		ts = Config.tokenSource()
	}

	// Build cloud identity API client
	svc, err := cloudidentity.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve client: %v", err)
	}
	return &Service{CloudIdentityService: svc,
		Groups: map[string]*cloudidentity.LookupGroupNameResponse{}}, nil
}

// Service has a CloudIdentityService field used for accessing google APIs
// and Groups for storing lookup information about GoogleGroups
type Service struct {
	CloudIdentityService *cloudidentity.Service
	Groups               map[string]*cloudidentity.LookupGroupNameResponse
}

// instance holds the cloudidentity.Service as constructed from credentials in services.yml file.
// It also has information of the Groups as fetched using getGroupByID() func.
var instance *Service

// AddGroup adds Group by groupID and the corresponding cloudidentity.LookupGroupNameResponse
// to the Groups field of Service.
func (s *Service) AddGroup(groupID string, groupResp *cloudidentity.LookupGroupNameResponse) {
	if _, ok := s.Groups[strings.ToLower(groupID)]; !ok {
		s.Groups[strings.ToLower(groupID)] = groupResp
	}
}

// GetGroup retrieves cloudidentity.LookupGroupNameResponse from the Groups field of
// Service by the groupID.
func (s *Service) GetGroup(groupID string) *cloudidentity.LookupGroupNameResponse {
	if r, ok := s.Groups[strings.ToLower(groupID)]; ok {
		return r
	}
	return nil
}

// StartService starts the Google service
func StartService() error {
	var err error
	instance, err = Config.NewService()
	return err
}
