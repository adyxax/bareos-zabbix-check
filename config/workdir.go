package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const (
	bareosWorkDir = "/var/lib/bareos"
	baculaWorkDir = "/var/lib/bacula"
)

// checkWorkDir checks if a work directory is valid
func (c *Config) checkWorkDir() {
	// Determine the work directory to use.
	if c.workDir != "" {
		info, err := os.Stat(c.workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			fmt.Printf("INFO Invalid work directory %s : it does not exist or is not a directory.\n", c.workDir)
			os.Exit(0)
		}
	} else {
		c.workDir = bareosWorkDir
		info, err := os.Stat(c.workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			c.workDir = baculaWorkDir
			info, err := os.Stat(c.workDir)
			if os.IsNotExist(err) || !info.IsDir() {
				fmt.Println("INFO Could not find a suitable work directory. Is bareos or bacula installed?")
				os.Exit(0)
			}
		}
	}
	c.workDir = filepath.Clean(c.workDir)
	if c.verbose {
		log.Println("Setting work directory to ", c.workDir)
	}
}
