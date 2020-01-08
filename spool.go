package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
)

// jobs is a map that matches a job name string to its last successfull run timestamp
type jobs map[string]uint64

func loadSpool() (entries jobs, err error) {
	var (
		file  *os.File
		lines [][]string
	)
	// We read the spool
	file, err = os.Open(path.Join(workDir, spoolFile))
	if err != nil {
		return nil, fmt.Errorf("INFO Couldn't open spool file: %s", err)
	}
	defer file.Close()
	lines, err = csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("INFO Corrupted spool file : %s", err)
	}
	if verbose {
		log.Printf("Spool file content : %v\n", lines)
	}

	entries = make(map[string]uint64)
	for _, line := range lines {
		var (
			i int
		)
		i, err = strconv.Atoi(line[1])
		if err != nil {
			return nil, fmt.Errorf("INFO Corrupted spool file : couldn't parse timestamp entry")
		}
		entries[line[0]] = uint64(i)
	}
	return
}

func saveSpool(entries jobs) (err error) {
	var (
		file    *os.File
		lines   [][]string
		jobName string
		ts      uint64
		i       int
	)
	file, err = os.Create(path.Join(workDir, spoolFile))
	if err != nil {
		return fmt.Errorf("INFO Couldn't open spool file for writing : %s", err)
	}
	defer file.Close()

	lines = make([][]string, len(entries))
	i = 0
	for jobName, ts = range entries {
		lines[i] = make([]string, 2)
		lines[i][0] = jobName
		lines[i][1] = fmt.Sprintf("%d", ts)
		i++
	}
	err = csv.NewWriter(file).WriteAll(lines)
	return
}
