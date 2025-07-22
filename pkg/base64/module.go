package base64

// CommandModule represents the collection of different command configs
type CommandModule struct {
	CommandDecodeConfig commandDecodeConfig
}

// NewCommandModule returns the command module
func NewCommandModule() *CommandModule {
	return &CommandModule{
		CommandDecodeConfig: fetchCommandDecodeConfig(),
	}
}
