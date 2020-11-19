package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

func AddUserRole(identityToken, stage, userID, role string) (err error) {
	user, err := getUser(identityToken, stage, userID)
	if err != nil {
		return
	}

	user.UserRoles = append(user.UserRoles, role)
	return updateUser(identityToken, stage, user)
}

func getUser(identityToken, stage, userID string) (user user, err error) {
	fmt.Printf("userID: %v\n", userID)
	if userID == "" {
		return user, fmt.Errorf("userID is required")
	}

	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s", stage, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		err = errors.Wrap(err, "http.NewRequest failed")
		return
	}

	req.Header.Set("Authorization", identityToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "client.Do failed")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Wrong status: %q", resp.Status)
		return
	}

	defer resp.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(resp.Body); err != nil {
		err = errors.New("couldn't read response body")
		return
	}

	if err = json.Unmarshal(body, &user); err != nil {
		err = errors.Errorf("couldn't decode response body, %s", string(body))
		return
	}

	return user, err
}

func updateUser(identityToken, stage string, user user) (err error) {
	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s", stage, user.ID)

	body, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "json.Marshal failed")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
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

type user struct {
	ID             string   `json:"id"`
	Email          string   `json:"email"`
	CompanyID      string   `json:"companyId"`
	UserRoles      []string `json:"userRoles"`
	Username       string   `json:"username"`
	UserStatus     string   `json:"userStatus"`
	EulaAgreedDate string   `json:"eulaAgreedDate"`
	ValidEula      bool     `json:"validEula"`
	Firstname      string   `json:"firstname"`
	Surname        string   `json:"surname"`
	Locale         string   `json:"locale"`
	CreatedDate    string   `json:"createdDate"`
	UpdatedDate    string   `json:"updatedDate"`
}
