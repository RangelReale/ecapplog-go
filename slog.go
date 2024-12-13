package ecapplog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type SLogHandler struct {
	client *Client
	attrs  []slog.Attr
	groups []string
}

func NewSLogHandler(client *Client) *SLogHandler {
	return &SLogHandler{
		client: client,
	}
}

func (l *SLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (l *SLogHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := slogConverter(true, l.attrs, l.groups, &record)

	priority := Priority_INFORMATION
	switch record.Level {
	case slog.LevelDebug:
		priority = Priority_DEBUG
	case slog.LevelInfo:
		priority = Priority_INFORMATION
	case slog.LevelWarn:
		priority = Priority_WARNING
	case slog.LevelError:
		priority = Priority_ERROR
	}

	category := CategoryDEFAULT
	if cat, ok := slogcommon.FindAttribute(attrs, l.groups, "category"); ok {
		category = slogcommon.ValueToString(cat.Value)
	}

	var options []LogOption
	if len(attrs) > 0 {
		payload := slogcommon.AttrsToMap(attrs...)
		payloadEnc, err := json.Marshal(payload)
		if err != nil {
			options = append(options, WithSource(fmt.Sprintf("failed to marshal attributes: %v", err)))
		} else {
			options = append(options, WithSource(fmt.Sprintf("%s %s", record.Message, string(payloadEnc))))
		}
	}

	l.client.Log(record.Time, priority, category, record.Message, options...)

	return nil
}

func (l *SLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return l
	}
	return &SLogHandler{
		client: l.client,
		attrs:  slogcommon.AppendAttrsToGroup(l.groups, l.attrs, attrs...),
		groups: l.groups,
	}
}

func (l *SLogHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return l
	}

	return &SLogHandler{
		client: l.client,
		attrs:  l.attrs,
		groups: append(l.groups, name),
	}
}

var sourceKey = "source"
var errorKeys = []string{"error", "err"}

func slogConverter(addSource bool, loggerAttr []slog.Attr, groups []string, record *slog.Record) []slog.Attr {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	attrs = slogcommon.ReplaceError(attrs, errorKeys...)
	if addSource {
		attrs = append(attrs, slogcommon.Source(sourceKey, record))
	}
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	return attrs
}
