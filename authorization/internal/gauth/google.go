package gauth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/trennepohl/istio-auth-poc/authorization/internal"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleUserInfoUrl = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
)

type googleAuth struct {
	oauthConfig *oauth2.Config
	stateString string
}

func (g *googleAuth) GetUserInfo(code string, state string) (user internal.User, err error) {
	if state != g.stateString {
		return user, fmt.Errorf("invalid state %s", state)
	}

	googleToken, err := g.oauthConfig.Exchange(context.TODO(), code)
	if err != nil {
		return user, err
	}

	res, err := http.Get(googleUserInfoUrl + googleToken.AccessToken)
	if err != nil {
		return user, err
	}

	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(&user)
	return user, err
}

func (g *googleAuth) Login() string {
	return g.oauthConfig.AuthCodeURL(g.stateString)
}

func (g *googleAuth) GetRedirectUri() string {
	return g.oauthConfig.RedirectURL
}

func NewGoogleAuth(settings *internal.ServiceConfig) internal.AuthenticationService {
	return &googleAuth{
		oauthConfig: &oauth2.Config{
			RedirectURL:  settings.GoogleCallbackURL,
			ClientSecret: settings.GoogleClientSecret,
			ClientID:     settings.GoogleClientID,
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint:     google.Endpoint,
		},
		stateString: settings.GoogleStateCode,
	}
}
