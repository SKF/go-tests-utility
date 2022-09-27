package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
	dd_http "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
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

func New() *HttpClient {
	return &HttpClient{}
}

func NewWithToken(token string) *HttpClient {
	return &HttpClient{token: token}
}

func (c *HttpClient) FetchToken(stage, username, password string) error {
	return c.FetchTokenWithContext(context.Background(), stage, username, password)
}

func (c *HttpClient) FetchTokenWithContext(ctx context.Context, stage, username, password string) error {
	in := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{username, password}

	out := &struct {
		Token string `json:"token"`
	}{}

	var url string
	if stage == "prod" {
		url = "https://api-auth.users.enlight.skf.com/login"
	} else {
		url = fmt.Sprintf("https://api-auth.%s.users.enlight.skf.com/login", stage)
	}

	bs := new(bytes.Buffer)
	if err := json.NewEncoder(bs).Encode(in); err != nil {
		return errors.Wrapf(err, "Failed marshal body for POST request to endpoint: %s", url)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bs)
	if err != nil {
		return errors.Wrapf(err, "Failed to create POST request to endpoint: %s", url)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")

	client := dd_http.WrapClient(
		&http.Client{},
		dd_http.RTWithResourceNamer(func(req *http.Request) string {
			return fmt.Sprintf("%s %s", req.Method, req.URL.String())
		}),
	)

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "POST request to endpoint: %s failed", url)
	}

	defer resp.Body.Close()
	bodybytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(bodybytes, out); err != nil {
		return errors.Wrapf(err, "Failed to unmarshal json response from POST request to endpoint: %s, body: '%s'", url, bodybytes)
	}

	c.token = out.Token
	return nil
}

func (c *HttpClient) Get(url string, out interface{}) (*HttpResponse, error) {
	return c.GetWithContext(context.Background(), url, out)
}

func (c *HttpClient) GetWithContext(ctx context.Context, url string, out interface{}) (*HttpResponse, error) {
	return c.send(ctx, "GET", url, nil, out)
}

func (c *HttpClient) Post(url string, in interface{}, out interface{}) (*HttpResponse, error) {
	return c.PostWithContext(context.Background(), url, in, out)
}

func (c *HttpClient) PostWithContext(ctx context.Context, url string, in interface{}, out interface{}) (*HttpResponse, error) {
	return c.send(ctx, "POST", url, in, out)
}

func (c *HttpClient) Put(url string, in interface{}, out interface{}) (*HttpResponse, error) {
	return c.PutWithContext(context.Background(), url, in, out)
}

func (c *HttpClient) PutWithContext(ctx context.Context, url string, in interface{}, out interface{}) (*HttpResponse, error) {
	return c.send(ctx, "PUT", url, in, out)
}

func (c *HttpClient) Delete(url string, out interface{}) (*HttpResponse, error) {
	return c.DeleteWithContext(context.Background(), url, out)
}

func (c *HttpClient) DeleteWithContext(ctx context.Context, url string, out interface{}) (*HttpResponse, error) {
	return c.send(ctx, "DELETE", url, nil, out)
}

func (c *HttpClient) send(ctx context.Context, method, url string, in interface{}, out interface{}) (*HttpResponse, error) {
	bs := new(bytes.Buffer)
	sendBody := in != nil && (method == "POST" || method == "PUT")

	if sendBody {
		if err := json.NewEncoder(bs).Encode(in); err != nil {
			return nil, errors.Wrapf(err, "Failed marshal body for %s request to endpoint: %s", method, url)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bs)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to create %s request to endpoint: %s", method, url)
	}
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", c.token)

	if sendBody {
		req.Header.Set("content-type", "application/json")
	}

	client := dd_http.WrapClient(
		&http.Client{},
		dd_http.RTWithResourceNamer(func(req *http.Request) string {
			return fmt.Sprintf("%s %s", req.Method, req.URL.String())
		}),
	)
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrapf(err, "%s request to endpoint: %s failed", method, url)
	}

	r, err := parseHttpResponse(resp)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse response from %s request to endpoint: %s", method, url)
	}

	if out != nil {
		if err = json.Unmarshal(r.Body, out); err != nil {
			return nil, errors.Wrapf(err, "Failed to unmarshal json response from %s request to endpoint: %s, Body: '%s'", method, url, r.Body)
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
	body, err := io.ReadAll(resp.Body)
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
