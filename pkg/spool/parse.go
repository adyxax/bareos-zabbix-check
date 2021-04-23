package spool

import (
	"encoding/csv"
	"io"
	"strconv"

	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/job"
	"github.com/pkg/errors"
)

// Parse parses a spool file
func Parse(handle io.Reader) (jobs []job.Job, err error) {
	lines, err := csv.NewReader(handle).ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Corrupted spool file")
	}
	for n := 0; n < len(lines); n++ {
		line := lines[n]
		i, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, errors.Wrapf(err, "Corrupted spool file : couldn't parse timestamp entry : %s", line[1])
		}
		jobs = append(jobs, job.Job{Name: line[0], Timestamp: uint64(i), Success: true})
	}
	return
}
