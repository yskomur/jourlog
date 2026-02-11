package jourlog

import (
	"fmt"
	"time"

	"github.com/coreos/go-systemd/v22/sdjournal"
)

// JournalReader provides functionality to read logs from the systemd journal.
type JournalReader struct {
	j         *sdjournal.Journal // The underlying journal instance
	cursor    string             // Current cursor position in the journal
	limit     int                // Maximum number of entries to retrieve (0 means no limit)
	retrieved int                // Number of entries retrieved so far
}

// NewJournalReader creates a new journal reader instance.
// Returns an error if the journal cannot be opened.
func NewJournalReader() (*JournalReader, error) {
	journal, err := sdjournal.NewJournal()
	if err != nil {
		return nil, fmt.Errorf("failed to open journal: %w", err)
	}

	return &JournalReader{
		j:         journal,
		limit:     0, // No limit by default
		retrieved: 0,
	}, nil
}

// SetLimit sets the maximum number of entries to retrieve.
// A value of 0 means no limit.
func (j *JournalReader) SetLimit(limit int) {
	if limit < 0 {
		limit = 0 // Ensure limit is non-negative
	}
	j.limit = limit
}

// SetService configures the reader to only show entries from the specified service.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetService(service string) error {
	if err := j.j.AddMatch("__SYSTEMD_SERVICE=" + service); err != nil {
		return fmt.Errorf("failed to add service match: %w", err)
	}
	return nil
}

// SetUnit configures the reader to only show entries from the specified systemd unit.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetUnit(unit string) error {
	if err := j.j.AddMatch("_SYSTEMD_UNIT=" + unit); err != nil {
		return fmt.Errorf("failed to add unit match: %w", err)
	}
	return nil
}

// SetSince configures the reader to only show entries since the specified time.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetSince(since time.Time) error {
	timestamp := since.UnixNano() / 1000 // Convert to microseconds
	if err := j.j.AddMatch("_REALTIME_TIMESTAMP>=" + fmt.Sprintf("%d", timestamp)); err != nil {
		return fmt.Errorf("failed to add time match: %w", err)
	}
	return nil
}

// SetUntil configures the reader to only show entries until the specified time.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetUntil(until time.Time) error {
	timestamp := until.UnixNano() / 1000 // Convert to microseconds
	if err := j.j.AddMatch("_REALTIME_TIMESTAMP<=" + fmt.Sprintf("%d", timestamp)); err != nil {
		return fmt.Errorf("failed to add time match: %w", err)
	}
	return nil
}

// SetPriority configures the reader to only show entries with the specified priority or higher.
// Priority levels are: 0 (emergency), 1 (alert), 2 (critical), 3 (error), 4 (warning), 5 (notice), 6 (info), 7 (debug)
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetPriority(priority int) error {
	if priority < 0 || priority > 7 {
		return fmt.Errorf("invalid priority level: %d (must be 0-7)", priority)
	}
	if err := j.j.AddMatch("PRIORITY<=" + fmt.Sprintf("%d", priority)); err != nil {
		return fmt.Errorf("failed to add priority match: %w", err)
	}
	return nil
}

// AddFilter adds a custom match filter to the journal reader.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) AddFilter(filter string) error {
	if err := j.j.AddMatch(filter); err != nil {
		return fmt.Errorf("failed to add filter match: %w", err)
	}
	return nil
}

// SetHostname configures the reader to only show entries from the specified hostname.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetHostname(hostname string) error {
	if err := j.j.AddMatch("_HOSTNAME=" + hostname); err != nil {
		return fmt.Errorf("failed to add hostname match: %w", err)
	}
	return nil
}

// SetExecutable configures the reader to only show entries from the specified executable.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetExecutable(executable string) error {
	if err := j.j.AddMatch("_EXE=" + executable); err != nil {
		return fmt.Errorf("failed to add executable match: %w", err)
	}
	return nil
}

// SetMessageFilter configures the reader to only show entries containing the specified text in the message.
// Note: This is not a direct journalctl filter but uses the MESSAGE field.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) SetMessageFilter(text string) error {
	if err := j.j.AddMatch("MESSAGE=" + text); err != nil {
		return fmt.Errorf("failed to add message filter: %w", err)
	}
	return nil
}

