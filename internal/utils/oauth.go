package utils

import (
	"context"
	"fmt"

	"github.com/coreos/go-oidc"
	"github.com/syed.fazil/vtask/internal/config"
	"golang.org/x/oauth2"
)

func NewOauth2Config(ctx context.Context) (oauth2Config *oauth2.Config, err error) {
	provider, err := oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("could not create new provider: %v", err)
	}
	oauth2Config = &oauth2.Config{
		ClientID:     config.App.OAuthGoogleClientID,
		ClientSecret: config.App.OAuthGoogleClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  fmt.Sprintf("%s/auth/google/callback", config.App.AppBaseURL),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
	return
}

func NewOIDCProvider(ctx context.Context) (oidcProvider *oidc.Provider, err error) {
	oidcProvider, err = oidc.NewProvider(ctx, "https://accounts.google.com")
	if err != nil {
		return nil, fmt.Errorf("could not create new provider: %v", err)
	}
	return
}
