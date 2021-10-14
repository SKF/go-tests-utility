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

func Try(function func() (bool, error), timeout time.Duration) error {
	nrOfRetries := 0
	timeSlept := time.Duration(0)

	for {
		success, err := function()
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		sleepDuration := time.Duration(math.Pow(powBase, float64(nrOfRetries))*startingMillisToWait) * time.Millisecond

		timeSlept += sleepDuration
		if timeSlept > timeout {
			break
		}

		time.Sleep(sleepDuration)
		nrOfRetries++
	}

	return fmt.Errorf("timeout <%s> reached before condition was met, #retries = %d", timeout, nrOfRetries)
}
