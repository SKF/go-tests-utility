package users

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

func Delete(accessToken, stage, userID string) error {
	url := fmt.Sprintf(identityMgmtBaseURL+"/users/%s", stage, userID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return errors.Errorf("wrong status: %q", resp.Status)
	}

	return nil
}
