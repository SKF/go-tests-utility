package auth

import (
	"context"
	"fmt"
	"sync"
	"github.com/SKF/go-utility/v2/auth"
	"github.com/SKF/go-utility/v2/auth/cachedauth"
)

const ssoBaseURL = "https://sso-api.%s.users.enlight.skf.com"
var lock sync.Mutex

// SignIn will sign in the user and if needed complete the change password challenge
func SignIn(stage, username, password string) (tokens Tokens, err error) {
	lock.Lock()
	defer lock.Unlock()

	cachedauth.Configure(cachedauth.Config{Stage: stage})
	err = cachedauth.SignIn(context.Background(), username, password)

	if err != nil {
		err = fmt.Errorf("failed to signin: %w",err)
	}

	return convertTokens(cachedauth.GetTokensByUser(username)), err
}

func convertTokens(in auth.Tokens) Tokens {
	return Tokens{
		AccessToken:   in.AccessToken,
		IdentityToken: in.IdentityToken,
		RefreshToken:  in.RefreshToken,
	}
}


type Tokens struct {
	AccessToken   string `json:"accessToken"`
	IdentityToken string `json:"identityToken"`
	RefreshToken  string `json:"refreshToken"`
}

