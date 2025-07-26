package main

import (
	"flag"
	"os"
	"testing"

	"github.com/aryannr97/unfold/pkg/registry"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "test main entrypoint",
			args: []string{"unfold", "google", "subcommand"},
		},
		{
			name: "test main entrypoint",
			args: []string{"unfold", "azure", "subcommand"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			main()
		})
	}
}

func Test_run(t *testing.T) {
	type args struct {
		reg registry.Registry
	}
	tests := []struct {
		name           string
		args           args
		cmdArgs        []string
		env            func()
		expectedOutput string
	}{
		{
			name: "test command success",
			args: args{
				reg: registry.Registry{
					"test": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold", "test", "subcommand"},
			expectedOutput: "test output",
		},
		{
			name: "test command failed missing command",
			args: args{
				reg: registry.Registry{
					"test1": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold"},
			expectedOutput: "[unfold] provide valid command",
		},
		{
			name: "test command failed missing subcommand",
			args: args{
				reg: registry.Registry{
					"test": {
						"subcommand1": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold", "test"},
			expectedOutput: "[unfold] provide valid sub-command or value for the command",
		},
		{
			name: "test command failed unsupported command",
			args: args{
				reg: registry.Registry{
					"test": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold", "test1", "subcommand"},
			expectedOutput: "[unfold] test1 command not found",
		},
		{
			name: "test command failed unsupported subcommand",
			args: args{
				reg: registry.Registry{
					"test": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold", "test", "subcommand1"},
			expectedOutput: "[unfold] test subcommand1 command not found",
		},
		{
			name: "test command failed unsupported flag",
			args: args{
				reg: registry.Registry{
					"test": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			cmdArgs:        []string{"unfold", "test", "subcommand", "-flag"},
			expectedOutput: "[unfold] flag provided but not defined: -flag",
		},
		{
			name: "azure command failed service not started",
			args: args{
				reg: registry.Registry{
					"azure": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			env: func() {
				os.Setenv("AZURE_OFFERS_FILE", "")
				os.Setenv("AZURE_CERT_FILE", "")
			},
			cmdArgs:        []string{"unfold", "azure", "subcommand"},
			expectedOutput: "[unfold] open : no such file or directory",
		},
		{
			name: "google command failed service not started",
			args: args{
				reg: registry.Registry{
					"google": {
						"subcommand": &MockCommand{
							Output:  "test output",
							FlagSet: flag.NewFlagSet("subcommand", flag.ContinueOnError),
						},
					},
				},
			},
			env: func() {
				os.Setenv("GOOGLE_KEYFILE", "")
			},
			cmdArgs:        []string{"unfold", "google", "subcommand"},
			expectedOutput: "[unfold] open : no such file or directory",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.cmdArgs
			if tt.env != nil {
				tt.env()
			}
			output := run(tt.args.reg)
			if output != tt.expectedOutput {
				t.Errorf("expected output %s, got %s", tt.expectedOutput, output)
			}
		})
	}
}

type MockCommand struct {
	Output  string
	FlagSet *flag.FlagSet
}

func (m *MockCommand) Execute() string {
	return m.Output
}

func (m *MockCommand) GetFlagSet() *flag.FlagSet {
	return m.FlagSet
}
