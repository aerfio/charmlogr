package charmlogr_test

import (
	"errors"
	"os"

	"github.com/charmbracelet/log"

	"github.com/aerfio/charmlogr"
)

func Example() {
	l := charmlogr.NewLogger(log.New(log.WithOutput(os.Stdout), log.WithLevel(log.DebugLevel)))
	l.Info("info msg")
	l.WithName("loggerName").V(1).Info("debug message")
	l.Error(errors.New("whoops"), "additional msg", "key", "value")
	l.Error(nil, "no error but err level")

	// Output:
	// INFO info msg v=0
	// DEBUG loggerName: debug message v=1
	// ERROR additional msg err=whoops key=value
	// ERROR no error but err level
}

func ExampleNewLogger() {
	l := charmlogr.NewLoggerWitOptions(
		log.New(log.WithOutput(os.Stdout)), // default info lvl
		charmlogr.WithErrorFieldName("error"),
		charmlogr.WithVerbosityFieldName("verbosity"),
	)
	l.V(1).Info("does not get logged")
	l.Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// INFO some message verbosity=0
	// ERROR log line with error error=whoops
}

func ExampleNewLogger_logfmt_formatter() {
	l := charmlogr.NewLoggerWitOptions(
		log.New(log.WithOutput(os.Stdout), log.WithFormatter(log.LogfmtFormatter)),
	)
	l.WithName("logfmt-logger").Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// lvl=info prefix=logfmt-logger: msg="some message" v=0
	// lvl=error msg="log line with error" err=whoops
}

func ExampleNewLogger_json_formatter() {
	l := charmlogr.NewLoggerWitOptions(
		log.New(log.WithOutput(os.Stdout), log.WithFormatter(log.JSONFormatter)),
	)
	l.WithName("json-logger").Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// {"lvl":"info","msg":"some message","prefix":"json-logger:","v":0}
	// {"err":"whoops","lvl":"error","msg":"log line with error"}
}
