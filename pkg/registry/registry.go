package registry

import (
	"flag"

	"github.com/aryannr97/unfold/pkg/azure"
	"github.com/aryannr97/unfold/pkg/commands"
	"github.com/aryannr97/unfold/pkg/google"
	"github.com/aryannr97/unfold/pkg/jwt"
)

// Operation represents the operation to be executed
type Operation interface {
	Execute() string
	GetFlagSet() *flag.FlagSet
}

// Registry represents the commands to be executed
type Registry map[string]map[string]Operation

// New returns the list of commands to be executed
func New() Registry {
	return Registry{
		commands.Azure: {
			commands.Get:       azure.NewCommandModule().CommandGetConfig,
			commands.Search:    azure.NewCommandModule().CommandSearchConfig,
			commands.Configure: azure.NewCommandModule().CommandConfigureConfig,
		},
		commands.Google: {
			commands.Get:       google.NewCommandModule().CommandGetConfig,
			commands.Search:    google.NewCommandModule().CommandSearchConfig,
			commands.Configure: google.NewCommandModule().CommandConfigureConfig,
		},
		commands.JWT: {
			commands.Decode: jwt.NewCommandModule().CommandDecodeConfig,
		},
	}
}
