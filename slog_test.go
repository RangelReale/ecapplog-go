package ecapplog

import "testing"

func TestSlogHandler(t *testing.T) {
	client := NewClient()
	handler := NewSLogHandler(client)
	if handler == nil {
		t.Error("handler is nil")
	}
}
