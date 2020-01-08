package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
)

// stateFileHeader : A structure to hold the header of the state file. It is statically aligned for amd64 architecture
type stateFileHeader struct {
	ID                        [14]byte
	_                         int16
	Version                   int32
	_                         int32
	LastJobsAddr              uint64
	EndOfRecentJobResultsList uint64
	Reserved                  [19]uint64
}

func (sfh stateFileHeader) String() string {
	return fmt.Sprintf("ID: \"%s\", Version: %d, LastJobsAddr: %d, EndOfRecentJobResultsList: %d", sfh.ID[:len(sfh.ID)-2], sfh.Version, sfh.EndOfRecentJobResultsList, sfh.Reserved)
}

// jobEntry : A structure to hold a job result from the state file
type jobEntry struct {
	Pad            [16]byte
	Errors         int32
	JobType        int32
	JobStatus      int32
	JobLevel       int32
	JobID          uint32
	VolSessionID   uint32
	VolSessionTime uint32
	JobFiles       uint32
	JobBytes       uint64
	StartTime      uint64
	EndTime        uint64
	Job            [maxNameLength]byte
}

func (je jobEntry) String() string {
	var matches = jobNameRegex.FindSubmatchIndex(je.Job[:])
	var jobNameLen int
	if len(matches) >= 4 {
		jobNameLen = matches[3]
	}
	return fmt.Sprintf("Errors: %d, JobType: %c, JobStatus: %c, JobLevel: %c, JobID: %d, VolSessionID: %d, VolSessionTime: %d, JobFiles: %d, JobBytes: %d, StartTime: %d, EndTime: %d, Job: %s",
		je.Errors, je.JobType, je.JobStatus, je.JobLevel, je.JobID, je.VolSessionID, je.VolSessionTime, je.JobFiles, je.JobBytes, je.StartTime, je.EndTime, je.Job[:jobNameLen])
}

const (
	// maxNameLength : the maximum length of a string, hard coded in bareos
	maxNameLength = 128
	// stateFileHeaderLength : the length of the state file header struct
	stateFileHeaderLength = 14 + 2 + 4 + 4 + 8 + 8 + 19*8
	// jobResultLength : the length of the job result struct
	jobResultLength = 16 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 8 + 8 + 8 + maxNameLength
)

var jobNameRegex = regexp.MustCompilePOSIX(`^([-A-Za-z0-9_]+)\.[0-9]{4}-[0-9]{2}-[0-9]{2}.*`)

// readNextBytes : Reads the next "number" bytes from a "file", returns the number of bytes actually read as well as the bytes read
func readNextBytes(file *os.File, number int) (int, []byte) {
	var bytes = make([]byte, number)

	var n, err = file.Read(bytes)
	if err != nil {
		fmt.Printf("INFO Corrupted state file : file.Read failed in %s : %s\n", stateFile, err)
	}

	return n, bytes
}

func parseStateFile() (successfulJobs jobs, errorJobs jobs, err error) {
	var (
		n               int
		stateFileHandle *os.File
		data            []byte
		buffer          *bytes.Buffer
		numberOfJobs    uint32
		matches         []int
	)
	// Open the state file
	stateFileHandle, err = os.Open(stateFile)
	if err != nil {
		return nil, nil, fmt.Errorf("INFO Couldn't open state file : %s", err)
	}
	defer stateFileHandle.Close()

	// Parsing the state file header
	var header stateFileHeader
	n, data = readNextBytes(stateFileHandle, stateFileHeaderLength)
	if n != stateFileHeaderLength {
		return nil, nil, fmt.Errorf("INFO Corrupted state file : invalid header length in %s", stateFile)
	}
	buffer = bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.LittleEndian, &header)
	if err != nil {
		return nil, nil, fmt.Errorf("INFO Corrupted state file : binary.Read failed on header in %s : %s", stateFile, err)
	}
	if verbose {
		log.Printf("Parsed header: %+s\n", header)
	}
	if id := string(header.ID[:len(header.ID)-1]); id != "Bareos State\n" && id != "Bacula State\n" {
		return nil, nil, fmt.Errorf("INFO Corrupted state file : Not a bareos or bacula state file %s", stateFile)
	}
	if header.Version != 4 {
		return nil, nil, fmt.Errorf("INFO Invalid state file : This script only supports bareos state file version 4, got %d", header.Version)
	}
	if header.LastJobsAddr == 0 {
		return nil, nil, fmt.Errorf("INFO No jobs exist in the state file")
	}

	// We seek to the jobs position in the state file
	stateFileHandle.Seek(int64(header.LastJobsAddr), 0)

	// We read how many jobs there are in the state file
	n, data = readNextBytes(stateFileHandle, 4)
	if n != 4 {
		return nil, nil, fmt.Errorf("INFO Corrupted state file : invalid numberOfJobs read length in %s", stateFile)
	}
	buffer = bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.LittleEndian, &numberOfJobs)
	if err != nil {
		return nil, nil, fmt.Errorf("INFO Corrupted state file : binary.Read failed on numberOfJobs in %s : %s", stateFile, err)
	}
	if verbose {
		log.Printf("%d jobs found in state file\n", numberOfJobs)
	}

	// We parse the job entries
	successfulJobs = make(map[string]uint64)
	errorJobs = make(map[string]uint64)
	for ; numberOfJobs > 0; numberOfJobs-- {
		var (
			jobResult jobEntry
			jobName   string
		)
		n, data = readNextBytes(stateFileHandle, jobResultLength)
		if n != jobResultLength {
			return nil, nil, fmt.Errorf("INFO Corrupted state file : invalid job entry in %s", stateFile)
		}
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &jobResult)
		if err != nil {
			return nil, nil, fmt.Errorf("INFO Corrupted state file : binary.Read failed on job entry in %s : %s", stateFile, err)
		}
		matches = jobNameRegex.FindSubmatchIndex(jobResult.Job[:])
		if len(matches) >= 4 {
			jobName = string(jobResult.Job[:matches[3]])
		} else {
			return nil, nil, fmt.Errorf("INFO Couldn't parse job name, this shouldn't happen : %s", jobResult.Job[:])
		}
		if verbose {
			log.Printf("Parsed job entry: %s\n", jobResult)
		}
		// If the job is of type backup (B == ascii 66)
		if jobResult.JobType == 66 {
			var (
				successExists  bool
				errorExists    bool
				currentSuccess uint64
				currentError   uint64
			)
			currentSuccess, successExists = successfulJobs[jobName]
			currentError, errorExists = errorJobs[jobName]
			// If the job is of status success (T == ascii 84)
			if jobResult.JobStatus == 84 {
				// if there is an older entry in errorJobs we delete it
				if errorExists && jobResult.StartTime > currentError {
					delete(errorJobs, jobName)
				}
				// if there are no entries more recent in successfulJobs we add this one
				if !successExists || successExists && jobResult.StartTime > currentSuccess {
					successfulJobs[jobName] = jobResult.StartTime
				}
			} else {
				if !errorExists || jobResult.StartTime > currentError {
					errorJobs[jobName] = jobResult.StartTime
				}
			}
		}
	}
	return
}
