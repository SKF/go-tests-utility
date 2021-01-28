package users

import (
	"context"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"
)

func Delete(identityToken, stage, userID string) error {
	return DeleteWithContext(context.Background(), identityToken, stage, userID)
}

func DeleteWithContext(ctx context.Context, identityToken, stage, userID string) error {
	req := client.Delete("/users/{id}").
		Assign("id", userID).
		SetHeader(headers.ContentType, "application/json")

	restClient := httpClientIdentityMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.Errorf("wrong response status: %q", resp.Status)
	}

	return nil
}
