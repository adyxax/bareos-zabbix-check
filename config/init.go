package config

import "flag"

// Init initialises a program config from the command line flags
func (c *Config) Init() {
	flag.BoolVar(&c.verbose, "v", false, "Activates verbose debugging output, defaults to false.")
	flag.BoolVar(&c.quiet, "q", false, "Suppress all output, suitable to force a silent update of the spool file.")
	flag.StringVar(&c.stateFile, "f", "", "Force the state file to use, defaults to "+bareosStateFile+" if it exists else "+baculaStateFile+".")
	flag.StringVar(&c.workDir, "w", "", "Force the work directory to use, defaults to "+bareosWorkDir+" if it exists else "+baculaWorkDir+".")

	// command line arguments parsing
	flag.Parse()
	c.checkWorkDir()
	c.checkStateFile()
}
