package godog

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	json_matcher "github.com/SKF/go-tests-utility/api/godog/json"
)

func TestGetRequest(t *testing.T) {
	api := BaseFeature{}
	api.SetBaseUrl("https://jsonplaceholder.typicode.com")

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
	api := BaseFeature{}
	api.SetBaseUrl("https://jsonplaceholder.typicode.com")

	err := api.CreatePathRequest(http.MethodPost, "/posts")
	require.NoError(t, err)

	err = api.ExecuteInvalidRequest()
	require.NoError(t, err)

	err = api.AssertResponseCode(http.StatusInternalServerError)
	require.NoError(t, err)
}
