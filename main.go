package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"time"
)

const (
	bareosWorkDir   = "/var/lib/bareos"
	bareosStateFile = "bareos-fd.9102.state"
	baculaWorkDir   = "/var/lib/bacula"
	baculaStateFile = "bacula-fd.9102.state"
	spoolFile       = "bareos-zabbix-check.spool"
)

// We declare globally the variables that will hold the command line arguments
var (
	verbose   bool
	quiet     bool
	stateFile string
	workDir   string
)

func main() {
	var (
		info           os.FileInfo
		err            error
		successfulJobs jobs
		errorJobs      jobs
		spoolJobs      jobs
		jobName        string
		ts             uint64
		now            uint64
		errorString    string
		missingString  string
	)

	// command line arguments parsing
	flag.BoolVar(&verbose, "v", false, "Activates verbose debugging output, defaults to false.")
	flag.BoolVar(&quiet, "q", false, "Suppress all output, suitable to force a silent update of the spool file.")
	flag.StringVar(&stateFile, "f", "", "Force the state file to use, defaults to "+bareosStateFile+" if it exists else "+baculaStateFile+".")
	flag.StringVar(&workDir, "w", "", "Force the work directory to use, defaults to "+bareosWorkDir+" if it exists else "+baculaWorkDir+".")
	flag.Parse()

	// Determine the work directory to use.
	if workDir != "" {
		info, err = os.Stat(workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			fmt.Printf("INFO Invalid work directory %s : it does not exist or is not a directory.\n", workDir)
			os.Exit(0)
		}
	} else {
		workDir = "/var/lib/bareos"
		info, err = os.Stat(workDir)
		if os.IsNotExist(err) || !info.IsDir() {
			workDir = "/var/lib/bacula"
			info, err := os.Stat(workDir)
			if os.IsNotExist(err) || !info.IsDir() {
				fmt.Println("INFO Could not find a suitable work directory. Is bareos or bacula installed?")
				os.Exit(0)
			}
		}
	}
	workDir = path.Clean(workDir)
	if verbose {
		log.Println("Setting work directory to ", workDir)
	}

	// Finds the state file to parse
	if stateFile != "" {
		stateFile = path.Join(workDir, stateFile)
		info, err = os.Stat(stateFile)
		if os.IsNotExist(err) || info.IsDir() {
			fmt.Printf("INFO The state file %s does not exist.\n", stateFile)
			os.Exit(0)
		}
	} else {
		stateFile = path.Join(workDir, bareosStateFile)
		info, err = os.Stat(stateFile)
		if os.IsNotExist(err) || info.IsDir() {
			stateFile = path.Join(workDir, baculaStateFile)
			info, err = os.Stat(stateFile)
			if os.IsNotExist(err) || info.IsDir() {
				fmt.Println("INFO Could not find a suitable state file. Has a job ever run?")
				os.Exit(0)
			}
		}
	}
	if verbose {
		log.Println("Using state file ", stateFile)
	}

	successfulJobs, errorJobs, err = parseStateFile()
	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}
	// We will check for errors in loading the spool file only if necessary. If all jobs ran successfully without errors in the state file and we manage to write
	// a new spool file without errors, then we will ignore any error here to avoid false positives during backup bootstrap
	spoolJobs, err = loadSpool()

	// if we have jobs in the spool we merge this list with successfull jobs from the state file
	if err == nil {
		for jobName, ts = range spoolJobs {
			var (
				current uint64
				ok      bool
			)
			current, ok = successfulJobs[jobName]
			if !ok || current < ts {
				successfulJobs[jobName] = ts
			}
		}
	}
	// we write this new spool
	if err2 := saveSpool(successfulJobs); err2 != nil {
		fmt.Printf("AVERAGE: Error saving the spool file : %s\n", err2)
		os.Exit(0)
	}

	// We build the error string listing the jobs in error
	for jobName, ts = range errorJobs {
		if errorString == "" {
			errorString = fmt.Sprintf("errors: %s", jobName)
		} else {
			errorString = fmt.Sprintf("%s, %s", errorString, jobName)
		}
	}
	now = uint64(time.Now().Unix())
	// Next we check if all jobs ran recently and build the missing string
	for jobName, ts = range successfulJobs {
		if ts < now-24*3600 {
			if missingString == "" {
				missingString = fmt.Sprintf("missing: %s", jobName)
			} else {
				missingString = fmt.Sprintf("%s, %s", missingString, jobName)
			}
		}
	}
	if errorString != "" || missingString != "" {
		fmt.Printf("AVERAGE: %s %s", errorString, missingString)
		if err != nil {
			fmt.Printf(" additionnal errors: %s", err)
		}
	} else {
		fmt.Printf("OK")
	}
}
