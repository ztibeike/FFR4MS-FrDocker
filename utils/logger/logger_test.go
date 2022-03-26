package logger

import "testing"

func TestLogger(t *testing.T) {
	Info("test")
	Info("test%s", "1")
}
