package main

import (
	"bareos-zabbix-check/config"
	"bareos-zabbix-check/job"
	"bareos-zabbix-check/spool"
	"bareos-zabbix-check/state"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	var (
		config        config.Config
		spool         spool.Spool
		errorString   string
		missingString string
	)
	config.Init()
	// Open the state file
	stateFile, err := os.Open(config.StateFile())
	if err != nil {
		fmt.Printf("INFO Couldn't open state file : %s", err)
		os.Exit(0)
	}
	defer stateFile.Close()
	// parse the state file
	header, err := state.ParseHeader(stateFile)
	if err != nil {
		fmt.Printf("INFO Could not parse state file header : %s", err)
		os.Exit(0)
	}
	if config.Verbose() {
		log.Printf("Parsed header: %+s\n", header)
	}
	// seek to the job entries in the state file
	offset, err := stateFile.Seek(int64(header.LastJobsAddr), 0)
	if err != nil {
		fmt.Printf("INFO Couldn't seek to jobs position in state file : %s", err)
	}
	if uint64(offset) != header.LastJobsAddr {
		fmt.Print("INFO Truncated state file")
	}
	// Then parse the jobs in the state file
	jobs, err := state.ParseJobs(stateFile)
	if err != nil {
		fmt.Printf("INFO Could not parse jobs in state file : %s", err)
	}
	if config.Verbose() {
		log.Printf("%d jobs found in state file\n", len(jobs))
		for i := 0; i < len(jobs); i++ {
			log.Print(jobs[i])
		}
	}

	// We will check for errors in loading the spool file only at the end. If all jobs ran successfully without errors
	// in the state file and we manage to write a new spool file without errors, then we will ignore any error here to
	// avoid false positives during backup bootstrap
	err = spool.Load(&config)

	jobs = job.KeepOldestOnly(append(jobs, spool.Jobs()...))
	spool.SetJobs(job.KeepSuccessOnly(jobs))

	// we write this new spool
	if err2 := spool.Save(); err2 != nil {
		fmt.Printf("AVERAGE: Error saving the spool file : %s\n", err2)
		os.Exit(0)
	}

	now := uint64(time.Now().Unix())
	// We build the error strings
	for _, job := range jobs {
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
		fmt.Printf("AVERAGE: %s %s", errorString, missingString)
		if err != nil {
			fmt.Printf(" additionnal errors: %s", err)
		}
	} else {
		fmt.Printf("OK")
	}
}
