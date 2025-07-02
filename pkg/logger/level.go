package logger

type Level string

const (
	LEVEL_DEBUG Level = "debug"
	LEVEL_INFO  Level = "info"
	LEVEL_WARN  Level = "warn"
	LEVEL_ERROR Level = "error"
	LEVEL_FATAL Level = "fatal"
)

func (l Level) toInt() int {
	switch l {
	default:
		return 0
	case LEVEL_DEBUG:
		return 10
	case LEVEL_INFO:
		return 20
	case LEVEL_WARN:
		return 30
	case LEVEL_ERROR:
		return 40
	case LEVEL_FATAL:
		return 50
	}
}
