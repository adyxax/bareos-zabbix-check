package job

import (
	"testing"
)

func TestJob_String(t *testing.T) {
	type fields struct {
		Name      string
		Timestamp uint64
		Success   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"default job", fields{}, "Job { Name: \"\", Timestamp: \"0\", Success: \"false\" }"},
		{"a job", fields{Name: "a", Timestamp: 10, Success: true}, "Job { Name: \"a\", Timestamp: \"10\", Success: \"true\" }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			job := Job{
				Name:      tt.fields.Name,
				Timestamp: tt.fields.Timestamp,
				Success:   tt.fields.Success,
			}
			if got := job.String(); got != tt.want {
				t.Errorf("Job.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
