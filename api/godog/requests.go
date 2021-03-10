package godog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"

	json_matcher "github.com/SKF/go-tests-utility/api/godog/json"
	http_model "github.com/SKF/go-utility/http-model"
	"github.com/SKF/go-utility/log"
)

func (api *BaseFeature) CreatePathRequest(method, path string) error {
	api.Request = Request{
		Headers: make(http.Header),
		Body:    make(map[string]interface{}),
		Method:  method,
		Url:     api.baseURL + path,
	}
	return nil
}

func (api *BaseFeature) SetRequestHeaderParameterTo(key, value string) (err error) {
	if strings.HasPrefix(value, ".") {
		if value, err = api.GetValue(value); err != nil {
			return err
		}
	}

	api.Request.Headers.Add(key, value)
	return
}

func (api *BaseFeature) SetsRequestPathParameterTo(key, value string) (err error) {
	if strings.HasPrefix(value, ".") {
		if value, err = api.GetValue(value); err != nil {
			return err
		}
	}

	keyPattern := "{" + key + "}"
	if !strings.Contains(api.Request.Url, keyPattern) {
		return errors.New("api.Request path does not contain variable: " + keyPattern)
	}

	api.Request.Url = strings.ReplaceAll(api.Request.Url, keyPattern, value)
	return
}

func (api *BaseFeature) SetRequestBodyParameterTo(key, value string) (err error) {
	if strings.HasPrefix(value, ".") {
		if value, err = api.GetValue(value); err != nil {
			return err
		}
	}

	keyParts := strings.Split(key, ".")
	prevMap := api.Request.Body
	for idx, key := range keyParts {
		if len(keyParts) == idx+1 {
			prevMap[key] = value
			break
		}

		if _, exists := prevMap[key]; !exists {
			prevMap[key] = make(map[string]interface{})
		}
		prevMap = prevMap[key].(map[string]interface{})
	}

	return
}

func (api *BaseFeature) SetRequestBodyStringListParameterTo(key, valuesstr string) (err error) {
	values := strings.Split(valuesstr, ",")
	list := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if strings.HasPrefix(value, ".") {
			if value, err = api.GetValue(value); err != nil {
				return err
			}
		}

		if len(value) > 0 {
			list = append(list, value)
		}
	}

	keyParts := strings.Split(key, ".")
	prevMap := api.Request.Body
	for idx, key := range keyParts {
		if len(keyParts) == idx+1 {
			prevMap[key] = list
			break
		}

		if _, exists := prevMap[key]; !exists {
			prevMap[key] = make(map[string]interface{})
		}
		prevMap = prevMap[key].(map[string]interface{})
	}

	return
}

func (api *BaseFeature) ExecuteTheRequest() error {
	return api.ExecuteTheRequestWithContext(context.Background())
}

func (api *BaseFeature) ExecuteTheRequestWithContext(ctx context.Context) (err error) {
	jsonBody, err := json.Marshal(api.Request.Body)
	if err != nil {
		return errors.Wrap(err, "json.Marshal failed")
	}

	if api.Request.Method == http.MethodGet {
		jsonBody = nil
	}

	return api.ExecuteTheRequestWithPayloadAndContext(ctx, jsonBody)
}

func (api *BaseFeature) ExecuteTheRequestWithPayload(payload []byte) error {
	return api.ExecuteTheRequestWithPayloadAndContext(context.Background(), payload)
}

func (api *BaseFeature) ExecuteTheRequestWithPayloadAndContext(ctx context.Context, payload []byte) (err error) {
	log.Debugf("Request %s: %s\n", api.Request.Method, payload)
	log.Debugf("req headers: %v\n", api.Request.Headers)

	if len(payload) > 0 {
		log.Debugf("Payload: %s\n", payload)
	}

	var bodyBuffer io.Reader
	if payload != nil {
		bodyBuffer = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequestWithContext(ctx, api.Request.Method, api.Request.Url, bodyBuffer)
	if err != nil {
		return errors.Wrapf(err, "http.NewRequest failed - Payload: `%s`", string(payload))
	}
	req.Header = api.Request.Headers

	api.Request.ExecutionTime = time.Now()
	client := dd_http.WrapClient(
		&http.Client{},
		dd_http.RTWithResourceNamer(func(req *http.Request) string {
			return fmt.Sprintf("%s %s", req.Method, req.URL.String())
		}),
	)

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "client.Do failed - header: `%+v`", req.Header)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadAll failed")
	}

	api.Response.Raw = resp
	api.Response.Body = body

	log.Debugf("Response: %s", body)
	return nil
}

func (api *BaseFeature) ExecuteInvalidRequest() error {
	return api.ExecuteInvalidRequestWithContext(context.Background())
}

func (api *BaseFeature) ExecuteInvalidRequestWithContext(ctx context.Context) error {
	invalidBody := []byte(`{ "param": "value",}`)
	return api.ExecuteTheRequestWithPayloadAndContext(ctx, invalidBody)
}

func (api *BaseFeature) AssertNotEmpty(responseKey string) error {
	value, err := json_matcher.Read(api.Response.Body, responseKey)
	if err != nil {
		return err
	}

	if value == "" {
		return errors.New(fmt.Sprintf("No value found for: %v", responseKey))
	}
	return nil
}

func (api *BaseFeature) AssertResponseBodyValueEquals(key, expected string) (err error) {
	switch key {
	case "len(.data)":
		return api.AssertDataLength(expected)
	}

	if expected, err = api.GetValue(expected); err != nil {
		return
	}

	actual, err := json_matcher.Read(api.Response.Body, key)
	if err != nil {
		return
	}

	if actual != expected {
		return errors.Errorf("Match error: Values mismatch, expected: '%s' actual: '%s'", expected, actual)
	}

	return
}

func (api *BaseFeature) AssertDataLength(expected string) error {
	expectedLen, err := strconv.Atoi(expected)
	if err != nil {
		return err
	}

	return json_matcher.ArrayLen(api.Response.Body, ".data", expectedLen)
}

func (api *BaseFeature) AssertResponseCode(code int) (err error) {
	if api.Response.Raw.StatusCode != code {
		err = errors.Errorf("expected status code: %d, got: %d \n response: %s, request: %+v", code, api.Response.Raw.StatusCode, string(api.Response.Body), api.Request)
		return
	}

	return
}

func (api *BaseFeature) AssertErrorIs(errorMessage string, code int) (err error) {
	if err = api.AssertResponseCode(code); err != nil {
		return
	}
	if err = api.AssertResponseBodyErrorMessageIs(errorMessage); err != nil {
		return
	}
	return
}

func (api *BaseFeature) AssertResponseBodyErrorMessageIs(errorMessage string) (err error) {
	var responseBody http_model.ErrorResponse
	if err = json.Unmarshal(api.Response.Body, &responseBody); err != nil {
		return
	}
	if responseBody.Error.Message != errorMessage {
		err = errors.Errorf("expected error message: %s, got: %s", errorMessage, responseBody.Error.Message)
		return
	}
	return
}
