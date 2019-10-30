package template

import (
	"github.com/coveooss/gotemplate/v3/utils"
	"github.com/coveooss/multilogger"
)

const (
	loggingBase = "Logging"
)

var (
	// TemplateLog is the logger used to log message during template processing
	TemplateLog = multilogger.New("gotemplate")
	// InternalLog is application logger used to follow the behaviour of the application
	InternalLog = multilogger.New("gotemplate-int")
)

var loggingFuncs = dictionary{
	"trace":    func(args ...interface{}) string { return logBase(TemplateLog.Trace, args...) },
	"debug":    func(args ...interface{}) string { return logBase(TemplateLog.Debug, args...) },
	"info":     func(args ...interface{}) string { return logBase(TemplateLog.Info, args...) },
	"notice":   func(args ...interface{}) string { return logBase(TemplateLog.Info, args...) },
	"warning":  func(args ...interface{}) string { return logBase(TemplateLog.Warning, args...) },
	"error":    func(args ...interface{}) string { return logBase(TemplateLog.Error, args...) },
	"critical": func(args ...interface{}) string { return logBase(TemplateLog.Fatal, args...) },
	"fatal":    func(args ...interface{}) string { return logBase(TemplateLog.Fatal, args...) },
	"panic":    func(args ...interface{}) string { return logBase(TemplateLog.Panic, args...) },
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
