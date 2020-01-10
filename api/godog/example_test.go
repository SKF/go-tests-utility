package godog

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	json_matcher "github.com/SKF/go-tests-utility/api/godog/json"
)

func TestGetRequest(t *testing.T) {
	api := BaseFeature{}
	api.SetBaseUrl("http://dummy.restapiexample.com/")

	err := api.CreatePathRequest(http.MethodGet, "/api/v1/employee/{id}")
	assert.NoError(t, err)

	id := "1"
	err = api.SetsRequestPathParameterTo("id", id)
	assert.NoError(t, err)


	err = api.ExecuteTheRequest()
	assert.NoError(t, err)

	err = api.AssertResponseCode(http.StatusOK)
	assert.NoError(t, err)

	res, err := json_matcher.Read(api.Response.Body, ".id")
	assert.NoError(t, err)
	assert.Equal(t, id, res)
}
