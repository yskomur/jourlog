package jourlog

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/coreos/go-systemd/v22/journal"
)

// JourLog represents a logger instance that writes to systemd's journal.
// It can be configured with different log levels and options.
type JourLog struct {
	logLevel      journal.Priority // Current log level threshold
	printLog      bool             // Whether to echo logs to console
	tracebackDeep int              // How many stack frames to skip when capturing caller info
}

// NewJourlog creates a new logger instance with default settings.
// By default, the log level is set to Info, console output is disabled,
// and the traceback depth is set to 3.
func NewJourlog() *JourLog {
	return &JourLog{
		logLevel:      journal.PriInfo,
		printLog:      false,
		tracebackDeep: 2,
	}
}

// SetLogLevel sets the minimum log level that will be processed.
// Messages with a lower priority will be ignored.
func (j *JourLog) SetLogLevel(l journal.Priority) {
	if l < journal.PriAlert || l > journal.PriDebug {
		panic("invalid log level")
	}
	j.logLevel = l
}

// GetLogLevel returns the current log level threshold.
func (j *JourLog) GetLogLevel() journal.Priority {
	return j.logLevel
}

// GetEcho returns whether console output is enabled.
func (j *JourLog) GetEcho() bool {
	return j.printLog
}

// SetEcho enables or disables printing log records to the console.
func (j *JourLog) SetEcho(state bool) {
	j.printLog = state
}

// getCaller captures caller information from the stack
func (j *JourLog) getCaller() (string, string, string) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(j.tracebackDeep+1, pc) // +1 because we're in a helper function
	if n > 0 {
		f := runtime.FuncForPC(pc[0])
		if f != nil {
			file, line := f.FileLine(pc[0])
			return file, strconv.Itoa(line), f.Name()
		}
	}
	return "unknown", "0", "unknown"
}

// journalLogger logs a message with the specified priority and captures caller info
func (j *JourLog) journalLogger(priority journal.Priority, format string, a ...interface{}) {
	file, line, function := j.getCaller()

	vars := make(map[string]string)
	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = line
	vars["CODE_FUNC"] = function

	// Format the message
	record := fmt.Sprintf(format, a...)

	// Send to journal
	err := journal.Send(record, priority, vars)
	if err != nil {
		// If journal logging fails, at least print to console
		fmt.Printf("Failed to log to journal: %v\n", err)
	}

	// Echo to console if enabled
	if j.printLog {
		fmt.Println(record)
	}
}

// Emerge logs a message with emergency priority.
// This level should be used for panic conditions that require immediate attention.
func (j *JourLog) Emerge(format string, a ...interface{}) {
	if j.logLevel >= journal.PriEmerg {
		j.journalLogger(journal.PriEmerg, format, a...)
	}
}

// Alert logs a message with alert priority.
// This level should be used for conditions that require immediate action.
func (j *JourLog) Alert(format string, a ...interface{}) {
	if j.logLevel >= journal.PriAlert {
		j.journalLogger(journal.PriAlert, format, a...)
	}
}

// Critical logs a message with critical priority.
// This level should be used for critical conditions like hard device errors.
func (j *JourLog) Critical(format string, a ...interface{}) {
	if j.logLevel >= journal.PriCrit {
		j.journalLogger(journal.PriCrit, format, a...)
	}
}

// Error logs a message with error priority.
// This level should be used for error conditions that don't require immediate action.
func (j *JourLog) Error(format string, a ...interface{}) {
	if j.logLevel >= journal.PriErr {
		j.journalLogger(journal.PriErr, format, a...)
	}
}

// Warning logs a message with warning priority.
// This level should be used for warning conditions that might require attention.
func (j *JourLog) Warning(format string, a ...interface{}) {
	if j.logLevel >= journal.PriWarning {
		j.journalLogger(journal.PriWarning, format, a...)
	}
}

// Notice logs a message with notice priority.
// This level should be used for normal but significant conditions.
func (j *JourLog) Notice(format string, a ...interface{}) {
	if j.logLevel >= journal.PriNotice {
		j.journalLogger(journal.PriNotice, format, a...)
	}
}

// Info logs a message with informational priority.
// This level should be used for informational messages.
func (j *JourLog) Info(format string, a ...interface{}) {
	if j.logLevel >= journal.PriInfo {
		j.journalLogger(journal.PriInfo, format, a...)
	}
}

// Debug logs a message with debug priority.
// This level should be used for detailed debug information.
func (j *JourLog) Debug(format string, a ...interface{}) {
	if j.logLevel >= journal.PriDebug {
		j.journalLogger(journal.PriDebug, format, a...)
	}
}
