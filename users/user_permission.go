package users

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	dd_tracer "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func AddUserRole(identityToken, stage, userID, role string) error {
	return AddUserRoleWithContext(context.Background(), identityToken, stage, userID, role)
}

func AddUserRoleWithContext(ctx context.Context, identityToken, stage, userID, role string) (err error) {
	user, err := getUser(ctx, identityToken, stage, userID)
	if err != nil {
		return
	}

	user.UserRoles = append(user.UserRoles, role)
	return updateUser(ctx, identityToken, stage, user)
}

func getUser(ctx context.Context, identityToken, stage, userID string) (user user, err error) {
	if userID == "" {
		return user, fmt.Errorf("userID is required")
	}

	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s", stage, userID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
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

func updateUser(ctx context.Context, identityToken, stage string, user user) (err error) {
	url := fmt.Sprintf(accessMgmtBaseURL+"/users/%s", stage, user.ID)

	body, err := json.Marshal(user)
	if err != nil {
		return errors.Wrap(err, "json.Marshal failed")
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(body))
	if err != nil {
		return errors.Wrap(err, "http.NewRequest failed")
	}
	req = req.WithContext(ctx)
	if span, ok := dd_tracer.SpanFromContext(ctx); ok {
		if err = dd_tracer.Inject(span.Context(), dd_tracer.HTTPHeadersCarrier(req.Header)); err != nil {
			return errors.Wrapf(err, "ddtracer.Inject: failed to inject trace headers")
		}
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
