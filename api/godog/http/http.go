package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type HttpClient struct {
	token string
}

type HttpResponse struct {
	Status     string
	StatusCode int
	Headers    map[string][]string
	Body       []byte
}

func New(token string) *HttpClient {
	return &HttpClient{token: token}
}

func (c *HttpClient) Get(url string, out interface{}) (*HttpResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create GET request to endpoint: %s", url)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", c.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "GET request to endpoint: %s failed", url)
	}

	r, err := parseHttpResponse(resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse response from GET request to endpoint: %s", url)
	}

	if out != nil {
		if err = json.Unmarshal(r.Body, out); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal json response from GET request to endpoint: %s", url)
		}
	}

	return r, err
}

func (c *HttpClient) Post(url string, in interface{}, out interface{}) (*HttpResponse, error) {
	client := &http.Client{}

	bs := new(bytes.Buffer)
	if err := json.NewEncoder(bs).Encode(in); err != nil {
		return nil, errors.Wrapf(err, "Failed marshal body for POST request to endpoint: %s", url)
	}

	req, err := http.NewRequest("POST", url, bs)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create POST request to endpoint: %s", url)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", c.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "POST request to endpoint: %s failed", url)
	}

	r, err := parseHttpResponse(resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse response from POST request to endpoint: %s", url)
	}

	if out != nil {
		if err = json.Unmarshal(r.Body, out); err == nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal json response from POST request to endpoint: %s", url)
		}
	}

	return r, err
}

func (c *HttpClient) Put(url string, in interface{}, out interface{}) (*HttpResponse, error) {
	client := &http.Client{}

	bs := new(bytes.Buffer)
	if err := json.NewEncoder(bs).Encode(in); err != nil {
		return nil, errors.Wrapf(err, "Failed marshal body for POST request to endpoint: %s", url)
	}

	req, err := http.NewRequest("PUT", url, bs)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create POST request to endpoint: %s", url)
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("authorization", c.token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "POST request to endpoint: %s failed", url)
	}

	r, err := parseHttpResponse(resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse response from POST request to endpoint: %s", url)
	}

	if out != nil {
		if err = json.Unmarshal(r.Body, out); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal json response from POST request to endpoint: %s", url)
		}
	}

	return r, err
}

func parseHttpResponse(resp *http.Response) (*HttpResponse, error) {
	headers := make(map[string][]string)

	for h, v := range resp.Header {
		headers[h] = v
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &HttpResponse{
		Status:     resp.Status,
		StatusCode: resp.StatusCode,
		Headers:    headers,
		Body:       body,
	}, nil
}
