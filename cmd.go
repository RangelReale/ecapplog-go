package ecapplog

import (
	"fmt"
	"time"
)

type command uint8

const (
	command_Banner command = 99
	command_Log    command = 0
)

type cmdLog struct {
	Time     cmdTime  `json:"time"`
	Priority Priority `json:"priority"`
	Category string   `json:"category"`
	Message  string   `json:"message"`
	Source   string   `json:"source,omitempty"`
}

const cmdTimeFormat = "2006-01-02T15:04:05.999999999"

type cmdTime struct {
	time.Time
}

func (t cmdTime) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", t.UTC().Format(cmdTimeFormat))
	return []byte(stamp), nil
}
