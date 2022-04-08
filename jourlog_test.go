package jourlog

import "testing"

func TestLog(t *testing.T) {
	log := NewJourlog()
	if log.GetEcho() == false {
		t.Log("Ok")
	}
	t.Log("Bir test yazildi")
}
