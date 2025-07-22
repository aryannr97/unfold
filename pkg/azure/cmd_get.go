package azure

import (
	"flag"
	"fmt"
	"strings"

	"github.com/aryannr97/unfold/pkg/commands"
)

// commandGetConfig represents the configuration for the get command
type commandGetConfig struct {
	FlagSet *flag.FlagSet
	Opts    struct {
		TenantFlag *string
		StatusFlag *string
	}
}

// Execute executes the get command
func (c commandGetConfig) Execute() string {
	if *c.Opts.TenantFlag != "" {
		atf := NewTenantFinder()
		tenants, err := atf.GetTenantBySubscriptionID(*c.Opts.TenantFlag)
		if err != nil {
			return fmt.Sprintf("[unfold] %s", err.Error())
		}
		return fmt.Sprintf("[unfold] retrieved tenant(s): %v", strings.Join(tenants, ","))
	} else if *c.Opts.StatusFlag != "" {
		return fmt.Sprintf("[unfold] %s", GetAzureJobStatus(*c.Opts.StatusFlag))
	}
	return "[unfold] something went wrong"
}

// GetFlagSet returns the flag set for the get command
func (c commandGetConfig) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

// fetchCommandGetConfig fetches the command get config
func fetchCommandGetConfig() commandGetConfig {
	flagSet := flag.NewFlagSet(commands.Get, flag.ContinueOnError)
	return commandGetConfig{
		Opts: struct {
			TenantFlag *string
			StatusFlag *string
		}{
			TenantFlag: flagSet.String("t", "", "used to fetch tenant for given subscription"),
			StatusFlag: flagSet.String("s", "", "used to fetch status for given job id"),
		},
		FlagSet: flagSet,
	}
}
