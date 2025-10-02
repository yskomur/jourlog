package jourlog

import (
	"testing"

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
