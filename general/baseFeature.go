package general

import (
	"net/http"
	"time"
)

type BaseFeature struct {
	StartedAt         time.Time
	Response          response
	Request           Request
	baseURL           string
	deprecatedBaseURL string

	GetValue func(key string) (value string, err error)
}

func GetBaseFeature(baseUrl, deprecatedBaseURL string) BaseFeature {
	return BaseFeature{
		baseURL: baseUrl,
		deprecatedBaseURL:deprecatedBaseURL,
	}
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
