package google

import (
	"flag"
	"fmt"

	"github.com/aryannr97/unfold/pkg/commands"
	"github.com/aryannr97/unfold/pkg/helpers"
)

// commandConfigureConfig represents the configuration for the configure command
type commandConfigureConfig struct {
	FlagSet       *flag.FlagSet
	AddRemoveOpts struct {
		RemoveFlag *bool
		EmailID    *string
		Group      *string
	}
}

// Execute executes the configure command
func (c commandConfigureConfig) Execute() string {
	if *c.AddRemoveOpts.RemoveFlag {
		err := RemoveMemberFromGroupID(*c.AddRemoveOpts.Group, *c.AddRemoveOpts.EmailID)
		if err != nil {
			return fmt.Sprintf("[unfold] unable to remove the member %s", helpers.RedValue(err.Error()))
		}
		return "[unfold] successfully removed the member from the group"
	}
	err := AddMemberToGroupID(*c.AddRemoveOpts.Group, *c.AddRemoveOpts.EmailID)
	if err != nil {
		return fmt.Sprintf("[unfold] failed to add member to the given group %s", helpers.RedValue(err.Error()))
	}

	return "[unfold] successfully added the member to the given group"
}

// GetFlagSet returns the flag set for the configure command
func (c commandConfigureConfig) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

// fetchCommandConfigureConfig fetches the command configure config
func fetchCommandConfigureConfig() commandConfigureConfig {
	flagSet := flag.NewFlagSet(commands.Configure, flag.ContinueOnError)
	return commandConfigureConfig{
		AddRemoveOpts: struct {
			RemoveFlag *bool
			EmailID    *string
			Group      *string
		}{
			RemoveFlag: flagSet.Bool("r", false, "remove emailID from respective google group"),
			EmailID:    flagSet.String("id", "", "provide a valid emaildID"),
			Group:      flagSet.String("g", "", "provide a valid google group"),
		},
		FlagSet: flagSet,
	}
}
