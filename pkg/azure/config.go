package azure

import (
	"context"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	oauth2cc "golang.org/x/oauth2/clientcredentials"
	"gopkg.in/yaml.v2"
)

const (
	// Resource indices refers to the fixed order of different
	// microsoft/azure domains added under azure.resources in services.yml.
	// Order of the resources is critical and should not be changed.
	// Note: Declare new constant variables after the resource indices.
	partnerResourceIndex       = iota
	managementResourceIndex    = iota
	graphResourceIndex         = iota
	partnerCenterResourceIndex = iota

	// ProviderShortName refers to the shortname field of the provider table for Azure
	ProviderShortName = "MSAZ"
	// AddMode is the name for mode operation add
	AddMode = "add"
	// RemoveMode is the name for mode operation remove
	RemoveMode = "remove"
)

// AZConfig contains azure service credentials
type AZConfig struct {
	ClientID       string                 `json:"clientID" yaml:"clientID"`
	ClientSecret   string                 `json:"clientSecret" yaml:"clientSecret"`
	TokenURL       string                 `json:"tokenURL" yaml:"tokenURL"`
	Resources      []string               `json:"resources" yaml:"resources"`
	Publisher      string                 `json:"publisher" yaml:"publisher"`
	TestOfferName  string                 `json:"testOfferName" yaml:"testOfferName"`
	IdentityCAFile string                 `json:"identityCAFile" yaml:"identityCAFile"`
	Offers         map[string]OfferConfig `json:"offers" yaml:"offers"`
}

// OfferConfig represents the offer config offer_name and product_durable_id
type OfferConfig struct {
	ProductDurableID string `json:"productDurableID" yaml:"productDurableID"`
}

type ServiceCollection map[string]AZService

// AZService implements service.Service interface along with additional details required for azure
type AZService struct {
	BaseURL         string
	Publisher       string
	IdentityCACerts *x509.CertPool
	httpClient      *http.Client
}

// Config is package var to store azure service credentials from yaml
var Config = AZConfig{
	ClientID:     os.Getenv("AZURE_CLIENT_ID"),
	ClientSecret: os.Getenv("AZURE_CLIENT_SECRET"),
	TokenURL:     fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", os.Getenv("AZURE_TENANT_ID")),
	Resources: []string{
		"https://cloudpartner.azure.com",
		"https://management.azure.com",
		"https://graph.microsoft.com",
	},
	Publisher:      os.Getenv("AZURE_OFFERS_PUBLISHER"),
	IdentityCAFile: os.Getenv("AZURE_CERT_FILE"),
}

// Instance contains the default AZService instance, initialized by main()
var Instance *AZService

// MgmtResourceInstance contains the AZService instance to call management APIs of Azure
var MgmtResourceInstance *AZService

// MgmtResourceInstance contains the AZService instance to call graph APIs of Microsoft
var GraphResourceInstance *AZService

// GetBaseURL returns base url
func (svc AZService) GetBaseURL() string {
	return svc.BaseURL
}

// GetClient returns oauth http client for interacting with azure apis
func (svc AZService) GetClient() (*http.Client, error) {
	return svc.httpClient, nil
}

// NewService creates a new service from the AZConfig
func (c *AZConfig) NewService(resourceIndex int) (*AZService, error) {
	if resourceIndex >= len(c.Resources) {
		return nil, fmt.Errorf("resourceIndex %v is exceeding available number of resources %v", resourceIndex, len(c.Resources))
	}

	oauthConf := &oauth2cc.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.TokenURL,
		EndpointParams: map[string][]string{
			"resource": {c.Resources[resourceIndex]},
		},
		AuthStyle: oauth2.AuthStyleInParams,
	}

	certs, err := loadCACerts(Config.IdentityCAFile)
	if err != nil {
		return nil, fmt.Errorf("unable to load Azure identity CA certs: %w", err)
	}

	return &AZService{
		BaseURL:         c.Resources[resourceIndex],
		Publisher:       c.Publisher,
		httpClient:      oauthConf.Client(context.TODO()),
		IdentityCACerts: certs,
	}, nil
}

// IsTestOfferNamePresent checks if test offer name is present
// Note: Empty test offer name implies that prod offers will be used
func (c *AZConfig) IsTestOfferNamePresent() bool {
	return c.TestOfferName != ""
}

// StartService starts the Azure service
func StartService() error {
	var err error

	offers, err := os.ReadFile(os.Getenv("AZURE_OFFERS_FILE"))
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(offers, &Config.Offers)
	if err != nil {
		return err
	}

	Instance, err = Config.NewService(partnerResourceIndex)
	if err != nil {
		return err
	}

	MgmtResourceInstance, err = Config.NewService(managementResourceIndex)
	if err != nil {
		return err
	}

	GraphResourceInstance, err = Config.NewService(graphResourceIndex)
	if err != nil {
		return err
	}

	return err
}

// loadCACerts loads the CA certs from the given path
func loadCACerts(path string) (*x509.CertPool, error) {
	if len(path) == 0 {
		return nil, nil
	}
	var err error
	pem, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(pem)
	if !ok {
		err = errors.New("unable to append certificates from PEM file")
	}
	return pool, err
}
