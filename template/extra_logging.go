package template

import (
	"os"
	"strconv"

	"github.com/coveo/gotemplate/utils"
	logging "github.com/op/go-logging"
)

const (
	// DebugEnvVar is the name of the environment variable used to set the debug logging level
	DebugEnvVar = "GOTEMPLATE_DEBUG"

	logger         = "gotemplate"
	loggerInternal = "gotemplate-int"
	loggingBase    = "Logging"
)

var loggingFuncs = dictionary{
	"critical":  func(args ...interface{}) string { return logBase(Log.Critical, args...) },
	"criticalf": func(format string, args ...interface{}) string { return logBasef(Log.Criticalf, format, args...) },
	"debug":     func(args ...interface{}) string { return logBase(Log.Debug, args...) },
	"debugf":    func(format string, args ...interface{}) string { return logBasef(Log.Debugf, format, args...) },
	"error":     func(args ...interface{}) string { return logBase(Log.Error, args...) },
	"errorf":    func(format string, args ...interface{}) string { return logBasef(Log.Errorf, format, args...) },
	"fatal":     func(args ...interface{}) string { return logBase(Log.Fatal, args...) },
	"fatalf":    func(format string, args ...interface{}) string { return logBasef(Log.Fatalf, format, args...) },
	"info":      func(args ...interface{}) string { return logBase(Log.Info, args...) },
	"infof":     func(format string, args ...interface{}) string { return logBasef(Log.Infof, format, args...) },
	"notice":    func(args ...interface{}) string { return logBase(Log.Notice, args...) },
	"noticef":   func(format string, args ...interface{}) string { return logBasef(Log.Noticef, format, args...) },
	"panic":     func(args ...interface{}) string { return logBase(Log.Panic, args...) },
	"panicf":    func(format string, args ...interface{}) string { return logBasef(Log.Panicf, format, args...) },
	"warning":   func(args ...interface{}) string { return logBase(Log.Warning, args...) },
	"warningf":  func(format string, args ...interface{}) string { return logBasef(Log.Warningf, format, args...) },
}

var loggingFuncsArgs = arguments{
	"criticalf": {"format", "args"},
	"debugf":    {"format", "args"},
	"errorf":    {"format", "args"},
	"fatalf":    {"format", "args"},
	"infof":     {"format", "args"},
	"noticef":   {"format", "args"},
	"panicf":    {"format", "args"},
	"warningf":  {"format", "args"},
}

var loggingFuncsAliases = aliases{
	"warning":  {"warn"},
	"warningf": {"warnf"},
}

var loggingFuncsHelp = descriptions{
	"critical":  "Logs a message using CRITICAL as log level (0).",
	"criticalf": "Logs a message with format string using CRITICAL as log level (0).",
	"debug":     "Logs a message using DEBUG as log level (5).",
	"debugf":    "Logs a message with format using DEBUG as log level (5).",
	"error":     "Logs a message using ERROR as log level (1).",
	"errorf":    "Logs a message with format using ERROR as log level (1).",
	"fatal":     "Equivalents to critical followed by a call to os.Exit(1).",
	"fatalf":    "Equivalents to criticalf followed by a call to os.Exit(1).",
	"info":      "Logs a message using INFO as log level (4).",
	"infof":     "Logs a message with format using INFO as log level (4).",
	"notice":    "Logs a message using NOTICE as log level (3).",
	"noticef":   "Logs a message with format using NOTICE as log level (3).",
	"panic":     "Equivalents to critical followed by a call to panic.",
	"panicf":    "Equivalents to criticalf followed by a call to panic.",
	"warning":   "Logs a message using WARNING as log level (2).",
	"warningf":  "Logs a message with format using WARNING as log level (2).",
}

func (t *Template) addLoggingFuncs() {
	t.AddFunctions(loggingFuncs, loggingBase, funcOptions{
		funcHelp:    loggingFuncsHelp,
		funcArgs:    loggingFuncsArgs,
		funcAliases: loggingFuncsAliases,
	})
}

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

// log is application logger used to follow the behaviour of the application
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

// ConfigureLogging allows configuration of the default logging level
func ConfigureLogging(level, internalLevel logging.Level, simple bool) {
	format := `[%{module}] %{time:2006/01/02 15:04:05.000} %{color}%{level:-8s} %{message}%{color:reset}`
	if simple {
		format = `[%{level}] %{message}`
	}
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(os.Stderr, "", 0), logging.MustStringFormatter(format)))
	SetLogLevel(level)
	logging.SetLevel(internalLevel, loggerInternal)
}

// InitLogging allows configuration of the default logging level
func InitLogging() int {
	if level, err := strconv.Atoi(utils.GetEnv(DebugEnvVar, "2")); err != nil {
		log.Warningf("Unable to convert %s into integer: %s", DebugEnvVar, os.Getenv(DebugEnvVar))
	} else {
		logging.SetLevel(logging.Level(level), loggerInternal)
	}
	return 0
}

// Default package init
var _ = InitLogging()
