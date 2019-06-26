package template

import (
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/fatih/color"
	logging "github.com/op/go-logging"
)

const (
	logger         = "gotemplate"
	loggerInternal = "gotemplate-int"
	loggingBase    = "Logging"
)

var loggingFuncs = dictionary{
	"critical": func(args ...interface{}) string { return logBase(Log.Critical, args...) },
	"debug":    func(args ...interface{}) string { return logBase(Log.Debug, args...) },
	"error":    func(args ...interface{}) string { return logBase(Log.Error, args...) },
	"fatal":    func(args ...interface{}) string { return logBase(Log.Fatal, args...) },
	"info":     func(args ...interface{}) string { return logBase(Log.Info, args...) },
	"notice":   func(args ...interface{}) string { return logBase(Log.Notice, args...) },
	"panic":    func(args ...interface{}) string { return logBase(Log.Panic, args...) },
	"warning":  func(args ...interface{}) string { return logBase(Log.Warning, args...) },
}

var loggingFuncsAliases = aliases{
	"critical": {"criticalf"},
	"debug":    {"debugf"},
	"error":    {"errorf"},
	"fatal":    {"fatalf"},
	"info":     {"infof"},
	"notice":   {"noticef"},
	"panic":    {"panicf"},
	"warning":  {"warn", "warnf", "warningf"},
}

var loggingFuncsHelp = descriptions{
	"critical": "Logs a message using CRITICAL as log level (0).",
	"debug":    "Logs a message using DEBUG as log level (5).",
	"error":    "Logs a message using ERROR as log level (1).",
	"fatal":    "Equivalents to critical followed by a call to os.Exit(1).",
	"info":     "Logs a message using INFO as log level (4).",
	"notice":   "Logs a message using NOTICE as log level (3).",
	"panic":    "Equivalents to critical followed by a call to panic.",
	"warning":  "Logs a message using WARNING as log level (2).",
}

func (t *Template) addLoggingFuncs() {
	t.AddFunctions(loggingFuncs, loggingBase, FuncOptions{
		FuncHelp:    loggingFuncsHelp,
		FuncAliases: loggingFuncsAliases,
	})
}

func logBase(f func(...interface{}), args ...interface{}) string {
	f(utils.FormatMessage(args...))
	return ""
}

// Log is the logger used to log message during template processing
var Log = logging.MustGetLogger(logger)

// log is application logger used to follow the behaviour of the application
var log = logging.MustGetLogger(loggerInternal)

var loggingMutex sync.Mutex

func getLogLevelInternal() logging.Level {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	return logging.GetLevel(loggerInternal)
}

// GetLogLevel returns the current logging level for gotemplate
func GetLogLevel() logging.Level {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	return logging.GetLevel(logger)
}

// SetLogLevel set the logging level for gotemplate
func SetLogLevel(level logging.Level) {
	loggingMutex.Lock()
	defer loggingMutex.Unlock()
	logging.SetLevel(level, logger)
}

// ConfigureLogging allows configuration of the default logging level
func ConfigureLogging(level, internalLevel logging.Level, simple bool) {
	format := `[%{module}] %{time:2006/01/02 15:04:05.000} %{color}%{level:-8s} %{message}%{color:reset}`
	if simple {
		format = `[%{level}] %{message}`
	}
	logging.SetBackend(logging.NewBackendFormatter(logging.NewLogBackend(color.Error, "", 0), logging.MustStringFormatter(format)))
	SetLogLevel(level)
	logging.SetLevel(internalLevel, loggerInternal)
}

// InitLogging allows configuration of the default logging level
func InitLogging() int {
	if level, err := strconv.Atoi(utils.GetEnv(EnvDebug, "2")); err != nil {
		log.Warningf("Unable to convert %s into integer: %s", EnvDebug, os.Getenv(EnvDebug))
	} else {
		logging.SetLevel(logging.Level(level), loggerInternal)
	}
	return 0
}

// Default package init
var _ = InitLogging()

// TryGetLoggingLevelFromString converts a string into a logging level
func TryGetLoggingLevelFromString(level string, defaultLevel logging.Level) (logging.Level, error) {
	level = strings.TrimSpace(level)
	if level == "" {
		return defaultLevel, nil
	}

	levelNum, err := strconv.Atoi(level)
	if err == nil {
		return logging.Level(levelNum), nil
	}

	return logging.LogLevel(level)
}

// GetLoggingLevelFromString converts a string into a logging level
func GetLoggingLevelFromString(level string) logging.Level {
	return must(TryGetLoggingLevelFromString(level, logging.INFO)).(logging.Level)
}
