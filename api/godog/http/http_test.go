package http

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const token = "let me in!"

func TestSmokeGet(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, "GET", req.Method)
		require.Equal(t, token, req.Header.Get("Authorization"))
		rw.Header().Set("Content-Type", "application/json; version=1.0")
		fmt.Fprintln(rw, `{"key":"this is the value"}`)
	}))
	defer ts.Close()

	client := NewWithToken(token)
	resp, err := client.Get(ts.URL, nil)
	require.Nil(t, err)
	require.Equal(t, "200 OK", resp.Status)
	require.Equal(t, "application/json; version=1.0", resp.Headers["Content-Type"][0])
	require.Equal(t, []byte("{\"key\":\"this is the value\"}\n"), resp.Body)
}

func TestSmokePost(t *testing.T) {
	realBody := struct {
		Key string `json:"key"`
	}{"this is the value"}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, "POST", req.Method)
		require.Equal(t, token, req.Header.Get("Authorization"))
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		body, err := io.ReadAll(req.Body)
		req.Body.Close()

		require.Nil(t, err)
		require.Equal(t, []byte("{\"key\":\"this is the value\"}\n"), body)
	}))
	defer ts.Close()

	client := NewWithToken(token)
	resp, err := client.Post(ts.URL, realBody, nil)
	require.Nil(t, err)
	require.Equal(t, "200 OK", resp.Status)
}

func TestSmokePut(t *testing.T) {
	realBody := struct {
		Key string `json:"key"`
	}{"this is the value"}

	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, "PUT", req.Method)
		require.Equal(t, token, req.Header.Get("Authorization"))
		require.Equal(t, "application/json", req.Header.Get("Content-Type"))

		body, err := io.ReadAll(req.Body)
		req.Body.Close()

		require.Nil(t, err)
		require.Equal(t, []byte("{\"key\":\"this is the value\"}\n"), body)
	}))
	defer ts.Close()

	client := NewWithToken(token)
	resp, err := client.Put(ts.URL, realBody, nil)
	require.Nil(t, err)
	require.Equal(t, "200 OK", resp.Status)
}
