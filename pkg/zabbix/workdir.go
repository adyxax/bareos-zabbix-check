package zabbix

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	bareosWorkDir = "/var/lib/bareos"
	baculaWorkDir = "/var/lib/bacula"
)

var root = "/"

// checkWorkDir checks if a work directory is valid
func checkWorkDir() error {
	// Determine the work directory to use.
	if workDir != "" {
		workDir = filepath.Join(root, workDir)
		info, err := os.Stat(workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			return fmt.Errorf("Invalid work directory %s : it does not exist or is not a directory", workDir)
		}
	} else {
		workDir = filepath.Join(root, bareosWorkDir)
		info, err := os.Stat(workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			workDir = filepath.Join(root, baculaWorkDir)
			info, err := os.Stat(workDir)
			if os.IsNotExist(err) || !info.IsDir() {
				return fmt.Errorf("Could not find a suitable work directory. Is bareos or bacula installed?")
			}
		}
	}
	workDir = filepath.Clean(workDir)
	return nil
}
