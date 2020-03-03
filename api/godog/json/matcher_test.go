package json

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMatcherSmoke(t *testing.T) {
	json := []byte(`{"key" : "value"}`)
	require.Nil(t, Match(json, ".key", "value"))
}

func TestMatcherNull(t *testing.T) {
	require.Nil(t, MatchNull([]byte(` null  `), ""))
	require.Nil(t, MatchNull([]byte(`null`), ""))
	require.Nil(t, MatchNull([]byte(`null`), "."))
	require.Nil(t, MatchNull([]byte(`{"key": null}`), ".key"))
	require.Nil(t, MatchNull([]byte(`{"key": [null]}`), ".key[0]"))

	require.NotNil(t, MatchNull([]byte(""), ""))
	require.NotNil(t, MatchNull([]byte(`{"key": []}`), ".key"))
	require.NotNil(t, MatchNull([]byte(`{"key": [null]}`), ".key"))
	require.NotNil(t, MatchNull([]byte(`{"key": {}}`), ".key"))
	require.NotNil(t, MatchNull([]byte(`{"key": "abc"}`), ".key"))
	require.NotNil(t, MatchNull([]byte(`null`), ".key"))
}

func TestMatcherNested(t *testing.T) {
	json := []byte(`{"keyA" : {"keyB" : "valueB"}}`)
	require.Nil(t, Match(json, ".keyA.keyB", "valueB"))
}

func TestMatcherError(t *testing.T) {
	json := []byte(`{"keyA" : {"keyB" : "valueA"}}`)
	require.Error(t, Match(json, ".keyA.keyB", "valueB"))
}

func TestMatcherNestedMulti(t *testing.T) {
	json := []byte(`{"keyA" : {"keyB" : "valueB", "keyC": "valueC"}}`)
	require.Nil(t, Match(json, ".keyA.keyC", "valueC"))
}

func TestMatcherNumbers(t *testing.T) {
	json := []byte(`{"key" : 12345}`)
	require.Nil(t, Match(json, ".key", `12345`))

	json = []byte(`{"key" : 123.456}`)
	require.Nil(t, Match(json, ".key", `\d{3}\.\d{3}`))

	json = []byte(`{"key" : -12345}`)
	require.Nil(t, Match(json, ".key", `-12345`))

	json = []byte(`{"key" : -123.456}`)
	require.Nil(t, Match(json, ".key", `-\d{3}\.\d{3}`))
}

func TestMatcherArray(t *testing.T) {
	json := []byte(`{"keyA" : ["value1", "value2"]}`)
	require.Nil(t, Match(json, ".keyA[1]", "value2"))
}

func TestRead(t *testing.T) {
	json := []byte(`{"key" : "value"}`)
	result, err := Read(json, ".key")
	require.Equal(t, "value", result)
	require.Nil(t, err)

	json = []byte(`{"key" : {"value" : 98} }`)
	_, err = Read(json, ".key")
	require.NotNil(t, err)

	json = []byte(`{"apa" : "value" }`)
	_, err = Read(json, ".key")
	require.NotNil(t, err)

}
