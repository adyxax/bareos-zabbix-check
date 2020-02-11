package job

import "fmt"

// Job is a bareos job
type Job struct {
	Name      string
	Timestamp uint64
	Success   bool
}

func (job Job) String() string {
	return fmt.Sprintf("Job { Name: \"%s\", Timestamp: \"%d\", Success: \"%t\" }", job.Name, job.Timestamp, job.Success)
}
