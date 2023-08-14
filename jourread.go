package jourlog

import (
	"github.com/coreos/go-systemd/v22/sdjournal"
)

type JournalReader struct {
	j      *sdjournal.Journal
	cursor string
}

func NewJournalReader() *JournalReader {
	return &JournalReader{
		j: &sdjournal.Journal{},
	}
}

func (j *JournalReader) LastHour() {
	_ = j.j.AddMatch("__SYSTEMD_SERVICE=ssh.service")
	_ = j.j.AddMatch("__REALTIME=now()-1hour :)")
}
func (j *JournalReader) Retrive() string {
	_, _ = j.j.Next()
	log, _ := j.j.GetEntry()
	for f, rec := range log.Fields {
		print(f, rec)
	}
	j.cursor, _ = j.j.GetCursor()
	return log.Fields["MESSAGE"]
}
