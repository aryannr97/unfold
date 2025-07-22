package jwt

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/aryannr97/unfold/pkg/commands"
	"github.com/golang-jwt/jwt/v5"
)

type commandDecodeConfig struct {
	FlagSet *flag.FlagSet
}

// Execute executes the decode command
func (c commandDecodeConfig) Execute() string {
	tokenString := os.Args[3]

	// Parse the token
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, jwt.MapClaims{})

	// Handle errors
	if err != nil {
		return fmt.Sprintf("[unfold] failed to parse token, %v", err)
	}

	// Check if the token is valid
	claims, _ := token.Claims.(jwt.MapClaims)
	// Convert claims to JSON
	jsonOutput, _ := json.MarshalIndent(claims, "", "  ")

	// Print the JSON output
	return fmt.Sprintf("[unfold] unfolded JWT Claims (JSON)\n%v", string(jsonOutput))
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