// ClearFilters removes all match filters from the journal reader.
func (j *JournalReader) ClearFilters() {
	j.j.FlushMatches()
}

// LastHour is a convenience method that configures the reader to only show entries from the last hour.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) LastHour() error {
	return j.SetSince(time.Now().Add(-time.Hour))
}

// LastDay is a convenience method that configures the reader to only show entries from the last 24 hours.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) LastDay() error {
	return j.SetSince(time.Now().Add(-24 * time.Hour))
}

// LastWeek is a convenience method that configures the reader to only show entries from the last 7 days.
// Returns an error if the match filter cannot be added.
func (j *JournalReader) LastWeek() error {
	return j.SetSince(time.Now().Add(-7 * 24 * time.Hour))
}

// Today is a convenience method that configures the reader to only show entries from today (since midnight).
// Returns an error if the match filter cannot be added.
func (j *JournalReader) Today() error {
	now := time.Now()
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return j.SetSince(midnight)
}

// Retrieve fetches the next log entry from the journal.
// Returns the message content and any error encountered.
func (j *JournalReader) Retrieve() (string, error) {
	// Check if we've reached the limit
	if j.limit > 0 && j.retrieved >= j.limit {
		return "", fmt.Errorf("reached limit of %d entries", j.limit)
	}

	// Move to the next entry
	n, err := j.j.Next()
	if err != nil {
		return "", fmt.Errorf("failed to advance to next entry: %w", err)
	}

	// Check if we reached the end
	if n == 0 {
		return "", fmt.Errorf("no more entries")
	}

	// Get the entry
	log, err := j.j.GetEntry()
	if err != nil {
		return "", fmt.Errorf("failed to get entry: %w", err)
	}

	// Print all fields for debugging
	for f, rec := range log.Fields {
		fmt.Printf("%s: %s\n", f, rec)
	}

	// Update cursor position
	cursor, err := j.j.GetCursor()
	if err != nil {
		return "", fmt.Errorf("failed to get cursor: %w", err)
	}
	j.cursor = cursor

	// Increment the retrieved counter
	j.retrieved++

	// Return the message content
	message, ok := log.Fields["MESSAGE"]
	if !ok {
		return "", fmt.Errorf("entry has no MESSAGE field")
	}

	return message, nil
}

// ResetCounter resets the counter for retrieved entries.
// This is useful when you want to reuse the same reader with a different limit.
func (j *JournalReader) ResetCounter() {
	j.retrieved = 0
}

// SeekHead moves the cursor to the beginning of the journal.
// Returns an error if the seek operation fails.
func (j *JournalReader) SeekHead() error {
	if err := j.j.SeekHead(); err != nil {
		return fmt.Errorf("failed to seek to head: %w", err)
	}
	j.retrieved = 0
	return nil
}

// SeekTail moves the cursor to the end of the journal.
// Returns an error if the seek operation fails.
func (j *JournalReader) SeekTail() error {
	if err := j.j.SeekTail(); err != nil {
		return fmt.Errorf("failed to seek to tail: %w", err)
	}
	j.retrieved = 0
	return nil
}

// SeekCursor moves the cursor to the specified cursor position.
// Returns an error if the seek operation fails.
func (j *JournalReader) SeekCursor(cursor string) error {
	if err := j.j.SeekCursor(cursor); err != nil {
		return fmt.Errorf("failed to seek to cursor: %w", err)
	}
	j.retrieved = 0
	return nil
}

// GetCursor returns the current cursor position.
// Returns an error if the cursor cannot be retrieved.
func (j *JournalReader) GetCursor() (string, error) {
	cursor, err := j.j.GetCursor()
	if err != nil {
		return "", fmt.Errorf("failed to get cursor: %w", err)
	}
	return cursor, nil
}

// Close closes the journal reader and releases any resources.
// It should be called when the reader is no longer needed.
func (j *JournalReader) Close() error {
	if j.j != nil {
		return j.j.Close()
	}
	return nil
}
