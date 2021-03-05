package godog

import (
	"github.com/pkg/errors"
)

func (api *BaseFeature) AssertEquals(actual, expected string) (err error) {
	if actual, err = api.GetValue(actual); err != nil {
		return
	}

	if expected, err = api.GetValue(expected); err != nil {
		return
	}

	if actual != expected {
		return errors.Errorf("match error: Values mismatch, expected: '%s' actual: '%s'", expected, actual)
	}

	return
}
