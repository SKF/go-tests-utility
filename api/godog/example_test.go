package godog

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	json_matcher "github.com/SKF/go-tests-utility/api/godog/json"
)

func TestGetRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/todos/1", r.URL.Path)

		body := `
{
	"userId": 1,
	"id": 1,
	"title": "delectus aut autem",
	"completed": false
}`

		fmt.Fprintln(w, body)
	}))
	defer s.Close()

	api := BaseFeature{}
	api.SetBaseUrl(s.URL)

	err := api.CreatePathRequest(http.MethodGet, "/todos/{id}")
	require.NoError(t, err)

	id := "1"
	err = api.SetsRequestPathParameterTo("id", id)
	require.NoError(t, err)

	err = api.ExecuteTheRequest()
	require.NoError(t, err)

	err = api.AssertResponseCode(http.StatusOK)
	assert.NoError(t, err)

	res, err := json_matcher.Read(api.Response.Body, ".id")
	require.NoError(t, err)
	assert.Equal(t, id, res)
}

func TestCreateInvalidRequest(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/posts", r.URL.Path)

		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer s.Close()

	api := BaseFeature{}
	api.SetBaseUrl(s.URL)

	err := api.CreatePathRequest(http.MethodPost, "/posts")
	require.NoError(t, err)
	err = api.SetRequestHeaderParameterTo("Content-Type", "application/json")
	require.NoError(t, err)

	err = api.ExecuteInvalidRequest()
	require.NoError(t, err)

	err = api.AssertResponseCode(http.StatusInternalServerError)
	require.NoError(t, err)
}
