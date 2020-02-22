package state

import (
	"bareos-zabbix-check/job"
	"bytes"
	"io"
	"reflect"
	"testing"
)

func Test_jobEntry_String(t *testing.T) {
	var badlyNamedJob [128]byte
	copy(badlyNamedJob[:], []byte("job_name"))
	var normalJob [128]byte
	copy(normalJob[:], []byte("normal_name.2012-06-01"))
	type fields struct {
		Errors         int32
		JobType        int32
		JobStatus      int32
		JobLevel       int32
		JobID          uint32
		VolSessionID   uint32
		VolSessionTime uint32
		JobFiles       uint32
		JobBytes       uint64
		StartTime      uint64
		EndTime        uint64
		Job            [maxNameLength]byte
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"normal job",
			fields{Errors: 1, JobType: 'B', JobStatus: 'T', JobLevel: 'F', JobID: 2, VolSessionID: 3, VolSessionTime: 4, JobFiles: 5, JobBytes: 6, Job: badlyNamedJob},
			"Errors: 1, JobType: B, JobStatus: T, JobLevel: F, JobID: 2, VolSessionID: 3, VolSessionTime: 4, JobFiles: 5, JobBytes: 6, StartTime: 1970-01-01 01:00:00 +0100 CET, EndTime: 1970-01-01 01:00:00 +0100 CET, Job: job_name",
		},
		{
			"badly named job",
			fields{Errors: 1, JobType: 'B', JobStatus: 'T', JobLevel: 'F', JobID: 2, VolSessionID: 3, VolSessionTime: 4, JobFiles: 5, JobBytes: 6, Job: normalJob},
			"Errors: 1, JobType: B, JobStatus: T, JobLevel: F, JobID: 2, VolSessionID: 3, VolSessionTime: 4, JobFiles: 5, JobBytes: 6, StartTime: 1970-01-01 01:00:00 +0100 CET, EndTime: 1970-01-01 01:00:00 +0100 CET, Job: normal_name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			je := jobEntry{
				Errors:         tt.fields.Errors,
				JobType:        tt.fields.JobType,
				JobStatus:      tt.fields.JobStatus,
				JobLevel:       tt.fields.JobLevel,
				JobID:          tt.fields.JobID,
				VolSessionID:   tt.fields.VolSessionID,
				VolSessionTime: tt.fields.VolSessionTime,
				JobFiles:       tt.fields.JobFiles,
				JobBytes:       tt.fields.JobBytes,
				StartTime:      tt.fields.StartTime,
				EndTime:        tt.fields.EndTime,
				Job:            tt.fields.Job,
			}
			if got := je.String(); got != tt.want {
				t.Errorf("jobEntry.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseJobs(t *testing.T) {
	readerEmpty := bytes.NewReader([]byte{})
	readerTooSmall := bytes.NewReader([]byte{
		1, // number of jobs
	})
	readerJobError := bytes.NewReader([]byte{
		1, 0, 0, 0, // number of jobs
	})
	readerJobTooSmall := bytes.NewReader([]byte{
		1, 0, 0, 0, // number of jobs
		0,
	})
	readerInvalidJobName := bytes.NewReader([]byte{
		1, 0, 0, 0, // number of jobs
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // pad
		0, 0, 0, 0, // Errors
		'B', 0, 0, 0, // JobType
		'T', 0, 0, 0, // JobStatus
		0, 0, 0, 0, // JobLevel
		0, 0, 0, 0, // JobID
		0, 0, 0, 0, // VolSessionID
		0, 0, 0, 0, // VolSessionTime
		0, 0, 0, 0, // JobFiles
		0, 0, 0, 0, 0, 0, 0, 0, // JobBytes
		1, 0, 0, 0, 0, 0, 0, 0, // StartTime
		0, 0, 0, 0, 0, 0, 0, 0, // EndTime
		't', 'e', 's', 't', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // Job
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	readerZeroJobs := bytes.NewReader([]byte{
		0, 0, 0, 0, // number of jobs
	})
	readerOneNonBackupJob := bytes.NewReader([]byte{
		1, 0, 0, 0, // number of jobs
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // pad
		0, 0, 0, 0, // Errors
		'R', 0, 0, 0, // JobType
		'T', 0, 0, 0, // JobStatus
		0, 0, 0, 0, // JobLevel
		0, 0, 0, 0, // JobID
		0, 0, 0, 0, // VolSessionID
		0, 0, 0, 0, // VolSessionTime
		0, 0, 0, 0, // JobFiles
		0, 0, 0, 0, 0, 0, 0, 0, // JobBytes
		1, 0, 0, 0, 0, 0, 0, 0, // StartTime
		0, 0, 0, 0, 0, 0, 0, 0, // EndTime
		't', 'e', 's', 't', '.', '2', '0', '1', '2', '-', '0', '2', '-', '0', '1', 0, // Job
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	readerOneSuccessfulBackupJob := bytes.NewReader([]byte{
		1, 0, 0, 0, // number of jobs
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, // pad
		0, 0, 0, 0, // Errors
		'B', 0, 0, 0, // JobType
		'T', 0, 0, 0, // JobStatus
		0, 0, 0, 0, // JobLevel
		0, 0, 0, 0, // JobID
		0, 0, 0, 0, // VolSessionID
		0, 0, 0, 0, // VolSessionTime
		0, 0, 0, 0, // JobFiles
		0, 0, 0, 0, 0, 0, 0, 0, // JobBytes
		1, 0, 0, 0, 0, 0, 0, 0, // StartTime
		0, 0, 0, 0, 0, 0, 0, 0, // EndTime
		't', 'e', 's', 't', '.', '2', '0', '1', '2', '-', '0', '2', '-', '0', '1', 0, // Job
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	})
	type args struct {
		handle io.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantJobs []job.Job
		wantErr  bool
	}{
		{"read empty", args{readerEmpty}, nil, true},
		{"read too small", args{readerTooSmall}, nil, true},
		{"read job error", args{readerJobError}, nil, true},
		{"read job too small", args{readerJobTooSmall}, nil, true},
		{"read invalid job name", args{readerInvalidJobName}, nil, true},
		{"read zero jobs", args{readerZeroJobs}, nil, false},
		{"read one non backup job", args{readerOneNonBackupJob}, nil, false},
		{"read one successful backup job", args{readerOneSuccessfulBackupJob}, []job.Job{{Name: "test", Timestamp: 1, Success: true}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJobs, err := ParseJobs(tt.args.handle)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseJobs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotJobs, tt.wantJobs) {
				t.Errorf("ParseJobs() = %v, want %v", gotJobs, tt.wantJobs)
			}
		})
	}
}
