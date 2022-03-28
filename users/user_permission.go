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

func RemoveUserRoleWithContext(ctx context.Context, identityToken, stage, userID, roleToBeRemoved string) (err error) {
	return updateRoleToAllUsersNodes(ctx, identityToken, stage, userID, roleToBeRemoved, removeRole)
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

	gunhr := getUserNodesHierarchiesResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&gunhr); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	for _, node := range gunhr.Data.Nodes {

		req := roleRequest{
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

		if resp.StatusCode != http.StatusAccepted {
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

type roleRequest struct {
	Roles []string
}

type getUserNodesHierarchiesResponse struct {
	Data nodeHierarchies `json:"data"`
}

type nodeHierarchies struct {
	Nodes []nodeHierarchy `json:"nodes"`
}
type nodeHierarchy struct {
	ID    string   `json:"id"`
	Roles []string `json:"roles"`
}
