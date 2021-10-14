package retry_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/SKF/go-tests-utility/api/godog/retry"
)

func Test_ReachTimeout(t *testing.T) {
	// When
	err := retry.Try(func() (bool, error) {
		return false, nil
	}, 1*time.Millisecond)

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout <1ms> reached")
}

func Test_Retries(t *testing.T) {
	tests := []struct {
		name              string
		nrOfRetriesNeeded int
		retryTimeout      time.Duration
		shouldFail        bool
	}{
		{
			name:              "success after timeout window",
			nrOfRetriesNeeded: 2,
			retryTimeout:      200 * time.Millisecond,
			shouldFail:        true,
		},
		{
			name:              "success within timeout window",
			nrOfRetriesNeeded: 1,
			retryTimeout:      200 * time.Millisecond,
			shouldFail:        false,
		},
		{
			name:              "exponential backoff, success after timeout window",
			nrOfRetriesNeeded: 4,
			retryTimeout:      1450 * time.Millisecond, // sleep for 700 ms (next sleep will be 800 ms), total 1500
			shouldFail:        true,
		},
		{
			name:              "exponential backoff, success within timeout window",
			nrOfRetriesNeeded: 3,
			retryTimeout:      1450 * time.Millisecond, // sleep for 700 ms (next sleep will be 800 ms), total 1500
			shouldFail:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			nrOfRetries := 0
			err := retry.Try(func() (bool, error) {
				if nrOfRetries == tt.nrOfRetriesNeeded {
					return true, nil
				}
				nrOfRetries++

				return false, nil
			}, tt.retryTimeout)

			// Then
			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
