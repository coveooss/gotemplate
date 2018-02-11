package template

import (
	"os"

	logging "github.com/op/go-logging"
)

const logger = "gotemplate"

// Log is the logger used to log message during template processing
var Log = logging.MustGetLogger(logger)

func getLogLevel() logging.Level {
	return logging.GetLevel(logger)
}

// InitLogging allows configuration of the default logging level
func InitLogging(level logging.Level, simple bool) {
	format := `[%{module}] %{time:2006/01/02 15:04:05.000} %{color}%{level:-8s} %{message}%{color:reset}`
	if simple {
		format = `[%{level}] %{message}`
	}
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(format)))
	logging.SetLevel(level, logger)
}
