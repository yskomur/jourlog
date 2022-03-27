package jourlog

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"runtime"
	"strconv"
)

var printLog bool
var tracebackDeep = 3

// SetEcho print log record console
func SetEcho(state bool) {
	printLog = state
}

// GetEcho return current state
func GetEcho() bool {
	return printLog
}

// JournalLogger journal.Send wrapper
func journalLogger(priority journal.Priority, format string, a ...interface{}) {
	var record string
	vars := make(map[string]string)

	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(tracebackDeep, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()
	record = fmt.Sprintf(format, a...)
	_ = journal.Send(record, priority, vars)
	if printLog {
		fmt.Println(record)
	}
}

func Emerge(format string, a ...interface{}) {
	journalLogger(journal.PriEmerg, format, a...)
}

func Notice(format string, a ...interface{}) {
	journalLogger(journal.PriNotice, format, a...)
}

func Warning(format string, a ...interface{}) {
	journalLogger(journal.PriWarning, format, a...)
}

func Debug(format string, a ...interface{}) {
	journalLogger(journal.PriDebug, format, a...)
}

func Info(format string, a ...interface{}) {
	journalLogger(journal.PriInfo, format, a...)
}

func Alert(format string, a ...interface{}) {
	journalLogger(journal.PriAlert, format, a...)
}

func Error(format string, a ...interface{}) {
	journalLogger(journal.PriErr, format, a...)
}

func Critical(format string, a ...interface{}) {
	journalLogger(journal.PriCrit, format, a...)
}

func init() {
	SetEcho(false)
	GetEcho()
}
