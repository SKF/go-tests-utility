package godog

import (
	"context"
	"net/http"
	"time"
)

type BaseFeature struct {
	StartedAt time.Time
	Response  response
	Request   Request
	baseURL   string
	ctx       context.Context

	GetValue func(key string) (value string, err error)
}

func (api *BaseFeature) SetContext(ctx context.Context) {
	api.ctx = ctx
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

type response struct {
	Body []byte
	Raw  *http.Response
}
