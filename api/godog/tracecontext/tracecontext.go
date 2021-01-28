package tracecontext

import (
	"os"

	"github.com/cucumber/godog"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	EnvTracerEnabled = "DATADOG_TRACING_ENABLED"
)

func Enabled() bool {
	return os.Getenv(EnvTracerEnabled) == "true"
}

func StartTracer(opts ...dd_tracer.StartOption) {
	if Enabled() {
		dd_tracer.Start(opts...)
	}
}

func StopTracer() {
	if Enabled() {
		dd_tracer.Stop()
	}
}

func RegisterSuiteHooks(ts *godog.TestSuiteContext, opts ...dd_tracer.StartOption) {
	ts.BeforeSuite(func() {
		StartTracer(opts...)
	})

	ts.AfterSuite(StopTracer)
}
