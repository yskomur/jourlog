package jourlog

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"runtime"
	"strconv"
)

// JournalLogger journal.Send wrapper
func journalLogger(priority journal.Priority, format string, a ...interface{}) {
	vars := make(map[string]string)

	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(3, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()
	_ = journal.Send(fmt.Sprintf(format, a...), priority, vars)
}

func Emerge(format string, a ...interface{}) {
	journalLogger(journal.PriEmerg, format, a)
}

func Notice(format string, a ...interface{}) {
	journalLogger(journal.PriNotice, format, a)
}

func Info(format string, a ...interface{}) {
	journalLogger(journal.PriInfo, format, a)
}

func Alert(format string, a ...interface{}) {
	journalLogger(journal.PriAlert, format, a)
}

func Error(format string, a ...interface{}) {
	journalLogger(journal.PriErr, format, a)
}

func Critical(format string, a ...interface{}) {
	journalLogger(journal.PriCrit, format, a)
}
