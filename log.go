package charmlogr

import (
	"github.com/charmbracelet/log"
	"github.com/go-logr/logr"
)

type charmLogger struct {
	l                  *log.Logger
	verbosityFieldName string
	errorFieldName     string
	nameSeparator      string
	//loggerName         string
}

// Underlier exposes access to the underlying logging implementation.  Since
// callers only have a logr.Logger, they have to know which implementation is
// in use, so this interface is less of an abstraction and more of way to test
// type conversion.
type Underlier interface {
	GetUnderlying() *log.Logger
}

var (
	_ logr.LogSink                = &charmLogger{}
	_ Underlier                   = &charmLogger{}
	_ logr.CallStackHelperLogSink = &charmLogger{}
	_ logr.CallDepthLogSink       = &charmLogger{}
)

func (cl *charmLogger) copy() *charmLogger {
	return &charmLogger{
		l:                  cl.l.With(),
		verbosityFieldName: cl.verbosityFieldName,
		errorFieldName:     cl.errorFieldName,
		nameSeparator:      cl.nameSeparator,
		//loggerName:         cl.loggerName,
	}
}

func (cl *charmLogger) GetUnderlying() *log.Logger {
	return cl.l
}

func (cl *charmLogger) Init(_ logr.RuntimeInfo) {}

func (cl *charmLogger) Enabled(level int) bool {
	return level >= int(cl.l.GetLevel())
}

func (cl *charmLogger) Info(level int, msg string, keysAndValues ...interface{}) {
	cl.l.Helper()
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
	cl.l.Helper()

	if err != nil {
		keysAndValues = append([]any{cl.errorFieldName, err.Error()}, keysAndValues...)
	}
	cl.l.Error(msg, keysAndValues...)
}

func (cl *charmLogger) WithValues(keysAndValues ...interface{}) logr.LogSink {
	cl.l.Helper()
	copied := cl.copy()
	copied.l = copied.l.With(keysAndValues...)
	return copied
}

func (cl *charmLogger) WithName(name string) logr.LogSink {
	cl.l.Helper()
	newLogger := *cl
	oldName := cl.l.GetPrefix()
	if oldName != "" {
		oldName += cl.nameSeparator
	}
	oldName += name
	newLogger.l = cl.l.WithPrefix(oldName)
	return &newLogger
}

// Option is additional parameter for NewLoggerWithOptions.
type Option func(*charmLogger)

// WithVerbosityFieldName updates the field key for logr.Info verbosity, which by default is set to "v". If set to "",
// the verbosity key is not added to log line
func WithVerbosityFieldName(name string) Option {
	return func(logger *charmLogger) {
		logger.verbosityFieldName = name
	}
}

// WithErrorFieldName changes the default field name from "err"
func WithErrorFieldName(name string) Option {
	return func(logger *charmLogger) {
		logger.errorFieldName = name
	}
}

// WithNameSeparator changes the default separator of name parts. Default value is "/"
func WithNameSeparator(separator string) Option {
	return func(logger *charmLogger) {
		logger.nameSeparator = separator
	}
}

func NewLoggerWithOptions(l *log.Logger, options ...Option) logr.Logger {
	return logr.New(NewLogSinkWitOptions(l, options...))
}

func NewLogger(l *log.Logger) logr.Logger {
	return NewLoggerWithOptions(l)
}

func NewLogSink(l *log.Logger) logr.LogSink {
	return NewLogSinkWitOptions(l)
}

func NewLogSinkWitOptions(l *log.Logger, options ...Option) logr.LogSink {
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

func (cl *charmLogger) GetCallStackHelper() func() {
	return cl.l.Helper
}

func (cl *charmLogger) WithCallDepth(depth int) logr.LogSink {
	copied := cl.copy()
	copied.l.SetCallerOffset(depth)
	return copied
}
