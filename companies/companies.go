package companies

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

const hierarchyBaseURL = "https://api.%s.hierarchy.enlight.skf.com"

func Create(identityToken, stage, parentNodeID, label, description string) (_ string, err error) {
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
		Type:        "company",
		SubType:     "company",
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

func Delete(identityToken, stage, companyID string) error {
	url := fmt.Sprintf(hierarchyBaseURL+"/nodes/%s", stage, companyID)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
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
