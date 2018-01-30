package template

import (
	logging "github.com/op/go-logging"
)

const logger = "gotemplate"

var log = logging.MustGetLogger(logger)

func getLogLevel() logging.Level {
	return logging.GetLevel(logger)
}

// SetLogLevel allows configuration of the default logging level
func SetLogLevel(level logging.Level) {
	logging.SetLevel(level, logger)
}
