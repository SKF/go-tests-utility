package users

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	disposable_emails "github.com/SKF/go-tests-utility/disposable-emails"
)

const identityMgmtBaseURL = "https://sso-api.%s.users.enlight.skf.com"

func Create(accessToken, stage, companyID, email string) (_ User, password string, err error) {
	startedAt := time.Now().Add(-1 * time.Second)

	requestBody := struct {
		Email     string `json:"email"`
		GivenName string `json:"givenName"`
		Surname   string `json:"surname"`
	}{
		Email:     email,
		GivenName: "Foo",
		Surname:   "Bar",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		err = errors.Wrap(err, "json.Marshal failed")
		return
	}

	url := fmt.Sprintf(identityMgmtBaseURL+"/companies/%s/users", stage, companyID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		err = errors.Wrap(err, "http.NewRequest failed")
		return
	}

	req.Header.Set("Authorization", accessToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "client.Do failed")
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.Wrap(err, "ioutil.ReadAll failed")
		return
	}

	var respBody struct {
		Data User `json:"data"`
	}

	if err = json.Unmarshal(body, &respBody); err != nil {
		err = errors.Wrapf(err, "json.Unmarshal failed, body: %s", string(body))
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("Wrong status: %q, req: %+v", resp.Status, req)
		return
	}

	const subject = "Welcome to SKF Enlight Centre"
	actualEmail, err := disposable_emails.PollForMessageWithSubject(email, subject, startedAt)
	if err != nil {
		return
	}

	temporaryPassword, err := getTemporaryPassword(actualEmail)
	if err != nil {
		return
	}

	return respBody.Data, temporaryPassword, nil
}

func getTemporaryPassword(emailMessage string) (string, error) {
	subMatches := temporaryPasswordRegexp.FindAllStringSubmatch(emailMessage, -1)
	if len(subMatches) == 0 {
		return "", errors.Errorf("couldn't retrieve temporary password from email: [%s]", emailMessage)
	}
	return subMatches[0][1], nil
}

var temporaryPasswordRegexp = regexp.MustCompile(`password=(\S+)\" style`)

type User struct {
	ID        string `json:"id"`
	CompanyID string `json:"companyId"`
	Email     string `json:"email"`
	GivenName string `json:"givenName"`
	Surname   string `json:"surname"`
	Language  string `json:"language"`
	Status    string `json:"status"`
}
