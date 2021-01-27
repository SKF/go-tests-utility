package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	disposable_emails "github.com/SKF/go-tests-utility/disposable-emails"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const identityMgmtBaseURL = "https://sso-api.%s.users.enlight.skf.com"

func Create(accessToken, stage, companyID, email string) (User, string, error) {
	return CreateWithContext(context.Background(), accessToken, stage, companyID, email)
}

func CreateWithContext(ctx context.Context, accessToken, stage, companyID, email string) (_ User, password string, err error) {
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
	req = req.WithContext(ctx)
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		if err = dd_tracer.Inject(span.Context(), dd_tracer.HTTPHeadersCarrier(req.Header)); err != nil {
			err = errors.Wrapf(err, "ddtracer.Inject: failed to inject trace headers")
			return
		}
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

	temporaryPassword, err := PollForTemporaryPassword(email, startedAt)
	if err != nil {
		return
	}

	return respBody.Data, temporaryPassword, nil
}

func PollForTemporaryPassword(email string, startedAt time.Time) (string, error) {
	const subject = "Welcome to SKF Digital Services"
	actualEmail, err := disposable_emails.PollForMessageWithSubject(email, subject, startedAt)
	if err != nil {
		return "", err
	}

	temporaryPassword, err := getTemporaryPassword(actualEmail)
	if err != nil {
		return "", err
	}

	return temporaryPassword, nil
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
