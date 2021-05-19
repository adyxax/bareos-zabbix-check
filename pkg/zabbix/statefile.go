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
		_, err := os.Stat(stateFileName)
		if err != nil {
			return fmt.Errorf("Could not open state file %s", stateFileName)
		}
	} else {
		stateFileName = filepath.Join(workDir, bareosStateFile)
		_, err := os.Stat(stateFileName)
		if err != nil {
			stateFileName = filepath.Join(workDir, baculaStateFile)
			_, err = os.Stat(stateFileName)
			if err != nil {
				return fmt.Errorf("Could not autodetect a suitable state file. Has a job ever run? Does the user you are running the check as has read access to bacula or bareos' /var/lib directory? Alternatively use the -w and -f flags to specify the work directory and state file to use.")
			}
		}
	}
	return nil
}
