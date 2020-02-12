package spool

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

// Save writes a spool on the disk
func (s *Spool) Save() (err error) {
	file, err := os.Create(filepath.Join(s.config.WorkDir(), spoolFile))
	if err != nil {
		return
	}
	defer file.Close()

	lines := make([][]string, len(s.jobs))
	var i int = 0
	for _, job := range s.jobs {
		lines[i] = make([]string, 2)
		lines[i][0] = job.Name
		lines[i][1] = fmt.Sprintf("%d", job.Timestamp)
		i++
	}
	err = csv.NewWriter(file).WriteAll(lines)
	return
}
