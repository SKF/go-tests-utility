# Tracing with godog
**tracecontext** is a package used for creating traces for godog scenarios aswell as steps. It will use godog After- and Beforehooks to start and stop traces.

**tracecontext** is implemented using the `context.Context` interface.

## Usage
To enable the tracer the environment variable `DATADOG_TRACING_ENABLED` needs to be set to `true`

The hardcoded values in the example might be good to fetch from environment variables or input arguments.
Like _dd_tracer.WithAgentAddr_, where you can configure which address the traces should be sent to, might be different than localhost if the code is running inside containers.


```go
import (
    "context"
    ...

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/SKF/go-tests-utility/api/godog/tracecontext"
	ddtracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type myFeature struct {
    ctx   context.Context
    items []string
}

func TestMain(m *testing.M) {
    status := godog.TestSuite{
    	status := godog.TestSuite{
		Name:                 "my-test-suite",
		ScenarioInitializer:  InitializeScenario,
		TestSuiteInitializer: InitializeTestSuite,
	}.Run()

	os.Exit(status)
}

func InitializeTestSuite(ts *godog.TestSuiteContext) {
    // Register hooks for starting the tracer
    // when a new suite is created
	tracecontext.RegisterSuiteHooks(ts,
		ddtracer.WithService("my-service"),
		ddtracer.WithEnv("sandbox"),
		ddtracer.WithServiceVersion("release-123"),
		ddtracer.WithAgentAddr("localhost:8126"),
	)
}

func InitializeScenario(s *godog.ScenarioContext) {
    // This will publish new traces when a scenerio is finished
    ctx := tracecontext.New(context.Background(), s)

    // A way of passing the context to the steps
    feature := myFeature{
        ctx: ctx,
        items: []string{},
    }

    s.Step(`^add "([^"]*)" to shopping cart$`, feature.addToCart)
    s.Step(`^go to checkout$`, feature.checkout)
}

func (f *myFeature) addToCart(item string) error {
    f.items = append(f.items, item)
    return nil
}

func (f *myFeature) checkout() error {
    // Here the ctx can be reached from f.ctx and used to pass it forward
    // or add additional information to its span

    tracecontext.SetTag(f.ctx, "count", len(f.items))
    
    ...
}
```

It can come handy to output all trace urls to a file,
all you will have to do to make this possible is to add a hook to your scenario.
```go
...
s.AfterScenario(func(s *godog.Scenario, err error) {
    tracecontext.WriteTraceURLToFile(ctx, false, "trace.log", err != nil, s.Name)
})
...
```