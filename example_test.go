package charmlogr_test

import (
	"errors"
	"os"
	"time"

	"github.com/charmbracelet/log"

	"github.com/aerfio/charmlogr"
)

func Example() {
	l := charmlogr.NewLogger(log.NewWithOptions(os.Stdout, log.Options{
		Level: log.DebugLevel,
	}))
	l.Info("info msg")
	l.WithName("loggerName").V(1).Info("debug message")
	l.Error(errors.New("whoops"), "additional msg", "key", "value")
	l.Error(nil, "no error but err level")

	// Output:
	// INFO info msg v=0
	// DEBU loggerName: debug message v=1
	// ERRO additional msg err=whoops key=value
	// ERRO no error but err level
}

func ExampleNewLoggerWithOptions() {
	l := charmlogr.NewLoggerWithOptions(
		log.New(os.Stdout), // default info lvl
		charmlogr.WithErrorFieldName("error"),
		charmlogr.WithVerbosityFieldName("verbosity"),
	)
	l.V(1).Info("does not get logged")
	l.Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// INFO some message verbosity=0
	// ERRO log line with error error=whoops
}

func ExampleNewLoggerWithOptions_LogFmtFormatter() {
	l := charmlogr.NewLoggerWithOptions(
		log.NewWithOptions(os.Stdout, log.Options{
			Formatter: log.LogfmtFormatter,
		}),
	)
	l.WithName("logfmt-logger").Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// lvl=info prefix=logfmt-logger msg="some message" v=0
	// lvl=error msg="log line with error" err=whoops
}

func ExampleNewLoggerWithOptions_JSONFormatter() {
	l := charmlogr.NewLoggerWithOptions(
		log.NewWithOptions(os.Stdout, log.Options{Formatter: log.JSONFormatter}),
	)
	l.WithName("json-logger").Info("some message")
	l.Error(errors.New("whoops"), "log line with error")
	// Output:
	// {"lvl":"info","msg":"some message","prefix":"json-logger","v":0}
	// {"err":"whoops","lvl":"error","msg":"log line with error"}
}

func ExampleNewLoggerWithOptions_MoreOptions() {
	l := charmlogr.NewLoggerWithOptions(
		log.NewWithOptions(os.Stdout, log.Options{
			TimeFunction: func() time.Time {
				return time.Date(1996, time.March, 24, 1, 2, 3, 4, time.UTC)
			},
			Prefix:          "test-prefix",
			ReportTimestamp: true,
			ReportCaller:    true,
			CallerOffset:    0,
			Fields:          []any{"key-pair", 1, "another-key", "value-for-that"},
			Formatter:       log.LogfmtFormatter,
		}),
	)
	l.
		WithName("json-logger").
		Info("some message")
	l.
		Error(errors.New("whoops"), "log line with error")
	// Output:
	// ts="1996/03/24 01:02:03" lvl=info caller=charmlogr/example_test.go:83 prefix=test-prefix/json-logger msg="some message" key-pair=1 another-key=value-for-that v=0
	// ts="1996/03/24 01:02:03" lvl=error caller=charmlogr/example_test.go:85 prefix=test-prefix msg="log line with error" key-pair=1 another-key=value-for-that err=whoops
}
