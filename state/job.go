package state

import (
	"bareos-zabbix-check/job"
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"
)

// jobLength : the length of the job result struct
const jobLength = 16 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 8 + 8 + 8 + maxNameLength

var jobNameRegex = regexp.MustCompilePOSIX(`^([-A-Za-z0-9_]+)\.[0-9]{4}-[0-9]{2}-[0-9]{2}.*`)

// jobEntry : A structure to hold a job result from the state file
// This comes from bareos repository file core/src/lib/recent_job_results_list.h:29 and file core/src/lib/recent_job_results_list.cc:44
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
	return fmt.Sprintf("Errors: %d, JobType: %c, JobStatus: %c, JobLevel: %c, JobID: %d, VolSessionID: %d, VolSessionTime: %d, JobFiles: %d, JobBytes: %d, StartTime: %s, EndTime: %s, Job: %s",
		je.Errors, je.JobType, je.JobStatus, je.JobLevel, je.JobID, je.VolSessionID, je.VolSessionTime, je.JobFiles, je.JobBytes, time.Unix(int64(je.StartTime), 0), time.Unix(int64(je.EndTime), 0), je.Job[:jobNameLen])
}

func (s *State) parseJobs(file *os.File) (err error) {
	// We seek to the jobs position in the state file
	file.Seek(int64(s.header.LastJobsAddr), 0)

	// We read how many jobs there are in the state file
	n, data, err := s.readNextBytes(file, 4)
	if err != nil {
		return fmt.Errorf("INFO Corrupted state file : %s", err)
	}
	if n != 4 {
		return fmt.Errorf("INFO Corrupted state file : invalid numberOfJobs read length in %s", s.config.StateFile())
	}
	buffer := bytes.NewBuffer(data)
	var numberOfJobs uint32
	err = binary.Read(buffer, binary.LittleEndian, &numberOfJobs)
	if err != nil {
		return fmt.Errorf("INFO Corrupted state file : binary.Read failed on numberOfJobs in %s : %s", s.config.StateFile(), err)
	}
	if s.config.Verbose() {
		log.Printf("%d jobs found in state file\n", numberOfJobs)
	}

	// We parse the job entries
	for ; numberOfJobs > 0; numberOfJobs-- {
		var (
			jobResult jobEntry
			jobName   string
		)
		n, data, err = s.readNextBytes(file, jobLength)
		if err != nil {
			return fmt.Errorf("INFO Corrupted state file : %s", err)
		}
		if n != jobLength {
			return fmt.Errorf("INFO Corrupted state file : invalid job entry in %s", s.config.StateFile())
		}
		buffer = bytes.NewBuffer(data)
		err = binary.Read(buffer, binary.LittleEndian, &jobResult)
		if err != nil {
			return fmt.Errorf("INFO Corrupted state file : binary.Read failed on job entry in %s : %s", s.config.StateFile(), err)
		}
		matches := jobNameRegex.FindSubmatchIndex(jobResult.Job[:])
		if len(matches) >= 4 {
			jobName = string(jobResult.Job[:matches[3]])
		} else {
			return fmt.Errorf("INFO Couldn't parse job name, this shouldn't happen : %s", jobResult.Job[:])
		}
		if s.config.Verbose() {
			log.Printf("Parsed job entry: %s\n", jobResult)
		}
		// If the job is of type backup (B == ascii 66)
		if jobResult.JobType == 66 {
			// If the job is of status success JobStatus is equals to 84 (T == ascii 84)
			s.jobs = append(s.jobs, job.Job{Name: jobName, Timestamp: jobResult.StartTime, Success: jobResult.JobStatus == 84})
		}
	}
	return
}
