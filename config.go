package ecapplog

type Priority string

const (
	Priority_TRACE       Priority = "TRACE"
	Priority_DEBUG       Priority = "DEBUG"
	Priority_INFORMATION Priority = "INFORMATION"
	Priority_NOTICE      Priority = "NOTICE"
	Priority_WARNING     Priority = "WARNING"
	Priority_FATAL       Priority = "FATAL"
	Priority_CRITICAL    Priority = "CRITICAL"
	Priority_ERROR       Priority = "ERROR"
)
