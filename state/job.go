package state

import (
	"bareos-zabbix-check/job"
	"bareos-zabbix-check/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"regexp"
	"time"

	"github.com/pkg/errors"
)

// maxnameLength : the maximum length of a job name, hardcoded in bareos
const maxNameLength = 128

// jobLength : the length of the job result struct
const jobLength = 16 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 4 + 8 + 8 + 8 + maxNameLength

var jobNameRegex = regexp.MustCompilePOSIX(`^([-A-Za-z0-9_]+)\.[0-9]{4}-[0-9]{2}-[0-9]{2}.*`)

// jobEntry : A structure to hold a job result from the state file
// This comes from bareos repository file core/src/lib/recent_job_results_list.h:29 and file core/src/lib/recent_job_results_list.cc:44
type jobEntry struct {
	_              [16]byte
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
	jobNameLen := utils.Clen(je.Job[:])
	if len(matches) >= 4 {
		jobNameLen = matches[3]
	}
	return fmt.Sprintf("Errors: %d, JobType: %c, JobStatus: %c, JobLevel: %c, JobID: %d, VolSessionID: %d, VolSessionTime: %d, JobFiles: %d, JobBytes: %d, StartTime: %s, EndTime: %s, Job: %s",
		je.Errors, je.JobType, je.JobStatus, je.JobLevel, je.JobID, je.VolSessionID, je.VolSessionTime, je.JobFiles, je.JobBytes, time.Unix(int64(je.StartTime), 0), time.Unix(int64(je.EndTime), 0), je.Job[:jobNameLen])
}

// ParseJobs parses the jobs in a state file
func ParseJobs(handle io.Reader) (jobs []job.Job, err error) {
	// We read how many jobs there are in the state file
	data := make([]byte, 4)
	n, err := handle.Read(data)
	if err != nil {
		return nil, errors.Wrap(err, "Corrupted state file")
	}
	if n != 4 {
		return nil, fmt.Errorf("Corrupted state file : invalid numberOfJobs read length")
	}
	buffer := bytes.NewBuffer(data)
	var numberOfJobs uint32
	_ = binary.Read(buffer, binary.LittleEndian, &numberOfJobs) // this call cannot fail since we checked the header length

	// We parse the job entries
	for ; numberOfJobs > 0; numberOfJobs-- {
		var (
			jobResult jobEntry
			jobName   string
		)
		data := make([]byte, jobLength)
		n, err = handle.Read(data)
		if err != nil {
			return nil, errors.Wrap(err, "Corrupted state file")
		}
		if n != jobLength {
			return nil, fmt.Errorf("Corrupted state file : invalid job entry")
		}
		buffer = bytes.NewBuffer(data)
		_ = binary.Read(buffer, binary.LittleEndian, &jobResult) // this call cannot fail since we checked the header length
		matches := jobNameRegex.FindSubmatchIndex(jobResult.Job[:])
		if len(matches) >= 4 {
			jobName = string(jobResult.Job[:matches[3]])
		} else {
			return nil, fmt.Errorf("Couldn't parse job name, this shouldn't happen : %s", jobResult.Job[:])
		}
		// If the job is of type backup (B == ascii 66)
		if jobResult.JobType == 66 {
			// If the job is of status success JobStatus is equals to 84 (T == ascii 84)
			jobs = append(jobs, job.Job{Name: jobName, Timestamp: jobResult.StartTime, Success: jobResult.JobStatus == 84})
		}
	}
	return
}
