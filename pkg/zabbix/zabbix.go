package zabbix

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/job"
	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/spool"
	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/state"
)

const (
	spoolFileName = "bareos-zabbix-check.spool"
)

var now = uint64(time.Now().Unix())

// Main the true main function of this program
func Main() string {
	err := processFlags()
	if err != nil {
		return fmt.Sprintf("INFO Failed to init program : %s", err)
	}
	// Open the state file
	stateFile, err := os.Open(stateFileName)
	if err != nil {
		return fmt.Sprintf("INFO Could not open state file : %s", err)
	}
	defer stateFile.Close()
	// parse the state file
	header, err := state.ParseHeader(stateFile)
	if err != nil {
		return fmt.Sprintf("INFO Could not parse state file header : %s", err)
	}
	// seek to the job entries in the state file
	offset, err := stateFile.Seek(int64(header.LastJobsAddr), 0)
	if err != nil {
		return fmt.Sprintf("INFO Couldn't seek to jobs position in state file : %s", err)
	}
	if uint64(offset) != header.LastJobsAddr {
		return fmt.Sprint("INFO Truncated state file")
	}
	// Then parse the jobs in the state file
	jobs, err := state.ParseJobs(stateFile)
	if err != nil {
		return fmt.Sprintf("INFO Could not parse jobs in state file : %s", err)
	}

	// We will check for errors in loading the spool file only at the end. If all jobs ran successfully without errors
	// in the state file and we manage to write a new spool file without errors, then we will ignore any error here to
	// avoid false positives during backup bootstrap
	// Open the spool file
	spoolFile, spoolErr := os.Open(filepath.Join(workDir, spoolFileName))
	var spoolJobs []job.Job
	if err == nil {
		defer spoolFile.Close()
		spoolJobs, spoolErr = spool.Parse(spoolFile)
	}

	jobs = job.KeepOldestOnly(append(jobs, spoolJobs...))

	// we write this new spool
	spoolFile, err = os.Create(filepath.Join(workDir, spoolFileName))
	if err == nil {
		defer spoolFile.Close()
		err = spool.Serialize(spoolFile, jobs)
	}
	if err != nil {
		return fmt.Sprintf("AVERAGE: Error saving the spool file : %s\n", err)
	}

	var (
		errorString   string
		missingString string
	)
	// We build the error strings
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		if job.Success {
			if job.Timestamp < now-24*3600 {
				if missingString == "" {
					missingString = fmt.Sprintf("missing: %s", job.Name)
				} else {
					missingString = fmt.Sprintf("%s, %s", missingString, job.Name)
				}
			}
		} else {
			if errorString == "" {
				errorString = fmt.Sprintf("errors: %s", job.Name)
			} else {
				errorString = fmt.Sprintf("%s, %s", errorString, job.Name)
			}
		}
	}
	// Finally we output
	if errorString != "" || missingString != "" {
		if spoolErr != nil {
			return fmt.Sprintf("AVERAGE: %s %s %s", errorString, missingString, spoolErr)
		}
		return fmt.Sprintf("AVERAGE: %s %s", errorString, missingString)
	}
	return "OK"
}
