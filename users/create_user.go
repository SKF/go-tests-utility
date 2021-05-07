package users

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/SKF/go-rest-utility/client"
	"github.com/SKF/go-rest-utility/client/auth"

	disposable_emails "github.com/SKF/go-tests-utility/disposable-emails"
)

const identityMgmtBaseURL = "https://sso-api.%s.users.enlight.skf.com"

func httpClientIdentityMgmt(stage, identityToken string) *client.Client {
	return client.NewClient(
		client.WithBaseURL(fmt.Sprintf(identityMgmtBaseURL, stage)),
		client.WithDatadogTracing(),
		client.WithTokenProvider(auth.RawToken(identityToken)),
	)
}

func Create(identityToken, stage, companyID, email string) (User, string, error) {
	return CreateWithContext(context.Background(), identityToken, stage, companyID, email)
}

func CreateWithContext(ctx context.Context, identityToken, stage, companyID, email string) (_ User, password string, err error) {
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

	req := client.Post("/companies/{companyId}/users").
		Assign("companyId", companyID).
		WithJSONPayload(requestBody)

	restClient := httpClientIdentityMgmt(stage, identityToken)
	resp, err := restClient.Do(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "failed to execute request")
		return
	}

	var respBody struct {
		Data User `json:"data"`
	}

	if err = resp.Unmarshal(&respBody); err != nil {
		err = errors.Wrap(err, "failed to unmarshal response")
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.Errorf("wrong response status: %q", resp.Status)
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
