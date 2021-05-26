package hierarchy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"
	"github.com/go-http-utils/headers"
	"github.com/pkg/errors"
)

const (
	hierarchyBaseURL = "https://api.%s.hierarchy.enlight.skf.com"
	companyType      = "company"
)

func httpClient(stage, identityToken string) *client.Client {
	return client.NewClient(
		client.WithBaseURL(fmt.Sprintf(hierarchyBaseURL, stage)),
		client.WithDatadogTracing(),
		client.WithTokenProvider(auth.RawToken(identityToken)),
	)
}

func CreateCompany(identityToken, stage, parentNodeID, label, description string) (_ string, err error) {
	return CreateCompanyWithContext(context.Background(), identityToken, stage, parentNodeID, label, description)
}

func CreateCompanyWithContext(ctx context.Context, identityToken, stage, parentNodeID, label, description string) (_ string, err error) {
	return CreateWithContext(ctx, identityToken, stage, parentNodeID, label, description, companyType, companyType)
}

func Create(identityToken, stage, parentNodeID, label, description, nodetype, subtype string) (_ string, err error) {
	return CreateWithContext(context.Background(), identityToken, stage, parentNodeID, label, description, nodetype, subtype)
}

func CreateWithContext(ctx context.Context, identityToken, stage, parentNodeID, label, description, nodetype, subtype string) (_ string, err error) {
	requestBody := struct {
		ParentID    string `json:"parentId"`
		Label       string `json:"label"`
		Description string `json:"description"`
		Type        string `json:"nodeType"`
		SubType     string `json:"nodeSubType"`
		Criticality string `json:"criticality"`
	}{
		ParentID:    parentNodeID,
		Label:       label,
		Description: description,
		Type:        nodetype,
		SubType:     subtype,
	}

	if requestBody.Type == "asset" {
		requestBody.Criticality = "criticality_b"
	}

	req := client.Post("/nodes").
		WithJSONPayload(requestBody)

	restClient := httpClient(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "failed to execute request")
		return
	}

	var responseBody struct {
		ID string `json:"nodeId"`
	}
	if err = resp.Unmarshal(&responseBody); err != nil {
		err = errors.Wrap(err, "failed to unmarshal response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("wrong response status: %q", resp.Status)
		return
	}

	return responseBody.ID, nil
}

func Delete(identityToken, stage, nodeID string) error {
	return DeleteWithContext(context.Background(), identityToken, stage, nodeID)
}

func DeleteWithContext(ctx context.Context, identityToken, stage, nodeID string) error {
	req := client.Delete("/nodes/{id}").
		Assign("id", nodeID).
		SetHeader(headers.ContentType, "application/json")

	restClient := httpClient(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		return errors.Wrap(err, "failed to execute request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("wrong response status: %q", resp.Status)
	}

	return nil
}
