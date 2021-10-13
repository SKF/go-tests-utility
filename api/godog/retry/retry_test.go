package retry

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_ReachTimeout(t *testing.T) {
	// When
	err := Try(func() (bool, error) {
		return false, nil
	}, 1*time.Millisecond)

	// Then
	require.Error(t, err)
	require.Contains(t, err.Error(), "timeout reached")
}

func Test_Retries(t *testing.T) {
	tests := []struct {
		name              string
		nrOfRetriesNeeded int
		nrOfTimesToRetry  int
		shouldFail        bool
	}{
		{
			name:              "timeout before success",
			nrOfRetriesNeeded: 2,
			nrOfTimesToRetry:  1,
			shouldFail:        true,
		},
		{
			name:              "success within timeout window",
			nrOfRetriesNeeded: 1,
			nrOfTimesToRetry:  2,
			shouldFail:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When
			nrOfRetries := 0
			err := Try(func() (bool, error) {
				if nrOfRetries == tt.nrOfRetriesNeeded {
					return true, nil
				}
				nrOfRetries++

				return false, nil
			}, time.Duration(tt.nrOfTimesToRetry*100)*time.Millisecond)

			// Then
			if tt.shouldFail {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
