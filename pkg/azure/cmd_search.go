package azure

import (
	"flag"
	"fmt"

	"github.com/aryannr97/unfold/pkg/commands"
)

// commandSearchConfig represents the configuration for the search command
type commandSearchConfig struct {
	AudienceOpts struct {
		ID    *string
		Offer *string
	}
	FlagSet *flag.FlagSet
}

// Execute executes the search command
func (c commandSearchConfig) Execute() string {
	if *c.AudienceOpts.ID != "" {
		return fmt.Sprintf("[unfold] %s", Search(*c.AudienceOpts.ID, Config.Offers[*c.AudienceOpts.Offer].ProductDurableID))
	}
	return "[unfold] id cannot be empty"
}

// GetFlagSet returns the flag set for the search command
func (c commandSearchConfig) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

// fetchCommandSearchConfig fetches the command search config
func fetchCommandSearchConfig() commandSearchConfig {
	flagSet := flag.NewFlagSet(commands.Search, flag.ContinueOnError)
	return commandSearchConfig{
		AudienceOpts: struct {
			ID    *string
			Offer *string
		}{
			ID:    flagSet.String("id", "", "used to search resource in private audience of offer"),
			Offer: flagSet.String("o", "", "provide a valid azure offer name"),
		},
		FlagSet: flagSet,
	}
}
