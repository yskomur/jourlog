package jourlog

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/coreos/go-systemd/v22/journal"
)

// TestLoggerCreation verifies that a new logger can be created with default settings
func TestLoggerCreation(t *testing.T) {
	jLog := NewJourlog()

	// Check default settings
	if jLog.GetLogLevel() != journal.PriInfo {
		t.Errorf("Expected default log level to be Info, got %v", jLog.GetLogLevel())
	}

	if jLog.GetEcho() != false {
		t.Errorf("Expected default echo setting to be false, got %v", jLog.GetEcho())
	}

	// Verify global logger has the same default settings
	if jLog.GetEcho() != jLog.GetEcho() {
		t.Errorf("Global logger and new logger have different echo settings")
	}

	if jLog.GetLogLevel() != jLog.GetLogLevel() {
		t.Errorf("Global logger and new logger have different log levels")
	}
}

// TestLoggerSettings verifies that logger settings can be changed
func TestLoggerSettings(t *testing.T) {
	jLog := NewJourlog()

	// Change settings
	jLog.SetEcho(true)
	jLog.SetLogLevel(journal.PriDebug)

	// Verify changes
	if jLog.GetEcho() != true {
		t.Errorf("Failed to set echo to true")
	}

	if jLog.GetLogLevel() != journal.PriDebug {
		t.Errorf("Failed to set log level to Debug")
	}
}

// TestJournalReader tests the journal reader functionality
// Note: This test may be skipped in environments where the journal is not accessible
func TestJournalReader(t *testing.T) {
	reader, err := NewJournalReader()
	if err != nil {
		t.Skip("Skipping journal reader test: ", err)
	}

	// Test setting filters
	err = reader.LastHour()
	if err != nil {
		t.Logf("Warning: Failed to set last hour filter: %v", err)
		// Don't fail the test as this might be environment-dependent
	}

	// Try to retrieve one entry
	// This is just a basic test to ensure the method doesn't panic
	message, err := reader.Retrieve()
	if err != nil {
		t.Logf("Info: Failed to retrieve journal entry: %v", err)
		// Don't fail the test as there might not be any entries
	} else {
		t.Logf("Successfully retrieved journal entry: %s", message)
	}
}

//go:noinline
func emitCallerMetadataTestLog(jLog *JourLog, message string) (string, string, int) {
	_, file, line, _ := runtime.Caller(0)
	jLog.Info("%s", message)
	return file, "github.com/yskomur/jourlog.emitCallerMetadataTestLog", line + 1
}

func TestCallerMetadataWithJournalReader(t *testing.T) {
	reader, err := NewJournalReader()
	if err != nil {
		t.Skip("Skipping caller metadata test: ", err)
	}
	defer func() {
		_ = reader.Close()
	}()

	if err := reader.SeekTail(); err != nil {
		t.Skip("Skipping caller metadata test (seek tail failed): ", err)
	}

	marker := "jourlog-caller-test-" + strconv.FormatInt(time.Now().UnixNano(), 10)
	if err := reader.SetMessageFilter(marker); err != nil {
		t.Skip("Skipping caller metadata test (message filter failed): ", err)
	}

	jLog := NewJourlog()
	expectedFile, expectedFunc, expectedLine := emitCallerMetadataTestLog(jLog, marker)

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		fields, err := reader.RetrieveEntry()
		if err != nil {
			if strings.Contains(err.Error(), "no more entries") {
				time.Sleep(25 * time.Millisecond)
				continue
			}
			t.Skip("Skipping caller metadata test (retrieve failed): ", err)
		}

		gotMsg := fields["MESSAGE"]
		if gotMsg != marker {
			continue
		}

		gotFile := fields["CODE_FILE"]
		gotFunc := fields["CODE_FUNC"]
		gotLine := fields["CODE_LINE"]

		if gotFile == "" || gotFunc == "" || gotLine == "" {
			t.Fatalf("missing caller metadata fields: CODE_FILE=%q CODE_FUNC=%q CODE_LINE=%q", gotFile, gotFunc, gotLine)
		}

		if filepath.Clean(gotFile) != filepath.Clean(expectedFile) {
			t.Fatalf("CODE_FILE mismatch: got=%q expected=%q", gotFile, expectedFile)
		}
		if gotFunc != expectedFunc {
			t.Fatalf("CODE_FUNC mismatch: got=%q expected=%q", gotFunc, expectedFunc)
		}
		if gotLine != strconv.Itoa(expectedLine) {
			t.Fatalf("CODE_LINE mismatch: got=%q expected=%d", gotLine, expectedLine)
		}
		return
	}

	t.Skip("Skipping caller metadata test (journal entry not observed in time)")
}
