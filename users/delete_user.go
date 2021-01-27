package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Delete(accessToken, stage, userID string) error {
	return DeleteWithContext(context.Background(), accessToken, stage, userID)
}

func DeleteWithContext(ctx context.Context, accessToken, stage, userID string) error {
	url := fmt.Sprintf(identityMgmtBaseURL+"/users/%s", stage, userID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
	}
	req = req.WithContext(ctx)
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		if err = dd_tracer.Inject(span.Context(), dd_tracer.HTTPHeadersCarrier(req.Header)); err != nil {
			return errors.Wrapf(err, "ddtracer.Inject: failed to inject trace headers")
		}
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.Errorf("wrong status: %q", resp.Status)
	}

	return nil
}
