package godog_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/SKF/go-tests-utility/api/godog"
)

func TestBaseFeature_assertEquals(t *testing.T) {
	state := godog.BaseFeature{}
	state.GetValue = func(key string) (value string, err error) {
		switch key {
		case "apa":
			return "apa", nil
		case "apa2":
			return "apa", nil
		case "bepa":
			return "bepa", nil
		default:
			return "", fmt.Errorf("key: %s not found", key)
		}
	}

	assert.NoError(t, state.AssertEquals("apa", "apa2"))

	assert.Error(t, state.AssertEquals("apa", "bepa"))
	assert.Error(t, state.AssertEquals("apa", "cepa"))
}
