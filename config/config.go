package config

// Config object
type Config struct {
	verbose   bool
	quiet     bool
	stateFile string
	workDir   string
}

// Verbose gets the verbose field of the configuration
func (config *Config) Verbose() bool { return config.verbose }

// Quiet gets the quiet field of the configuration
func (config *Config) Quiet() bool { return config.quiet }

// StateFile gets the stateFile field of the configuration
func (config *Config) StateFile() string { return config.stateFile }

// WorkDir gets the workDir field of the configuration
func (config *Config) WorkDir() string { return config.workDir }
