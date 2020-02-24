package zabbix

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	bareosStateFile = "bareos-fd.9102.state"
	baculaStateFile = "bacula-fd.9102.state"
)

func checkStateFile() error {
	// Finds the state file to parse
	if stateFileName != "" {
		stateFileName = filepath.Join(workDir, stateFileName)
		info, err := os.Stat(stateFileName)
		if os.IsNotExist(err) || info.IsDir() {
			return fmt.Errorf("The state file %s does not exist", stateFileName)
		}
	} else {
		stateFileName = filepath.Join(workDir, bareosStateFile)
		info, err := os.Stat(stateFileName)
		if os.IsNotExist(err) || info.IsDir() {
			stateFileName = filepath.Join(workDir, baculaStateFile)
			info, err = os.Stat(stateFileName)
			if os.IsNotExist(err) || info.IsDir() {
				return fmt.Errorf("Could not find a suitable state file. Has a job ever run?")
			}
		}
	}
	return nil
}
