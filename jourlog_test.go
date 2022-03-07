package jourlog

import "testing"

func TestLog(t *testing.T) {
	Info("%s: info", "Info")
	Emerge("%s: info", "Emerge")
	Notice("%s: info", "Notice")
	Alert("%s: info", "Alert")
	Error("%s: info", "Error")
	Warning("%s: info", "Warning")
	Debug("%s: info", "Debug")
	Critical("%s: info", "Critical")
}
