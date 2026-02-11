package jourlog_test

import (
	"fmt"
	"log"

	"github.com/coreos/go-systemd/v22/journal"
	"github.com/yskomur/jourlog"
)

func ExampleNewJourlog() {
	logger := jourlog.NewJourlog()
	logger.SetEcho(true)
	logger.SetLogLevel(journal.PriInfo)
	logger.Info("service started")
}

func ExampleNewJournalReader() {
	reader, err := jourlog.NewJournalReader()
	if err != nil {
		log.Printf("journal unavailable: %v", err)
		return
	}
	defer reader.Close()

	if err := reader.LastHour(); err != nil {
		log.Printf("unable to apply time filter: %v", err)
		return
	}

	message, err := reader.Retrieve()
	if err != nil {
		log.Printf("unable to retrieve entry: %v", err)
		return
	}

	fmt.Println(message)
}
