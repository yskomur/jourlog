package jourlog

import (
	"fmt"
	"github.com/coreos/go-systemd/v22/journal"
	"runtime"
	"strconv"
)

// JournalLogger journal.Send wrapper
func JournalLogger(priority journal.Priority, format string, a ...interface{}) {
	vars := make(map[string]string)

	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()
	_ = journal.Send(fmt.Sprintf(format, a...), priority, vars)
}
