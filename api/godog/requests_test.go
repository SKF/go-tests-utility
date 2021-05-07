package godog

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseFeature_SetRequestBodyStringListParameterTo(t *testing.T) {
	api := BaseFeature{}
	api.GetValue = func(key string) (value string, err error) {
		return key + "fixed", nil
	}

	err := api.CreatePathRequest(http.MethodPost, "")
	require.NoError(t, err)

	err = api.SetRequestBodyStringListParameterTo("names", "apa, .bepa, cepa")
	require.NoError(t, err)

	jsonBody, err := json.Marshal(api.Request.Body)
	require.NoError(t, err)

	result := struct {
		Names []string `json:"names"`
	}{}
	err = json.Unmarshal(jsonBody, &result)
	require.NoError(t, err)

	require.Equal(t, []string{"apa", ".bepafixed", "cepa"}, result.Names)
}
