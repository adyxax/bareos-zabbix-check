package spool

import (
	"bareos-zabbix-check/job"
	"bytes"
	"io"
	"reflect"
	"testing"
	"testing/iotest"
)

func TestParse(t *testing.T) {
	readerError := iotest.TimeoutReader(bytes.NewReader([]byte("\n")))
	readerCorruptedTimestamp := bytes.NewReader([]byte("test,x"))
	readerOneJob := bytes.NewReader([]byte("test,1"))
	type args struct {
		handle io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantJobs []job.Job
		wantErr  bool
	}{
		{"empty", args{readerError}, nil, true},
		{"corrupted timestamp", args{readerCorruptedTimestamp}, nil, true},
		{"one job", args{readerOneJob}, []job.Job{{Name: "test", Timestamp: 1, Success: true}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJobs, err := Parse(tt.args.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJobs, tt.wantJobs) {
				t.Errorf("Parse() = %v, want %v", gotJobs, tt.wantJobs)
			}
		})
	}
}
