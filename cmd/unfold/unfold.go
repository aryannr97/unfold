package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aryannr97/unfold/pkg/azure"
	"github.com/aryannr97/unfold/pkg/commands"
	"github.com/aryannr97/unfold/pkg/google"
	"github.com/aryannr97/unfold/pkg/helpers"
	"github.com/aryannr97/unfold/pkg/registry"
	"github.com/aryannr97/unfold/pkg/spinner"
)

func main() {
	// Collect the cli command registry
	reg := registry.New()

	output := process(reg)
	log.Println(output)
}

func process(reg registry.Registry) string {
	defer helpers.GracefullyExit()
	// Check if the command is provided
	if len(os.Args) < 2 {
		return "[unfold] provide valid command"
	}

	// Get the command from the arguments
	inputCommand := os.Args[1]
	switch inputCommand {
	case commands.Azure:
		// Initialize the azure service
		err := azure.StartService()
		if err != nil {
			return fmt.Sprintf("[unfold] %s", err.Error())
		}
	case commands.Google:
		// Initialize the google service
		err := google.StartService()
		if err != nil {
			return fmt.Sprintf("[unfold] %s", err.Error())
		}

	}

	// Capture the output in common variable
	output := ""

	// Initialize the spinner
	spinner := spinner.Get(spinner.BrailDot)
	go spinner.Start()
	defer spinner.Clear()

	// Check if the sub-command or value is provided
	if len(os.Args) < 3 {
		output = "[unfold] provide valid sub-command or value for the command"
	} else {
		inputSubCommand := os.Args[2]

		if base, ok := reg[inputCommand]; !ok {
			output = fmt.Sprintf("[unfold] %s command not found", inputCommand)
		} else {
			if cmd, ok := base[inputSubCommand]; !ok {
				output = fmt.Sprintf("[unfold] %s %s command not found", inputCommand, inputSubCommand)
			} else {
				err := cmd.GetFlagSet().Parse(os.Args[3:])
				if err != nil {
					output = fmt.Sprintf("[unfold] %s", err.Error())
				} else {
					output = cmd.Execute()
				}
			}
		}
	}

	return output
}
