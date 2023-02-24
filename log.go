package charmlogr

import (
	"github.com/charmbracelet/log"
	"github.com/go-logr/logr"
)

type charmLogger struct {
	l                  log.Logger
	verbosityFieldName string
	errorFieldName     string
	nameSeparator      string
	loggerName         string
}

// Underlier exposes access to the underlying logging implementation.  Since
// callers only have a logr.Logger, they have to know which implementation is
// in use, so this interface is less of an abstraction and more of way to test
// type conversion.
type Underlier interface {
	GetUnderlying() log.Logger
}

var (
	_ logr.LogSink = &charmLogger{}
	_ Underlier    = &charmLogger{}
)

func (cl *charmLogger) GetUnderlying() log.Logger {
	return cl.l
}

func (cl *charmLogger) Init(_ logr.RuntimeInfo) {}

func (cl *charmLogger) Enabled(level int) bool {
	return level >= int(cl.l.GetLevel())
}

func (cl *charmLogger) Info(level int, msg string, keysAndValues ...interface{}) {
	if cl.verbosityFieldName != "" {
		keysAndValues = append([]any{cl.verbosityFieldName, level}, keysAndValues...)
	}

	switch {
	case level <= 0:
		cl.l.Info(msg, keysAndValues...)
	default:
		cl.l.Debug(msg, keysAndValues...)
	}
}

func (cl *charmLogger) Error(err error, msg string, keysAndValues ...interface{}) {
	if err != nil {
		keysAndValues = append([]any{cl.errorFieldName, err.Error()}, keysAndValues...)
	}
	cl.l.Error(msg, keysAndValues...)
}

func (cl *charmLogger) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return &charmLogger{l: cl.l.With(keysAndValues...)}
}

func (cl *charmLogger) WithName(name string) logr.LogSink {
	newLogger := *cl
	newLogger.l = newLogger.l.With() // copies logger
	if newLogger.loggerName != "" {
		newLogger.loggerName += cl.nameSeparator
	}
	newLogger.loggerName += name
	newLogger.l.SetPrefix(newLogger.loggerName)
	return &newLogger
}

// option is additional parameter for NewLoggerWithOptions.
type option func(*charmLogger)

// WithVerbosityFieldName updates the field key for logr.Info verbosity, which by default is set to "v". If set to "",
// the verbosity key is added to log line
func WithVerbosityFieldName(name string) option {
	return func(logger *charmLogger) {
		logger.verbosityFieldName = name
	}
}

// WithErrorFieldName changes the default field name from "err"
func WithErrorFieldName(name string) option {
	return func(logger *charmLogger) {
		logger.errorFieldName = name
	}
}

// WithNameSeparator changes the default separator of name parts. Default value is "/"
func WithNameSeparator(separator string) option {
	return func(logger *charmLogger) {
		logger.nameSeparator = separator
	}
}

func NewLoggerWitOptions(l log.Logger, options ...option) logr.Logger {
	return logr.New(NewLogSinkWitOptions(l, options...))
}

func NewLogger(l log.Logger) logr.Logger {
	return NewLoggerWitOptions(l)
}

func NewLogSink(l log.Logger) logr.LogSink {
	return NewLogSinkWitOptions(l)
}

func NewLogSinkWitOptions(l log.Logger, options ...option) logr.LogSink {
	logger := &charmLogger{
		l:                  l,
		verbosityFieldName: "v",
		errorFieldName:     "err",
		nameSeparator:      "/",
	}
	for _, opt := range options {
		opt(logger)
	}
	return logger
}
