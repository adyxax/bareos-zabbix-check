package config

import (
	"fmt"
	"log"
	"os"
	"path"
)

const (
	bareosStateFile = "bareos-fd.9102.state"
	baculaStateFile = "bacula-fd.9102.state"
)

func (c *Config) checkStateFile() {
	// Finds the state file to parse
	if c.stateFile != "" {
		c.stateFile = path.Join(c.workDir, c.stateFile)
		info, err := os.Stat(c.stateFile)
		if os.IsNotExist(err) || info.IsDir() {
			fmt.Printf("INFO The state file %s does not exist.\n", c.stateFile)
			os.Exit(0)
		}
	} else {
		c.stateFile = path.Join(c.workDir, bareosStateFile)
		info, err := os.Stat(c.stateFile)
		if os.IsNotExist(err) || info.IsDir() {
			c.stateFile = path.Join(c.workDir, baculaStateFile)
			info, err = os.Stat(c.stateFile)
			if os.IsNotExist(err) || info.IsDir() {
				fmt.Println("INFO Could not find a suitable state file. Has a job ever run?")
				os.Exit(0)
			}
		}
	}
	if c.verbose {
		log.Println("Using state file ", c.stateFile)
	}
}
