package jourlog

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"runtime"
	"strconv"
)

type JourLog struct {
	logLevel      journal.Priority
	printLog      bool
	tracebackDeep int
}

var JLog = &JourLog{
	logLevel:      journal.PriInfo,
	printLog:      false,
	tracebackDeep: 3,
}

// SetEcho print log record console
func (j *JourLog) SetEcho(state bool) {
	j.printLog = state
}

// SetLogLevel set new loglevel
func (j *JourLog) SetLogLevel(l journal.Priority) {
	j.logLevel = l
}

// GetLogLevel return current loglevel
func (j *JourLog) GetLogLevel() journal.Priority {
	return j.logLevel
}

// GetEcho return current state
func (j *JourLog) GetEcho() bool {
	return j.printLog
}

// JournalLogger journal.Send wrapper
func (j *JourLog) journalLogger(priority journal.Priority, format string, a ...interface{}) {
	var record string
	vars := make(map[string]string)

	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(j.tracebackDeep, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()
	record = fmt.Sprintf(format, a...)
	_ = journal.Send(record, priority, vars)
	if j.printLog {
		fmt.Println(record)
	}
}

func (j *JourLog) Emerge(format string, a ...interface{}) {
	if j.logLevel >= journal.PriEmerg {
		j.journalLogger(journal.PriEmerg, format, a...)
	}
}

func (j *JourLog) Alert(format string, a ...interface{}) {
	if j.logLevel >= journal.PriAlert {
		j.journalLogger(journal.PriAlert, format, a...)
	}
}

func (j *JourLog) Critical(format string, a ...interface{}) {
	if j.logLevel >= journal.PriCrit {
		j.journalLogger(journal.PriCrit, format, a...)
	}
}

func (j *JourLog) Error(format string, a ...interface{}) {
	if j.logLevel >= journal.PriErr {
		j.journalLogger(journal.PriErr, format, a...)
	}
}

func (j *JourLog) Warning(format string, a ...interface{}) {
	if j.logLevel >= journal.PriWarning {
		j.journalLogger(journal.PriWarning, format, a...)
	}
}

func (j *JourLog) Notice(format string, a ...interface{}) {
	if j.logLevel >= journal.PriNotice {
		j.journalLogger(journal.PriNotice, format, a...)
	}
}

func (j *JourLog) Info(format string, a ...interface{}) {
	if j.logLevel >= journal.PriInfo {
		j.journalLogger(journal.PriInfo, format, a...)
	}
}

func (j *JourLog) Debug(format string, a ...interface{}) {
	if j.logLevel >= journal.PriDebug {
		j.journalLogger(journal.PriDebug, format, a...)
	}
}
