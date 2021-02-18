package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"
)

func AddUserRole(identityToken, stage, userID, role string) error {
	return AddUserRoleWithContext(context.Background(), identityToken, stage, userID, role)
}

func AddUserRoleWithContext(ctx context.Context, identityToken, stage, userID, role string) (err error) {
	user, err := getUser(ctx, identityToken, stage, userID)
	if err != nil {
		return
	}

	user.UserRoles = append(user.UserRoles, role)
	return updateUser(ctx, identityToken, stage, user)
}

func RemoveUserRole(identityToken, stage, userID, role string) error {
	return RemoveUserRoleWithContext(context.Background(), identityToken, stage, userID, role)
}

func RemoveUserRoleWithContext(ctx context.Context, identityToken, stage, userID, roleToBeRemoved string) (err error) {
	user, err := getUser(ctx, identityToken, stage, userID)
	if err != nil {
		return
	}

	var newUserRoles = make([]string, 0, len(user.UserRoles))
	for _, role := range user.UserRoles {
		if role == roleToBeRemoved {
			continue
		}

		newUserRoles = append(newUserRoles, role)
	}

	if len(newUserRoles) == len(user.UserRoles) {
		// Nothing to update
		return
	}

	user.UserRoles = newUserRoles
	return updateUser(ctx, identityToken, stage, user)
}

func getUser(ctx context.Context, identityToken, stage, userID string) (user user, err error) {
	if userID == "" {
		return user, fmt.Errorf("userID is required")
	}

	req := client.Get("/users/{id}").
		Assign("id", userID).
		SetHeader(headers.ContentType, "application/json")

	restClient := httpClientAccessMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "failed to execute request")
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("wrong response status: %q", resp.Status)
		return
	}

	if err = resp.Unmarshal(&user); err != nil {
		err = errors.Wrap(err, "failed to unmarshal body")
		return
	}

	return user, err
}

func updateUser(ctx context.Context, identityToken, stage string, user user) error {
	req := client.Put("/users/{id}").
		Assign("id", user.ID).
		WithJSONPayload(user)

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

type user struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	CompanyID      string   `json:"companyId"`
	UserRoles      []string `json:"userRoles"`
	Username       string   `json:"username"`
	UserStatus     string   `json:"userStatus"`
	EulaAgreedDate string   `json:"eulaAgreedDate"`
	ValidEula      bool     `json:"validEula"`
	Firstname      string   `json:"firstname"`
	Surname        string   `json:"surname"`
	Locale         string   `json:"locale"`
	CreatedDate    string   `json:"createdDate"`
	UpdatedDate    string   `json:"updatedDate"`
}
