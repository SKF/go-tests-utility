package retry

import (
	"fmt"
	"math"
	"time"
)

const (
	startingMillisToWait = 100
	powBase              = 2
)

type Until struct {
	Condition func(body []byte) bool
	Timeout   time.Duration
}

type UntilWithError struct {
	Condition func(body []byte) (bool, error)
	Timeout   time.Duration
}

func Try(function func() (bool, error), timeout time.Duration) error {
	nrOfRetries := 0
	endBefore := time.Now().Add(timeout)

	for {
		success, err := function()
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		sleepDuration := time.Duration(math.Pow(powBase, float64(nrOfRetries))*startingMillisToWait) * time.Millisecond
		nextRuntime := time.Now().Add(sleepDuration)

		if nextRuntime.After(endBefore) {
			break
		}

		time.Sleep(sleepDuration)
		nrOfRetries++
	}

	return fmt.Errorf("timeout <%s> reached before condition was met, #retries = %d", timeout, nrOfRetries)
}
