package state

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

// c.StateFile()HeaderLength : the length of the state file header struct
const headerLength = 14 + 2 + 4 + 4 + 8 + 8 + 19*8

// header : A structure to hold the header of the state file. It is statically aligned for amd64 architecture
// This comes from bareos repository file core/src/lib/bsys.cc:525 and core/src/lib/bsys.cc:652
type header struct {
	ID                        [14]byte
	_                         int16
	Version                   int32
	_                         int32
	LastJobsAddr              uint64
	EndOfRecentJobResultsList uint64
	Reserved                  [19]uint64
}

func (sfh header) String() string {
	return fmt.Sprintf("ID: \"%s\", Version: %d, LastJobsAddr: %d, EndOfRecentJobResultsList: %d", sfh.ID[:len(sfh.ID)-2], sfh.Version, sfh.EndOfRecentJobResultsList, sfh.Reserved)
}

func (s *State) parseHeader(file *os.File) (err error) {
	// Parsing the state file header
	n, data, err := s.readNextBytes(file, headerLength)
	if err != nil {
		return fmt.Errorf("INFO Corrupted state file : %s", err)
	}
	if n != headerLength {
		return fmt.Errorf("INFO Corrupted state file : invalid header length in %s", s.config.StateFile())
	}
	buffer := bytes.NewBuffer(data)
	err = binary.Read(buffer, binary.LittleEndian, &s.header)
	if err != nil {
		return fmt.Errorf("INFO Corrupted state file : binary.Read failed on header in %s : %s", s.config.StateFile(), err)
	}
	if s.config.Verbose() {
		log.Printf("Parsed header: %+s\n", s.header)
	}
	if id := string(s.header.ID[:len(s.header.ID)-1]); id != "Bareos State\n" && id != "Bacula State\n" {
		return fmt.Errorf("INFO Corrupted state file : Not a bareos or bacula state file %s", s.config.StateFile())
	}
	if s.header.Version != 4 {
		return fmt.Errorf("INFO Invalid state file : This script only supports bareos state file version 4, got %d", s.header.Version)
	}
	if s.header.LastJobsAddr == 0 {
		return fmt.Errorf("INFO No jobs exist in the state file")
	}
	return
}
