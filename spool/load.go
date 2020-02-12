package spool

import (
	"bareos-zabbix-check/config"
	"bareos-zabbix-check/job"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// Load loads a spool file in path
func (s *Spool) Load(c *config.Config) (err error) {
	s.config = c
	// We read the spool
	file, err := os.Open(filepath.Join(c.WorkDir(), spoolFile))
	if err != nil {
		return fmt.Errorf("Couldn't open spool file, starting from scratch: %s", err)
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return fmt.Errorf("Corrupted spool file, starting from scratch : %s", err)
	}
	if c.Verbose() {
		log.Printf("Spool file content : %v\n", lines)
	}

	for _, line := range lines {
		var i int
		i, err = strconv.Atoi(line[1])
		if err != nil {
			return fmt.Errorf("Corrupted spool file : couldn't parse timestamp entry")
		}
		s.jobs = append(s.jobs, job.Job{Name: line[0], Timestamp: uint64(i), Success: true})
	}
	return
}
