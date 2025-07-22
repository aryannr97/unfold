package base64

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"

	"github.com/aryannr97/unfold/pkg/commands"
)

type commandDecodeConfig struct {
	FlagSet *flag.FlagSet
}

// Execute executes the decode command
func (c commandDecodeConfig) Execute() string {
	tokenString := os.Args[3]

	decodeBytes, err := base64.StdEncoding.DecodeString(tokenString)
	if err != nil {
		return fmt.Sprintf("[unfold] failed to decode given string %v", err)
	}

	return fmt.Sprintf("[unfold] decoded content for base64 string\n%s", string(decodeBytes))
}

// GetFlagSet returns the flag set for the decode command
func (c commandDecodeConfig) GetFlagSet() *flag.FlagSet {
	return c.FlagSet
}

func fetchCommandDecodeConfig() commandDecodeConfig {
	flagSet := flag.NewFlagSet(commands.Decode, flag.ContinueOnError)
	return commandDecodeConfig{
		FlagSet: flagSet,
	}
}
