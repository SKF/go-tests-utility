package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

const ssoBaseURL = "https://sso-api.%s.users.enlight.skf.com"

// SignIn will sign in the user and if needed complete the change password challenge
func SignIn(stage, username, password string) (tokens Tokens, err error) {
	var resp SignInResponse

	if resp, err = initiateSignIn(stage, username, password); err != nil {
		return
	}

	if resp.Data.Challenge.Type == "" {
		tokens = resp.Data.Tokens
		return
	}

	if resp, err = completeSignIn(stage, resp.Data.Challenge, username, password); err != nil {
		return
	}

	tokens = resp.Data.Tokens
	return
}

func initiateSignIn(stage, username, password string) (signInResp SignInResponse, err error) {
	url := fmt.Sprintf(ssoBaseURL+"/sign-in/initiate", stage)

	jsonBody := `{"username": "` + username + `", "password": "` + password + `"}`
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return
		}

		err = errors.Errorf("StatusCode: %s, Error Message: %s \n", resp.Status, errorResp.Error.Message)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&signInResp); err != nil {
		return
	}

	return signInResp, err
}

func completeSignIn(stage string, challenge Challenge, username, newPassword string) (signInResp SignInResponse, err error) {
	url := fmt.Sprintf(ssoBaseURL+"/sign-in/complete", stage)

	baseJSON := `{"username": "%s", "id": "%s", "type": "%s", "properties": {"newPassword": "%s"}}`
	jsonBody := fmt.Sprintf(baseJSON, username, challenge.ID, challenge.Type, newPassword)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(jsonBody))
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}

		if err = json.NewDecoder(resp.Body).Decode(&errorResp); err != nil {
			return
		}

		err = errors.Errorf("StatusCode: %s, Error Message: %s \n", resp.Status, errorResp.Error.Message)
		return
	}

	if err = json.NewDecoder(resp.Body).Decode(&signInResp); err != nil {
		return
	}

	return signInResp, err
}

type SignInResponse struct {
	Data struct {
		Tokens    Tokens    `json:"tokens"`
		Challenge Challenge `json:"challenge"`
	} `json:"data"`
}

type Tokens struct {
	AccessToken   string `json:"accessToken"`
	IdentityToken string `json:"identityToken"`
	RefreshToken  string `json:"refreshToken"`
}

type Challenge struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}
