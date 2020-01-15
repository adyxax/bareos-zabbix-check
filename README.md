# Zabbix check for bareos backups

This repository contains code for a go program that can inspect a bareos status file to check the last run jobs. It outputs errors if a job's last run did not end successfully, or if a job is missing (ie it did not run). It should also be compatible with bacula.

This program was born from a need to query the status of the backups from the client machine and report it in zabbix at my workplace. Being a zabbix check it must exit with a code 0 even when reporting errors, be warned if you intend to use it with something else than zabbix. Changing this behaviour to suit your needs should not be hard at all though.

## Dependencies

go is required. Only go version >= 1.13.5 on linux amd64 has been tested.

## Building

For a debug build, use :
```
go build
```

For a release build, use :
```
go build -ldflags="-s -w"
```

## Usage

The common way to run this check is without any argument :
```
./bareos-zabbix-check
```

There are several flags available if you need to override the defaults :
  - -f string : Force the state file to use, defaults to bareos-fd.9102.state if it exists else bacula-fd.9102.state.
  - -q bool : Suppress all output, suitable to force a silent update of the spool file.
  - -v bool : Activates verbose debugging output, defaults to false.
  - -w string : Force the work directory to use, defaults to /var/lib/bareos if it exists else /var/lib/bacula.

## Output

As all zabbix checks, the program will exit 0 whatever happens. You will use the output in your triggers.

If there were no errors and there is no missing jobs, the program simply outputs : `OK`. The program outputs an `INFO <message>` if there were no backups ever (bootstrap situation mainly) or any special error. The program outputs an `AVERAGE <message>` if there was an error during the last run of a job, or if a job didn't run successfully in the last 24 hours.

Here is a list of the possible error messages and their meaning :
  - `AVERAGE: errors:%s missing:%s additionnal errors: %s` : there are backup errors or missing jobs.
  - `AVERAGE: Couldn't save spool : %s` : the program could not save its spool file in the work directory.
  - `INFO Invalid work directory %s : it does not exist or is not a directory.` : you manually specified a work directory with the `-w` flag and it is invalid.
  - `INFO Could not find a suitable work directory. Is bareos or bacula installed?` : neither /var/lib/bareos nor /var/lib/bacula seem to exist.
  - `INFO The state file %s does not exist.\n` : you manually specified a state file with the `-f` flag and it is invalid or does not exist in the working directory.
  - `INFO Could not find a suitable state file. Has a job ever run?` : neither bareos-fd.9102.state nor bacula-fd.9102.state seem to exist in the working directory.
  - `INFO Couldn't open state file : %s` : the bacula or bareos state file could not be opened.
  - `INFO Invalid state file : This script only supports bareos state file version 4, got %d` : The bacula or bareos version installed is not supported (yet!).
  - `INFO Corrupted state file : %s` : the bacula or bareos state file could not be parsed.
  - `INFO No jobs exist in the state file` : no jobs were found in the state file.
  - `INFO Couldn't parse job name, this shouldn't happen : %s` : the program uses a regex to strip time and date from a job entry and it did not work. This is a bug in this program! Please open an issue.

## Spool file

Stored in `/var/lib/bareos/bareos-zabbix-check.spool`, this spool data is a simple csv vile format where every line contains a job name and the timestamp of the last successful execution for this job.

## Limitations

### No alerts if a job fails to start on its first run

The Bareos file daemon holds no status reference for a job that never started properly. Therefore any director misconfiguration will not be caught up by this program unless the job ran successfully at least once. If it happened the job will have a status missing.

### False positives

Bareos status file only holds the last 10 jobs that ran on the host. This should be enough for nearly all use cases, but if a host has many jobs it won't do.

The solution to this is to have a `Client Run After Job` entry that runs this program after each job in order to have the program record that successful run in its spool.

### Missing job alert when you legitimately remove a job in the director's configuration

Because of the way we record jobs in a spool file in order to track missing jobs, if you remove a job in the director's configuration you will get a missing job alert the next day. To avoid this you just need to :
  - stop the bareos file daemon
  - delete the bareos file daemon status file (/var/lib/bareos/bareos-fd.9102.state by default)
  - start the bareos file daemon again
  - run any job in order to have the file daemon recreate a valid status file
  - delete the line referencing this job in the spool file (/var/lib/bareos/bareos-zabbix-check.spool by default)