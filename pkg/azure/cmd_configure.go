package azure

import (
	"flag"
	"fmt"

	"github.com/aryannr97/unfold/pkg/commands"
)

// commandConfigureConfig represents the configuration for the configure command
type commandConfigureConfig struct {
	FlagSet       *flag.FlagSet
	AddRemoveOpts struct {
		RemoveFlag     *bool
		SubscriptionID *string
		TenantID       *string
		Offer          *string
	}
}

// Execute executes the configure command
func (c commandConfigureConfig) Execute() string {
	var resource, resourceID string
	mode := "add"

	if *c.AddRemoveOpts.RemoveFlag {
		mode = "remove"
	}
	if c.IsSub() {
		resource = "sub"
		resourceID = *c.AddRemoveOpts.SubscriptionID
	} else {
		resource = "tenant"
		resourceID = *c.AddRemoveOpts.TenantID
	}
	return fmt.Sprintf("[unfold] %v", MakeConfigurationRequest(*c.AddRemoveOpts.Offer, resourceID, resource, mode))
}

// GetFlagSet returns the flag set for the configure command
func (c commandConfigureConfig) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

// IsSub checks if the subscription ID is provided
func (c *commandConfigureConfig) IsSub() bool {
	return c.AddRemoveOpts.SubscriptionID != nil && *c.AddRemoveOpts.SubscriptionID != ""
}

// fetchCommandConfigureConfig fetches the command configure config
func fetchCommandConfigureConfig() commandConfigureConfig {
	flagSet := flag.NewFlagSet(commands.Configure, flag.ContinueOnError)
	return commandConfigureConfig{
		AddRemoveOpts: struct {
			RemoveFlag     *bool
			SubscriptionID *string
			TenantID       *string
			Offer          *string
		}{
			RemoveFlag:     flagSet.Bool("r", false, "remove resource from respective private audience"),
			SubscriptionID: flagSet.String("sid", "", "provide a valid azure subscription id"),
			TenantID:       flagSet.String("tid", "", "provide a valid azure tenant id"),
			Offer:          flagSet.String("o", "", "provide a valid azure offer name"),
		},
		FlagSet: flagSet,
	}
}
