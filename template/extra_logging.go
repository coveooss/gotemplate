package template

import (
	"os"

	logging "github.com/op/go-logging"
)

const (
	logger         = "gotemplate"
	loggerInternal = "gotemplate-int"
	loggingBase    = "Logging"
)

var loggingFuncs = funcTableMap{
	"fatal":    {f: func(args ...interface{}) string { return logBase(Log.Fatal, args...) }, group: loggingBase, desc: ""},
	"fatalf":   {f: func(format string, args ...interface{}) string { return logBasef(Log.Fatalf, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
	"error":    {f: func(args ...interface{}) string { return logBase(Log.Error, args...) }, group: loggingBase, desc: ""},
	"errorf":   {f: func(format string, args ...interface{}) string { return logBasef(Log.Errorf, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
	"warning":  {f: func(args ...interface{}) string { return logBase(Log.Warning, args...) }, group: loggingBase, desc: ""},
	"warningf": {f: func(format string, args ...interface{}) string { return logBasef(Log.Warningf, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
	"notice":   {f: func(args ...interface{}) string { return logBase(Log.Notice, args...) }, group: loggingBase, desc: ""},
	"noticef":  {f: func(format string, args ...interface{}) string { return logBasef(Log.Noticef, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
	"info":     {f: func(args ...interface{}) string { return logBase(Log.Info, args...) }, group: loggingBase, desc: ""},
	"infof":    {f: func(format string, args ...interface{}) string { return logBasef(Log.Infof, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
	"debug":    {f: func(args ...interface{}) string { return logBase(Log.Debug, args...) }, group: loggingBase, desc: ""},
	"debugf":   {f: func(format string, args ...interface{}) string { return logBasef(Log.Debugf, format, args...) }, group: loggingBase, args: []string{"format"}, desc: ""},
}

func (t *Template) addLoggingFuncs() { t.AddFunctions(loggingFuncs) }

func logBase(f func(...interface{}), args ...interface{}) string {
	f(args...)
	return ""
}

func logBasef(f func(string, ...interface{}), format string, args ...interface{}) string {
	f(format, args...)
	return ""
}

// Log is the logger used to log message during template processing
var Log = logging.MustGetLogger(logger)

var log = logging.MustGetLogger(loggerInternal)

func getLogLevelInternal() logging.Level {
	return logging.GetLevel(loggerInternal)
}

// GetLogLevel returns the current logging level for gotemplate
func GetLogLevel() logging.Level {
	return logging.GetLevel(logger)
}

// SetLogLevel set the logging level for gotemplate
func SetLogLevel(level logging.Level) {
	logging.SetLevel(level, logger)
}

// InitLogging allows configuration of the default logging level
func InitLogging(level, internalLevel logging.Level, simple bool) {
	format := `[%{module}] %{time:2006/01/02 15:04:05.000} %{color}%{level:-8s} %{message}%{color:reset}`
	if simple {
		format = `[%{level}] %{message}`
	}
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(format)))
	SetLogLevel(level)
	logging.SetLevel(internalLevel, loggerInternal)
}

// Default package init
var _ = func() int {
	logging.SetLevel(logging.WARNING, loggerInternal)
	return 0
}()
