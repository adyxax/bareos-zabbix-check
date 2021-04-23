package zabbix

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(t *testing.T) {
	os.RemoveAll("tmp")
	cwd, _ := os.Getwd()
	err := os.MkdirAll("tmp/ok-18.2", 0777)
	if err != nil {
		t.Skipf("skipping main tests because tmp directory cannot be created : %s", err)
	}
	wd, err := os.Getwd()
	if err != nil {
		t.Skipf("skipping main tests because cannot get working directory : %s", err)
	}
	os.MkdirAll("tmp/ok-17.2", 0777)
	os.MkdirAll("tmp/no_state_file", 0777)
	os.MkdirAll("tmp/bacula_auto_detect_failed/var/lib/bacula", 0777)
	os.MkdirAll("tmp/bareos_auto_detect_failed/var/lib/bareos", 0777)
	os.MkdirAll("tmp/error", 0777)
	os.Symlink("../../testdata/bareos-fd-17.2.state", "tmp/ok-17.2/state")
	os.Symlink("../../testdata/bareos-fd-18.2.state", "tmp/ok-18.2/state")
	os.Symlink("../../testdata/bareos-fd-18.2.state-with-error", "tmp/error/state")
	tests := []struct {
		name      string
		timestamp uint64
		rootDir   string
		args      []string
		want      string
	}{
		{"failed bacula_auto_detect", 0, "tmp/bacula_auto_detect_failed", []string{}, "INFO Failed to init programm : Could not find a suitable state file. Has a job ever run?"},
		{"failed bareos_auto_detect", 0, "tmp/bareos_auto_detect_failed", []string{}, "INFO Failed to init programm : Could not find a suitable state file. Has a job ever run?"},
		{"failed auto_detect", 0, "tmp/non_existent", []string{}, "INFO Failed to init programm : Could not find a suitable work directory. Is bareos or bacula installed?"},
		{"no work directory", 0, "tmp", []string{"-w", "/non_existent"}, fmt.Sprintf("INFO Failed to init programm : Invalid work directory %s/tmp/non_existent : it does not exist or is not a directory", wd)},
		{"no state file auto_detect", 0, "tmp", []string{"-w", "/no_state_file"}, "INFO Failed to init programm : Could not find a suitable state file. Has a job ever run?"},
		{"no state file", 0, "tmp", []string{"-w", "/no_state_file", "-f", "test"}, fmt.Sprintf("INFO Failed to init programm : The state file %s/tmp/no_state_file/test does not exist", wd)},
		{"ok bareos 18.2", 1582579731, "tmp/ok-18.2", []string{"-w", "/", "-f", "state"}, "OK"},
		{"ok bareos 17.2", 1582579731, "tmp/ok-17.2", []string{"-w", "/", "-f", "state"}, "OK"},
		{"missing", 1582709331, "tmp/ok-18.2", []string{"-w", "/", "-f", "state"}, "AVERAGE:  missing: awhphpipam1_percona_xtrabackup, awhphpipam1_LinuxAll, awhphpipam1_www"},
		{"error", 1582579731, "tmp/error", []string{"-w", "/", "-f", "state"}, "AVERAGE: errors: awhphpipam1_percona_xtrabackup, awhphpipam1_www  Corrupted spool file: invalid argument"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now = tt.timestamp
			root = filepath.Join(cwd, tt.rootDir)
			flag.CommandLine = flag.NewFlagSet("bareos-zabbix-check", flag.ExitOnError) //flags are now reset
			os.Args = append([]string{"bareos-zabbix-check"}, tt.args...)
			if got := Main(); got != tt.want {
				t.Log(workDir)
				t.Errorf("Main() = %v, want %v", got, tt.want)
			}
		})
	}
	os.RemoveAll("tmp")
}
