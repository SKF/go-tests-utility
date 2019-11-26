package users

import (
	"fmt"
	"net/http"

	"github.com/SKF/go-utility/uuid"
	"github.com/pkg/errors"
)

const accessMgmtBaseURL = "https://api-web.%s.users.enlight.skf.com"

func AddUserAccess(identityToken, stage, userID, companyID string) (err error) {
	if !uuid.IsValid(userID) {
		return errors.Errorf("Invalid User ID: %q", userID)
	}

	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s/hierarchies/%s", stage, userID, companyID)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
	}

	req.Header.Set("Authorization", identityToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "client.Do failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("Wrong status: %q", resp.Status)
	}

	return err
}
