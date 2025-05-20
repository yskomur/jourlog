// Package jourlog provides a Go interface for logging to systemd's journal.
// It supports all standard log levels and can be configured to echo logs to the console.
// The package also captures file, line, and function information for each log entry.
package jourlog

import (
	"context"
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
		tracebackDeep: 3,
	}
}

// JLog is a global logger instance that can be used directly.
var JLog = &JourLog{
	logLevel:      journal.PriInfo,
	printLog:      false,
	tracebackDeep: 3,
}

// SetEcho enables or disables printing log records to the console.
func (j *JourLog) SetEcho(state bool) {
	j.printLog = state
}

// SetLogLevel sets the minimum log level that will be processed.
// Messages with a lower priority will be ignored.
func (j *JourLog) SetLogLevel(l journal.Priority) {
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

// journalLogger is an internal function that sends a log message to the systemd journal.
// It captures the caller's file, line, and function name and includes them as metadata.
func (j *JourLog) journalLogger(priority journal.Priority, format string, a ...interface{}) {
	vars := make(map[string]string)

	// Capture caller information
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(j.tracebackDeep, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	// Add metadata
	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()

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

// journalLoggerWithContext is like journalLogger but accepts a context that can contain additional metadata.
func (j *JourLog) journalLoggerWithContext(ctx context.Context, priority journal.Priority, format string, a ...interface{}) {
	vars := make(map[string]string)

	// Capture caller information
	pc := make([]uintptr, 10)
	runtime.Callers(j.tracebackDeep, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])

	// Add metadata
	vars["CODE_FILE"] = file
	vars["CODE_LINE"] = strconv.Itoa(line)
	vars["CODE_FUNC"] = f.Name()

	// Extract metadata from context if available
	if ctx != nil {
		if reqID, ok := ctx.Value("request_id").(string); ok {
			vars["REQUEST_ID"] = reqID
		}
		if userID, ok := ctx.Value("user_id").(string); ok {
			vars["USER_ID"] = userID
		}
		// Add any other context values you might want to extract
	}

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

// EmergeWithContext logs a message with emergency priority and includes context metadata.
func (j *JourLog) EmergeWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriEmerg {
		j.journalLoggerWithContext(ctx, journal.PriEmerg, format, a...)
	}
}

// Alert logs a message with alert priority.
// This level should be used for conditions that require immediate action.
func (j *JourLog) Alert(format string, a ...interface{}) {
	if j.logLevel >= journal.PriAlert {
		j.journalLogger(journal.PriAlert, format, a...)
	}
}

// AlertWithContext logs a message with alert priority and includes context metadata.
func (j *JourLog) AlertWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriAlert {
		j.journalLoggerWithContext(ctx, journal.PriAlert, format, a...)
	}
}

// Critical logs a message with critical priority.
// This level should be used for critical conditions like hard device errors.
func (j *JourLog) Critical(format string, a ...interface{}) {
	if j.logLevel >= journal.PriCrit {
		j.journalLogger(journal.PriCrit, format, a...)
	}
}

// CriticalWithContext logs a message with critical priority and includes context metadata.
func (j *JourLog) CriticalWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriCrit {
		j.journalLoggerWithContext(ctx, journal.PriCrit, format, a...)
	}
}

// Error logs a message with error priority.
// This level should be used for error conditions that don't require immediate action.
func (j *JourLog) Error(format string, a ...interface{}) {
	if j.logLevel >= journal.PriErr {
		j.journalLogger(journal.PriErr, format, a...)
	}
}

// ErrorWithContext logs a message with error priority and includes context metadata.
func (j *JourLog) ErrorWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriErr {
		j.journalLoggerWithContext(ctx, journal.PriErr, format, a...)
	}
}

// Warning logs a message with warning priority.
// This level should be used for warning conditions that might require attention.
func (j *JourLog) Warning(format string, a ...interface{}) {
	if j.logLevel >= journal.PriWarning {
		j.journalLogger(journal.PriWarning, format, a...)
	}
}

// WarningWithContext logs a message with warning priority and includes context metadata.
func (j *JourLog) WarningWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriWarning {
		j.journalLoggerWithContext(ctx, journal.PriWarning, format, a...)
	}
}

// Notice logs a message with notice priority.
// This level should be used for normal but significant conditions.
func (j *JourLog) Notice(format string, a ...interface{}) {
	if j.logLevel >= journal.PriNotice {
		j.journalLogger(journal.PriNotice, format, a...)
	}
}

// NoticeWithContext logs a message with notice priority and includes context metadata.
func (j *JourLog) NoticeWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriNotice {
		j.journalLoggerWithContext(ctx, journal.PriNotice, format, a...)
	}
}

// Info logs a message with informational priority.
// This level should be used for informational messages.
func (j *JourLog) Info(format string, a ...interface{}) {
	if j.logLevel >= journal.PriInfo {
		j.journalLogger(journal.PriInfo, format, a...)
	}
}

// InfoWithContext logs a message with informational priority and includes context metadata.
func (j *JourLog) InfoWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriInfo {
		j.journalLoggerWithContext(ctx, journal.PriInfo, format, a...)
	}
}

// Debug logs a message with debug priority.
// This level should be used for detailed debug information.
func (j *JourLog) Debug(format string, a ...interface{}) {
	if j.logLevel >= journal.PriDebug {
		j.journalLogger(journal.PriDebug, format, a...)
	}
}

// DebugWithContext logs a message with debug priority and includes context metadata.
func (j *JourLog) DebugWithContext(ctx context.Context, format string, a ...interface{}) {
	if j.logLevel >= journal.PriDebug {
		j.journalLoggerWithContext(ctx, journal.PriDebug, format, a...)
	}
}
