package google

import (
	"flag"
	"fmt"

	"github.com/aryannr97/unfold/pkg/commands"
)

// commandGetConfig represents the configuration for the get command
type commandGetConfig struct {
	FlagSet *flag.FlagSet
	Opts    struct {
		GroupFlag *string
	}
}

// Execute executes the get command
func (c commandGetConfig) Execute() string {
	if *c.Opts.GroupFlag != "" {
		res, err := GetGroupByID(*c.Opts.GroupFlag)
		if err != nil {
			return fmt.Sprintf("[unfold] %s", err.Error())
		}
		return fmt.Sprintf("[unfold] retrieved group resource name: %v", res.Name)
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
			GroupFlag *string
		}{
			GroupFlag: flagSet.String("g", "", "used to fetch group details for given group id"),
		},
		FlagSet: flagSet,
	}
}
