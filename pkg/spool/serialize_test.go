package spool

import (
	"bytes"
	"testing"

	"git.adyxax.org/adyxax/bareos-zabbix-check/pkg/job"
)

func TestSerialize(t *testing.T) {
	type args struct {
		jobs []job.Job
	}
	tests := []struct {
		name       string
		args       args
		wantHandle string
		wantErr    bool
	}{
		{"One job", args{[]job.Job{{Name: "a", Timestamp: 1}}}, "a,1\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handle := &bytes.Buffer{}
			if err := Serialize(handle, tt.args.jobs); (err != nil) != tt.wantErr {
				t.Errorf("Serialize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHandle := handle.String(); gotHandle != tt.wantHandle {
				t.Errorf("Serialize() = %v, want %v", gotHandle, tt.wantHandle)
			}
		})
	}
}
