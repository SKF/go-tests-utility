package hierarchy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	hierarchyBaseURL = "https://api.%s.hierarchy.enlight.skf.com"
	companyType      = "company"
)

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
	}{
		ParentID:    parentNodeID,
		Label:       label,
		Description: description,
		Type:        nodetype,
		SubType:     subtype,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal failed")
		return
	}

	url := fmt.Sprintf(hierarchyBaseURL+"/nodes", stage)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		err = errors.Wrap(err, "http.NewRequest failed")
		return
	}
	req = req.WithContext(ctx)
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		if err = dd_tracer.Inject(span.Context(), dd_tracer.HTTPHeadersCarrier(req.Header)); err != nil {
			err = errors.Wrapf(err, "ddtracer.Inject: failed to inject trace headers")
			return
		}
	}

	req.Header.Set("Authorization", identityToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "client.Do failed")
		return
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "ioutil.ReadAll failed")
		return
	}

	var responseBody struct {
		ID string `json:"nodeId"`
	}

	if err = json.Unmarshal(body, &responseBody); err != nil {
		err = errors.Wrapf(err, "json.Unmarshal failed, body: %s", string(body))
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Wrong status: %q, body: %s", resp.Status, string(body))
		return
	}

	return responseBody.ID, nil
}

func Delete(identityToken, stage, nodeID string) error {
	return DeleteWithContext(context.Background(), identityToken, stage, nodeID)
}

func DeleteWithContext(ctx context.Context, identityToken, stage, nodeID string) error {
	url := fmt.Sprintf(hierarchyBaseURL+"/nodes/%s", stage, nodeID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
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

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return errors.Wrap(err, "ioutil.ReadAll failed")
		}

		return errors.Errorf("wrong status: %q, body: %s", resp.Status, string(body))
	}

	return nil
}
