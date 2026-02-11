// Package jourlog provides structured logging and journal reading utilities for systemd-based
// systems.
//
// # Logging
//
// Create a logger with NewJourlog, choose the minimum accepted level with SetLogLevel,
// and write records with priority-specific methods such as Info, Error, or Debug.
//
//	logger := jourlog.NewJourlog()
//	logger.SetLogLevel(journal.PriInfo)
//	logger.Info("service started: pid=%d", os.Getpid())
//
// Each record is sent to systemd-journald with CODE_FILE, CODE_LINE, and CODE_FUNC
// metadata fields derived from the caller.
//
// # Journal reading
//
// Open a reader with NewJournalReader, apply optional filters, then retrieve entries
// sequentially with Retrieve.
//
//	reader, err := jourlog.NewJournalReader()
//	if err != nil {
//		// handle error
//	}
//	defer reader.Close()
//
//	_ = reader.LastHour()
//	for {
//		message, err := reader.Retrieve()
//		if err != nil {
//			break
//		}
//		fmt.Println(message)
//	}
//
// The reader supports common filters such as SetUnit, SetPriority, SetHostname,
// and convenience time windows like LastHour, LastDay, and Today.
package jourlog
