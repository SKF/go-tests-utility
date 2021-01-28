package tracecontext

import (
	"context"
	"errors"
	"fmt"
	"os"

	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	datadogHost = "datadoghq.eu"
)

func WriteTraceURLToFile(ctx context.Context, includeSpan bool, filename string, didFail bool, label string) error {
	url, ok := GetTraceURL(ctx, includeSpan)
	if !ok {
		return errors.New("couldn't get trace url")
	}

	verdict := "OK\t"
	if didFail {
		verdict = "FAIL"
	}

	data := fmt.Sprintf("%s\t%s\t%s\n", verdict, url, label)
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(data))
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}

func GetTraceURL(ctx context.Context, includeSpan bool) (string, bool) {
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		spanCtx := span.Context()

		url := fmt.Sprintf("https://%s/apm/trace/%d", datadogHost, spanCtx.TraceID())
		if spanCtx.TraceID() != spanCtx.SpanID() && includeSpan {
			url = fmt.Sprintf("%s?spanID=%d", url, spanCtx.SpanID())
		}
		return url, true
	}
	return "", false
}
