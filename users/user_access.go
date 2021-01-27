package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SKF/go-utility/v2/log"
	"github.com/SKF/go-utility/v2/uuid"
	"github.com/pkg/errors"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const accessMgmtBaseURL = "https://api-web.%s.users.enlight.skf.com"

func AddUserAccess(identityToken, stage, userID, nodeID string) error {
	return AddUserAccessWithContext(context.Background(), identityToken, stage, userID, nodeID)
}

func AddUserAccessWithContext(ctx context.Context, identityToken, stage, userID, nodeID string) (err error) {
	log.Debugf("Adding access %s - %s", userID, nodeID)
	if !uuid.IsValid(userID) {
		return fmt.Errorf("Invalid User ID: %q", userID)
	}

	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s/hierarchies/%s", stage, userID, nodeID)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest failed: %w", err)
	}
	req = req.WithContext(ctx)
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		if err = dd_tracer.Inject(span.Context(), dd_tracer.HTTPHeadersCarrier(req.Header)); err != nil {
			return errors.Wrapf(err, "ddtracer.Inject: failed to inject trace headers")
		}
	}

	req.Header.Set("Authorization", identityToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Wrong status code: %q", resp.Status)
	}

	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}

	return nil
}
