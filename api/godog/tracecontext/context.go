package tracecontext

import (
	"context"

	"github.com/cucumber/godog"
	dd_ext "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/ext"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/SKF/go-utility/v2/log"
)

const (
	spanType = "test"

	testOperationName = "testcase"
	tagTestID         = "test.id"
	tagTestTags       = "test.tags"

	stepOperationName = "step"
	tagStepID         = "step.id"
)

type GodogScenarioContext interface {
	BeforeStep(func(*godog.Step))
	AfterStep(func(*godog.Step, error))
	BeforeScenario(func(*godog.Scenario))
	AfterScenario(func(*godog.Scenario, error))
}

type traceCtx struct {
	context.Context
	rootContext context.Context

	rootSpan dd_tracer.Span
	stepSpan dd_tracer.Span

	// Workaround for this issue: https://github.com/cucumber/godog/issues/370
	stepErrorFound bool
}

func New(ctx context.Context, sc GodogScenarioContext) context.Context {
	if !Enabled() {
		log.Debug("Datadog tracing is not enabled")
		return ctx
	}

	tc := traceCtx{
		Context: ctx,
	}
	sc.BeforeScenario(tc.beforeScenario)
	sc.AfterScenario(tc.afterScenario)
	sc.BeforeStep(tc.beforeStep)
	sc.AfterStep(tc.afterStep)
	return &tc
}

func (tc *traceCtx) beforeScenario(s *godog.Scenario) {
	tc.stepErrorFound = false
	tc.rootSpan, tc.Context = dd_tracer.StartSpanFromContext(tc.Context, testOperationName,
		dd_tracer.SpanType(spanType),
	)
	tc.rootContext = tc.Context

	tc.rootSpan.SetTag(dd_ext.ResourceName, s.Name)
	tc.rootSpan.SetTag(tagTestID, s.Id)

	tags := []string{}
	for _, tag := range s.Tags {
		tags = append(tags, tag.Name)
	}
	tc.rootSpan.SetTag(tagTestTags, tags)
}

func (tc *traceCtx) afterScenario(s *godog.Scenario, err error) {
	tc.rootSpan.Finish(dd_tracer.WithError(err))
}

func (tc *traceCtx) beforeStep(s *godog.Step) {
	if tc.stepErrorFound {
		return
	}

	tc.stepSpan = dd_tracer.StartSpan(stepOperationName,
		dd_tracer.SpanType(spanType),
		dd_tracer.ChildOf(tc.rootSpan.Context()),
	)

	tc.stepSpan.SetTag(dd_ext.ResourceName, s.Text)
	tc.stepSpan.SetTag(tagStepID, s.Id)
	tc.Context = dd_tracer.ContextWithSpan(tc.Context, tc.stepSpan)
}

func (tc *traceCtx) afterStep(s *godog.Step, err error) {
	if err != nil {
		tc.stepErrorFound = true
	}

	tc.stepSpan.Finish(dd_tracer.WithError(err))
	tc.Context = tc.rootContext
}

func (tc *traceCtx) SetTag(key string, value interface{}) {
	SetTag(tc.Context, key, value)
}

func SetTag(ctx context.Context, key string, value interface{}) {
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		span.SetTag(key, value)
	}
}
