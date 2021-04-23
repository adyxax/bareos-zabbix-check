package spool

import (
	"encoding/csv"
	"fmt"
	"io"

	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/job"
)

// Serialize writes a spool on the disk
func Serialize(handle io.Writer, jobs []job.Job) error {
	lines := make([][]string, len(jobs))
	for i := 0; i < len(jobs); i++ {
		job := jobs[i]
		lines[i] = make([]string, 2)
		lines[i][0] = job.Name
		lines[i][1] = fmt.Sprintf("%d", job.Timestamp)
	}
	return csv.NewWriter(handle).WriteAll(lines)
}
