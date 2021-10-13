package retry

import (
	"fmt"
	"math"
	"time"
)

type Until struct {
	Condition func(body []byte) bool
	Timeout   time.Duration
}

func Try(function func() (bool, error), timeout time.Duration) error {
	for i := 0; i < int(timeout.Milliseconds()); i++ {
		success, err := function()
		if err != nil {
			return err
		}

		if success {
			return nil
		}

		sleepDuration := time.Duration(math.Pow(2, float64(i))*100) * time.Millisecond // nolint: gomnd
		if sleepDuration > timeout {
			break
		}

		time.Sleep(sleepDuration)
	}

	return fmt.Errorf("timeout reached before condition was met")
}
