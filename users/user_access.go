package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"
	"github.com/SKF/go-utility/v2/log"
	"github.com/SKF/go-utility/v2/uuid"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"
)

const accessMgmtBaseURL = "https://access-api.%s.users.enlight.skf.com"

func httpClientAccessMgmt(stage, identityToken string) *client.Client {
	return client.NewClient(
		client.WithBaseURL(fmt.Sprintf(accessMgmtBaseURL, stage)),
		client.WithDatadogTracing(),
		client.WithTokenProvider(auth.RawToken(identityToken)),
	)
}

func AddUserAccess(identityToken, stage, userID, nodeID string) error {
	return AddUserAccessWithContext(context.Background(), identityToken, stage, userID, nodeID)
}

func AddUserAccessWithContext(ctx context.Context, identityToken, stage, userID, nodeID string) (err error) {
	log.Debugf("Adding access %s - %s", userID, nodeID)
	if !uuid.IsValid(userID) {
		return fmt.Errorf("Invalid User ID: %q", userID)
	}

	req := client.Put("/users/{userId}/nodes/{nodeId}").
		Assign("userId", userID).
		Assign("nodeId", nodeID).
		SetHeader(headers.ContentType, "application/json")

	restClient := httpClientAccessMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("wrong response status: %q", resp.Status)
	}

	return nil
}

func RemoveUserAccess(identityToken, stage, userID, nodeID string) error {
	return RemoveUserAccessWithContext(context.Background(), identityToken, stage, userID, nodeID)
}

func RemoveUserAccessWithContext(ctx context.Context, identityToken, stage, userID, nodeID string) (err error) {
	log.Debugf("Removing access %s - %s", userID, nodeID)
	if !uuid.IsValid(userID) {
		return fmt.Errorf("Invalid User ID: %q", userID)
	}

	req := client.Delete("/users/{userId}/nodes/{nodeId}").
		Assign("userId", userID).
		Assign("nodeId", nodeID).
		SetHeader(headers.ContentType, "application/json")

	restClient := httpClientAccessMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("wrong response status: %q", resp.Status)
	}

	return nil
}
