package state

import (
	"bareos-zabbix-check/config"
	"fmt"
	"os"
)

// Parse parses a bareos state file
func (s *State) Parse(c *config.Config) (err error) {
	s.config = c
	// Open the state file
	file, err := os.Open(c.StateFile())
	if err != nil {
		return fmt.Errorf("INFO Couldn't open state file : %s", err)
	}
	defer file.Close()

	err = s.parseHeader(file)
	if err != nil {
		return err
	}
	err = s.parseJobs(file)
	if err != nil {
		return err
	}

	return
}

// readNextBytes : Reads the next "number" bytes from a "file", returns the number of bytes actually read as well as the bytes read
func (s *State) readNextBytes(file *os.File, number int) (n int, bytes []byte, err error) {
	bytes = make([]byte, number)
	n, err = file.Read(bytes)
	if err != nil {
		return 0, nil, fmt.Errorf("file.Read failed in %s : %s", s.config.StateFile(), err)
	}

	return
}
