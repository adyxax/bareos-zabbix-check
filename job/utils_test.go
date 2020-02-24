package job

import (
	"reflect"
	"testing"
)

func TestKeepOldestOnly(t *testing.T) {
	emptyList := []Job{}
	oneJob := []Job{{Name: "a", Timestamp: 10, Success: true}}
	twoJobs := []Job{
		{Name: "a", Timestamp: 5, Success: true},
		{Name: "a", Timestamp: 10, Success: true},
	}
	threeJobs := []Job{
		{Name: "a", Timestamp: 5, Success: true},
		{Name: "a", Timestamp: 10, Success: true},
		{Name: "a", Timestamp: 8, Success: false},
	}
	type args struct {
		jobs []Job
	}
	tests := []struct {
		name string
		args args
		want []Job
	}{
		{"empty list", args{emptyList}, nil},
		{"one job", args{oneJob}, oneJob},
		{"two jobs", args{twoJobs}, oneJob},
		{"three jobs", args{threeJobs}, oneJob},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KeepOldestOnly(tt.args.jobs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KeepOldestOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeepSuccessOnly(t *testing.T) {
	emptyList := []Job{}
	oneJob := []Job{{Name: "a", Timestamp: 10, Success: true}}
	twoJobs := []Job{
		{Name: "a", Timestamp: 10, Success: true},
		{Name: "a", Timestamp: 5, Success: false},
	}
	type args struct {
		jobs []Job
	}
	tests := []struct {
		name       string
		args       args
		wantResult []Job
	}{
		{"empty list", args{emptyList}, emptyList},
		{"one job", args{oneJob}, oneJob},
		{"two jobs", args{twoJobs}, oneJob},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := KeepSuccessOnly(tt.args.jobs); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("KeepSuccessOnly() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
