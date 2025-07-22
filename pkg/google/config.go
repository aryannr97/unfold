package google

import (
	"fmt"
	"net/mail"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/context"
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
	TestGroup             string `json:"testGroup" yaml:"testGroup"`

	// JwkURL is the URL where the JWK to validate the instance identity token can be found.
	JwkURL string `json:"jwkURL" yaml:"jwkURL"`
}

// Config contains the actual values from service.yml file
var Config = GCEConfig{
	ServiceAccountKeyFile: os.Getenv("GOOGLE_KEYFILE"),
	TestGroup:             "",
	JwkURL:                os.Getenv("GOOGLE_JWK_URL"),
}

// NewService creates a new cloudidentity.service from the GCEConfig.
func (c GCEConfig) NewService() (*Service, error) {
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

// Instance holds the cloudidentity.Service as constructed from credentials in services.yml file.
// It also has information of the Groups as fetched using getGroupByID() func.
var Instance *Service

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

// IsValidAccountString checks if the given acc string is a validprojectID or EmailID
func IsValidAccountString(acc string, isGoldImageReq bool) bool {
	// If acc is of valid email format, return true.
	_, err := mail.ParseAddress(acc)
	if err == nil {
		return true
	} else if isGoldImageReq {
		// if acc is not of emailID format and its GoldImageReq return false.
		return false
	}

	// else check if its of format matching the below regex
	idReg := regexp.MustCompile("^[a-zA-Z0-9-]*$")
	if acc == "" || !idReg.MatchString(acc) || len(acc) > 64 {
		return false
	}

	return true
}

// StartService starts the Google service
func StartService() error {
	var err error
	Instance, err = Config.NewService()
	return err
}
