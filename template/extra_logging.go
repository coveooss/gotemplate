package template

import (
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/coveooss/multilogger"
	"github.com/sirupsen/logrus"
)

const (
	loggingBase = "Logging"
)

var (
	// templateLog is the logger used to log message during template processing
	templateLog *multilogger.MultiLogger = multilogger.New(logrus.InfoLevel, multilogger.DisabledLevel, "", "gotemplate")
	// InternalLog is application logger used to follow the behaviour of the application
	InternalLog *multilogger.MultiLogger = multilogger.New(logrus.InfoLevel, multilogger.DisabledLevel, "", "gotemplate")
)

var loggingFuncs = dictionary{
	"trace":    func(args ...interface{}) string { return logBase(templateLog.Trace, args...) },
	"debug":    func(args ...interface{}) string { return logBase(templateLog.Debug, args...) },
	"info":     func(args ...interface{}) string { return logBase(templateLog.Info, args...) },
	"notice":   func(args ...interface{}) string { return logBase(templateLog.Info, args...) },
	"warning":  func(args ...interface{}) string { return logBase(templateLog.Warning, args...) },
	"error":    func(args ...interface{}) string { return logBase(templateLog.Error, args...) },
	"critical": func(args ...interface{}) string { return logBase(templateLog.Fatal, args...) },
	"fatal":    func(args ...interface{}) string { return logBase(templateLog.Fatal, args...) },
	"panic":    func(args ...interface{}) string { return logBase(templateLog.Panic, args...) },
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

// ConfigureLogging allows configuration of the default logging level
func ConfigureLogging(templateLevel, internalLevel, internalLogFileLevel, internalLogFilePath string) {
	templateLog.SetConsoleLevel(templateLevel)
	InternalLog.SetConsoleLevel(internalLevel)
	InternalLog.ConfigureFileLogger(internalLogFileLevel, internalLogFilePath)
}
