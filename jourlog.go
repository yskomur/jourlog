package jourlog

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"runtime"
	"strconv"
)

type JourLog struct {
	printLog      bool
	tracebackDeep int
}

// NewJourlog Create logger new instance
func NewJourlog() *JourLog {
	return &JourLog{
		printLog:      false,
		tracebackDeep: 3,
	}
}

// SetEcho print log record console
func (j *JourLog) SetEcho(state bool) {
	j.printLog = state
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
	j.journalLogger(journal.PriEmerg, format, a...)
}

func (j *JourLog) Notice(format string, a ...interface{}) {
	j.journalLogger(journal.PriNotice, format, a...)
}

func (j *JourLog) Warning(format string, a ...interface{}) {
	j.journalLogger(journal.PriWarning, format, a...)
}

func (j *JourLog) Debug(format string, a ...interface{}) {
	j.journalLogger(journal.PriDebug, format, a...)
}

func (j *JourLog) Info(format string, a ...interface{}) {
	j.journalLogger(journal.PriInfo, format, a...)
}

func (j *JourLog) Alert(format string, a ...interface{}) {
	j.journalLogger(journal.PriAlert, format, a...)
}

func (j *JourLog) Error(format string, a ...interface{}) {
	j.journalLogger(journal.PriErr, format, a...)
}

func (j *JourLog) Critical(format string, a ...interface{}) {
	j.journalLogger(journal.PriCrit, format, a...)
}
