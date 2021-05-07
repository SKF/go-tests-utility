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
	require.NotNil(t, MatchNull([]byte(`{"key": ""}`), ".key"))
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

func TestMatcherArrayRoot(t *testing.T) {
	require.Nil(t, Match([]byte(`["value1", "value2"]`), "[1]", "value2"))
	require.Nil(t, Match([]byte(`["value1", "value2"]`), ".[1]", "value2"))
}

func TestArrayLen(t *testing.T) {
	require.Nil(t, ArrayLen([]byte(`{"data" ["value1", "value2"]}`), ".data", 2))
	require.Nil(t, ArrayLen([]byte(`{"data" ["value1", "value2"]}`), ".data", 2))
	require.Nil(t, ArrayLen([]byte(`["value1", "value2"]`), ".", 2))
	require.Nil(t, ArrayLen([]byte(`["value1"]`), "", 1))
	require.Nil(t, ArrayLen([]byte(`["value1", 1]`), "", 2))

	require.NotNil(t, ArrayLen([]byte(`{"data" "value1"}`), ".data", -1))
	require.NotNil(t, ArrayLen([]byte(`{"data" null}`), ".data", -1))
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

func TestReadStringArr(t *testing.T) {
	result, err := ReadStringArr([]byte(`{"key" : ["value1", "value2"]}`), ".key")
	require.Nil(t, err)
	assert.Equal(t, []string{"value1", "value2"}, result)

	result, err = ReadStringArr([]byte(`["apa"]`), "")
	require.Nil(t, err)
	assert.Equal(t, []string{"apa"}, result)

	result, err = ReadStringArr([]byte(`["apa", ""]`), "")
	require.Nil(t, err)
	assert.Equal(t, []string{"apa", ""}, result)

	_, err = ReadStringArr([]byte(`{"key" : "value" }`), ".key")
	require.NotNil(t, err)

	_, err = ReadStringArr([]byte(`{"key" : {"value" : 98} }`), ".key")
	require.NotNil(t, err)

	_, err = ReadStringArr([]byte(`{"apa" : ["value1", "value2"] }`), ".key")
	require.NotNil(t, err)

	_, err = ReadStringArr([]byte(`["apa", {"a":1}]`), "")
	require.NotNil(t, err)
}

func TestKeyIsMissing(t *testing.T) {
	require.NoError(t, KeyIsMissing([]byte(`{}`), "a"))

	require.Error(t, KeyIsMissing([]byte(`{"a":1}`), "a"))
	require.Error(t, KeyIsMissing([]byte(`{"a":1}`), ".a"))

	require.Error(t, KeyIsMissing([]byte(`{"a": { "b": "test"} }`), ".a.b"))
	require.NoError(t, KeyIsMissing([]byte(`{"a": { "b": "test"} }`), "a.c"))

	require.NoError(t, KeyIsMissing([]byte(` null  `), ""))

}
