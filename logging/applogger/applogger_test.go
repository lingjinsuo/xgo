package applogger

import "testing"

func TestLogger(t *testing.T) {
	Info("%s", "aaaa")
	Debug("%s", "aaaa")
	Warn("%s", "aaaa")
	Error("%s", "aaaa")
	Fatal("%s", "aaaa")
}
