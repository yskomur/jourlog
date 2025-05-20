package jourlog

import (
	"context"
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
	if JLog.GetEcho() != jLog.GetEcho() {
		t.Errorf("Global logger and new logger have different echo settings")
	}

	if JLog.GetLogLevel() != jLog.GetLogLevel() {
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

// TestContextLogging tests the context-aware logging functionality
// This is a basic test that just ensures the methods don't panic
func TestContextLogging(t *testing.T) {
	jLog := NewJourlog()

	// Create a context with test values
	ctx := context.Background()
	ctx = context.WithValue(ctx, "request_id", "test-req-123")
	ctx = context.WithValue(ctx, "user_id", "test-user-456")

	// Test each log level with context
	t.Run("EmergWithContext", func(t *testing.T) {
		jLog.EmergeWithContext(ctx, "Test emergency log with context")
	})

	t.Run("AlertWithContext", func(t *testing.T) {
		jLog.AlertWithContext(ctx, "Test alert log with context")
	})

	t.Run("CriticalWithContext", func(t *testing.T) {
		jLog.CriticalWithContext(ctx, "Test critical log with context")
	})

	t.Run("ErrorWithContext", func(t *testing.T) {
		jLog.ErrorWithContext(ctx, "Test error log with context")
	})

	t.Run("WarningWithContext", func(t *testing.T) {
		jLog.WarningWithContext(ctx, "Test warning log with context")
	})

	t.Run("NoticeWithContext", func(t *testing.T) {
		jLog.NoticeWithContext(ctx, "Test notice log with context")
	})

	t.Run("InfoWithContext", func(t *testing.T) {
		jLog.InfoWithContext(ctx, "Test info log with context")
	})

	t.Run("DebugWithContext", func(t *testing.T) {
		jLog.DebugWithContext(ctx, "Test debug log with context")
	})
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
