package zabbix

import (
	"flag"
)

var (
	stateFileName string
	workDir       string
)

func processFlags() (err error) {
	flag.StringVar(&stateFileName, "f", "", "Force the state file to use, defaults to "+bareosStateFile+" if it exists else "+baculaStateFile+".")
	flag.StringVar(&workDir, "w", "", "Force the work directory to use, defaults to "+bareosWorkDir+" if it exists else "+baculaWorkDir+".")
	flag.Parse()
	err = checkWorkDir()
	if err == nil {
		err = checkStateFile()
	}
	return
}
