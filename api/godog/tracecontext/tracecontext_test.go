package tracecontext_test

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/SKF/go-tests-utility/api/godog/tracecontext"
)

type testScenarioContext struct {
	callbackBeforeScenario func(*godog.Scenario)
	callbackAfterScenario  func(*godog.Scenario, error)
	callbackBeforeStep     func(*godog.Step)
	callbackAfterStep      func(*godog.Step, error)
}

func (sc *testScenarioContext) BeforeScenario(fn func(*godog.Scenario)) {
	sc.callbackBeforeScenario = fn
}

func (sc *testScenarioContext) AfterScenario(fn func(*godog.Scenario, error)) {
	sc.callbackAfterScenario = fn
}

func (sc *testScenarioContext) BeforeStep(fn func(*godog.Step)) {
	sc.callbackBeforeStep = fn
}

func (sc *testScenarioContext) AfterStep(fn func(*godog.Step, error)) {
	sc.callbackAfterStep = fn
}

func (sc *testScenarioContext) Run(s *godog.Scenario, err error) {
	sc.callbackBeforeScenario(s)
	for _, step := range s.Steps {
		sc.callbackBeforeStep(step)
		sc.callbackAfterStep(step, nil)
	}
	sc.callbackAfterScenario(s, err)
}

func Test_TraceContextHappyCase(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long running tasks")
		return
	}

	os.Setenv(tracecontext.EnvTracerEnabled, "true")
	defer os.Unsetenv(tracecontext.EnvTracerEnabled)

	tracecontext.StartTracer(
		dd_tracer.WithService("tracecontext"),
		dd_tracer.WithEnv("local"),
		dd_tracer.WithAgentAddr("localhost:8126"),
	)
	defer tracecontext.StopTracer()

	sc := &testScenarioContext{}
	s := &godog.Scenario{
		Id:   "1",
		Name: "Scenario",
		Steps: []*messages.Pickle_PickleStep{
			{
				Id:   "1",
				Text: "Step 1",
			},
			{
				Id:   "2",
				Text: "Step 2",
			},
			{
				Id:   "3",
				Text: "Step 3",
			},
		},
		Tags: []*messages.Pickle_PickleTag{
			{
				Name: "ScenarioTag",
			},
		},
	}

	ctx := tracecontext.New(context.Background(), sc)
	sc.Run(s, nil)

	span, ok := dd_tracer.SpanFromContext(ctx)
	assert.True(t, ok)
	assert.NotEqual(t, uint64(0), span.Context().TraceID())
	assert.NotEqual(t, uint64(0), span.Context().SpanID())

	_, ok = dd_tracer.SpanFromContext(ctx)
	assert.True(t, ok)

	err := tracecontext.WriteTraceURLToFile(ctx, false, "trace.log", false, s.Name)
	assert.NoError(t, err)
}
