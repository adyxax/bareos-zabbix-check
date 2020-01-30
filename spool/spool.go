package spool

import (
	"bareos-zabbix-check/config"
	"bareos-zabbix-check/job"
)

const (
	spoolFile = "bareos-zabbix-check.spool"
)

// Spool is an object for manipulating a bareos spool file
type Spool struct {
	config *config.Config
	jobs   []job.Job
}

// Jobs exports a spool to a jobs list
func (s *Spool) Jobs() []job.Job {
	return s.jobs
}

// SetJobs sets a jobs list
func (s *Spool) SetJobs(jobs []job.Job) {
	s.jobs = jobs
}
