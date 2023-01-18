package jourlog

import "testing"

func TestLog(t *testing.T) {
	if JLog.GetEcho() == false {
		t.Log("Ok")
	}
	t.Log("Bir test yazildi")
}
