package state

import (
	"bareos-zabbix-check/utils"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// c.StateFile()HeaderLength : the length of the state file header struct
const headerLength = 14 + 2 + 4 + 4 + 8 + 8 // + 19*8

// Header is a structure to hold the header of the state file. It is statically aligned for amd64 architecture
// This comes from bareos repository file core/src/lib/bsys.cc:525 and core/src/lib/bsys.cc:652
type Header struct {
	ID                        [14]byte
	_                         int16
	Version                   int32
	_                         int32
	LastJobsAddr              uint64
	EndOfRecentJobResultsList uint64
	//Reserved                  [19]uint64
}

func (sfh *Header) String() string {
	return fmt.Sprintf("ID: \"%s\", Version: %d, LastJobsAddr: %d, EndOfRecentJobResultsList: %d",
		string(sfh.ID[:utils.Clen(sfh.ID[:])]), sfh.Version, sfh.LastJobsAddr, sfh.EndOfRecentJobResultsList)
}

// ParseHeader parses a Header struct
func ParseHeader(handle io.Reader) (h *Header, err error) {
	// Parsing the state file header
	data := make([]byte, headerLength)
	n, err := handle.Read(data)
	if err != nil {
		return nil, errors.Wrap(err, "Corrupted state file")
	}
	if n != headerLength {
		return nil, fmt.Errorf("Corrupted state file : invalid header length")
	}
	buffer := bytes.NewBuffer(data)
	h = &Header{}
	_ = binary.Read(buffer, binary.LittleEndian, h) // this call cannot fail since we checked the header length
	if id := string(h.ID[:utils.Clen(h.ID[:])]); id != "Bareos State\n" && id != "Bacula State\n" {
		return nil, fmt.Errorf("Corrupted state file : Not a bareos or bacula state file : %s", id)
	}
	if h.Version != 4 {
		return nil, fmt.Errorf("Invalid state file : This script only supports bareos state file version 4, got %d", h.Version)
	}
	if h.LastJobsAddr == 0 {
		return nil, fmt.Errorf("No jobs exist in the state file")
	}
	return
}
