package users

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-utility/v2/array"
	"github.com/pkg/errors"
)

type updateRoleFunc func(roles []string, roleToUpdate string) []string

func AddUserRole(identityToken, stage, userID, role string) error {
	return AddUserRoleWithContext(context.Background(), identityToken, stage, userID, role)
}

func AddUserRoleWithContext(ctx context.Context, identityToken, stage, userID, role string) (err error) {
	return updateRoleToAllUsersNodes(ctx, identityToken, stage, userID, role, addRole)
}

func RemoveUserRole(identityToken, stage, userID, role string) error {
	return updateRoleToAllUsersNodes(context.Background(), identityToken, stage, userID, role, removeRole)
}

func updateRoleToAllUsersNodes(ctx context.Context, identityToken, stage string, userID string, role string, roleFunc updateRoleFunc) error {
	getNodesRequest := client.Get("/users/{id}/nodes-only").
		Assign("id", userID)

	restClient := httpClientAccessMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, getNodesRequest)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("wrong response status: %q", resp.Status)
	}

	gunhr := GetUserNodesHierarchiesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&gunhr); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	for _, node := range gunhr.Data {

		req := RoleRequest{
			Roles: roleFunc(node.Roles, role),
		}

		putRolesRequest := client.Put("/users/{id}/nodes/{nodeId}").
			Assign("id", userID).
			Assign("nodeId", node.ID).
			WithJSONPayload(req)

		resp, err := restClient.Do(ctx, putRolesRequest)
		if err != nil {
			return errors.Wrap(err, "failed to execute request")
		}

		if resp.StatusCode != http.StatusOK {
			return errors.Errorf("wrong response status: %q", resp.Status)
		}
	}

	return nil
}

func addRole(roles []string, newRole string) []string {
	if array.ContainsString(roles, newRole) {
		return roles
	}
	return append(roles, newRole)
}

func removeRole(roles []string, roleToBeRemoved string) []string {
	var newUserRoles = make([]string, 0, len(roles))
	for _, role := range roles {
		if role == roleToBeRemoved {
			continue
		}

		newUserRoles = append(newUserRoles, role)
	}
	return newUserRoles
}

type RoleRequest struct {
	Roles []string
}

type GetUserNodesHierarchiesResponse struct {
	Data []NodeHierarchy `json:"data"`
}
type NodeHierarchy struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Roles   []string `json:"roles"`
	SubType string   `json:"subType"`
	Type    string   `json:"type"`
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
