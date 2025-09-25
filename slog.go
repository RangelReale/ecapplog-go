package ecapplog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	slogcommon "github.com/samber/slog-common"
)

type SLogHandler struct {
	client        *Client
	attrs         []slog.Attr
	groups        []string
	categoryKey   string
	customLevelFn func(slog.Level) Priority
	options       slog.HandlerOptions
}

func NewSLogHandler(client *Client, options ...SlogHandlerOption) *SLogHandler {
	ret := &SLogHandler{
		client:      client,
		categoryKey: categoryKey,
		options: slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		},
	}
	for _, option := range options {
		option(ret)
	}
	return ret
}

func (l *SLogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelDebug
	if l.options.Level != nil {
		minLevel = l.options.Level.Level()
	}
	return level >= minLevel
}

func (l *SLogHandler) Handle(ctx context.Context, record slog.Record) error {
	attrs := slogConverter(l.options, l.attrs, l.groups, &record)

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
	default:
		if l.customLevelFn != nil {
			priority = l.customLevelFn(record.Level)
		}
	}

	category := CategoryDEFAULT
	if cat, ok := slogcommon.FindAttribute(attrs, l.groups, l.categoryKey); ok {
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

type SlogHandlerOption func(*SLogHandler)

func WithSlogHandlerOptions(opts slog.HandlerOptions) SlogHandlerOption {
	return func(handler *SLogHandler) {
		handler.options = opts
	}
}

func WithWithSlogHandlerCustomLevelFn(fn func(slog.Level) Priority) SlogHandlerOption {
	return func(handler *SLogHandler) {
		handler.customLevelFn = fn
	}
}

func WithWithSlogHandlerCategoryKey(categoryKey string) SlogHandlerOption {
	return func(handler *SLogHandler) {
		handler.categoryKey = categoryKey
	}
}

var (
	sourceKey   = "source"
	categoryKey = "category"
	errorKeys   = []string{"error", "err"}
)

func slogConverter(options slog.HandlerOptions, loggerAttr []slog.Attr, groups []string, record *slog.Record) []slog.Attr {
	// aggregate all attributes
	attrs := slogcommon.AppendRecordAttrsToAttrs(loggerAttr, groups, record)

	// developer formatters
	if options.ReplaceAttr != nil {
		attrs = slogcommon.ReplaceAttrs(options.ReplaceAttr, groups, attrs...)
	}
	attrs = slogcommon.ReplaceError(attrs, errorKeys...)
	if options.AddSource {
		attrs = append(attrs, slogcommon.Source(sourceKey, record))
	}
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	return attrs
}
