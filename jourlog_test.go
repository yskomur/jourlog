package jourlog

import "testing"

func TestLog(t *testing.T) {
	jLog := NewJourlog()
	if JLog.GetEcho() == jLog.GetEcho() {
		t.Log("Ok")
	}
	t.Log("Bir test yazildi")
}
