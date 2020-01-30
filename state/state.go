package state

import (
	"bareos-zabbix-check/config"
	"bareos-zabbix-check/job"
)

// maxNameLength : the maximum length of a string, hard coded in bareos
const maxNameLength = 128

// State is an object for manipulating a bareos state file
type State struct {
	config *config.Config
	header header
	jobs   []job.Job
}

// Jobs returns the jobs from the state file
func (s *State) Jobs() []job.Job {
	return s.jobs
}
