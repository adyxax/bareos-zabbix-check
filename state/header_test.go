package state

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_header_String(t *testing.T) {
	var id [14]byte
	copy(id[:], []byte("test"))
	type fields struct {
		ID                        [14]byte
		Version                   int32
		LastJobsAddr              uint64
		EndOfRecentJobResultsList uint64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"default header", fields{ID: id, Version: 1, LastJobsAddr: 2, EndOfRecentJobResultsList: 3}, "ID: \"test\", Version: 1, LastJobsAddr: 2, EndOfRecentJobResultsList: 3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sfh := Header{
				ID:                        tt.fields.ID,
				Version:                   tt.fields.Version,
				LastJobsAddr:              tt.fields.LastJobsAddr,
				EndOfRecentJobResultsList: tt.fields.EndOfRecentJobResultsList,
			}
			if got := sfh.String(); got != tt.want {
				t.Errorf("header.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseHeader(t *testing.T) {
	readerEmpty := bytes.NewReader([]byte(""))
	readerTooSmall := bytes.NewReader([]byte("abcd"))
	readerNotBareosNorBacula := bytes.NewReader([]byte{
		't', 'e', 's', 't', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // ID
		0, 0, // padding
		4, 0, 0, 0, // version
		0, 0, 0, 0, //padding
		0, 0, 0, 0, 0, 0, 0, 0, // last job address
		0, 0, 0, 0, 0, 0, 0, 0, // EndOfRecentJobResultsLis
	})
	readerBadVersion := bytes.NewReader([]byte{
		'B', 'a', 'r', 'e', 'o', 's', ' ', 'S', 't', 'a', 't', 'e', '\n', 0, // ID
		0, 0, // padding
		3, 0, 0, 0, // version
		0, 0, 0, 0, //padding
		0, 0, 0, 0, 0, 0, 0, 0, // last job address
		0, 0, 0, 0, 0, 0, 0, 0, // EndOfRecentJobResultsLis
	})
	readerNoJobs := bytes.NewReader([]byte{
		'B', 'a', 'r', 'e', 'o', 's', ' ', 'S', 't', 'a', 't', 'e', '\n', 0, // ID
		0, 0, // padding
		4, 0, 0, 0, // version
		0, 0, 0, 0, //padding
		0, 0, 0, 0, 0, 0, 0, 0, // last job address
		0, 0, 0, 0, 0, 0, 0, 0, // EndOfRecentJobResultsList
	})
	readerValid := bytes.NewReader([]byte{
		'B', 'a', 'r', 'e', 'o', 's', ' ', 'S', 't', 'a', 't', 'e', '\n', 0, // ID
		0, 0, // padding
		4, 0, 0, 0, // version
		0, 0, 0, 0, //padding
		192, 0, 0, 0, 0, 0, 0, 0, // last job address
		254, 0, 0, 0, 0, 0, 0, 0, // EndOfRecentJobResultsList
	})
	type args struct {
		handle io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantH   *Header
		wantErr bool
	}{
		{"read error", args{readerEmpty}, nil, true},
		{"invalid header length", args{readerTooSmall}, nil, true},
		{"reader not bareos nor bacula", args{readerNotBareosNorBacula}, nil, true},
		{"reader bad version", args{readerBadVersion}, nil, true},
		{"reader no jobs", args{readerNoJobs}, nil, true},
		{"reader valid", args{readerValid}, &Header{
			ID:                        [14]byte{'B', 'a', 'r', 'e', 'o', 's', ' ', 'S', 't', 'a', 't', 'e', '\n', 0},
			Version:                   4,
			LastJobsAddr:              192,
			EndOfRecentJobResultsList: 254,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotH, err := ParseHeader(tt.args.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHeader() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotH, tt.wantH) {
				t.Errorf("parseHeader() = %v, want %v", gotH, tt.wantH)
			}
		})
	}
}
