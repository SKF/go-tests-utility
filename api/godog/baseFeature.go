package godog

import (
	"fmt"
	"net/http"
	"time"
)

type BaseFeature struct {
	StartedAt time.Time
	Response  response
	Request   Request
	baseURL   string

	GetValue func(key string) (value string, err error)
}

func (api *BaseFeature) SetBaseUrl(baseUrl string) {
	api.baseURL = baseUrl
}

type Request struct {
	Url           string
	Body          map[string]interface{}
	Headers       http.Header
	Method        string
	ExecutionTime time.Time
}

func (r * Request) String() string{
	return fmt.Sprintf("%s :: %s headers: %s time: %v\n %s", r.Method, r.Url, r.Headers, r.ExecutionTime, r.Body)
}

type response struct {
	Body []byte
	Raw  *http.Response
}
