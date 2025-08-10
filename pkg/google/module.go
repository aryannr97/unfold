package google

// CommandModule represents the collection of different command configs
type CommandModule struct {
	CommandGetConfig       commandGetConfig
	CommandConfigureConfig commandConfigureConfig
	CommandSearchConfig    commandSearchConfig
}

// NewCommandModule returns the command module
func NewCommandModule() *CommandModule {
	return &CommandModule{
		CommandGetConfig:       fetchCommandGetConfig(),
		CommandConfigureConfig: fetchCommandConfigureConfig(),
		CommandSearchConfig:    fetchCommandSearchConfig(),
	}
}
