package auth

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"slices"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	googleAuthURL          = "https://accounts.google.com/o/oauth2/v2/auth"
	googleTokenURL         = "https://oauth2.googleapis.com/token"
	portabilityScopePrefix = "https://www.googleapis.com/auth/dataportability."
	maxClientBytes         = 1 << 20
)

func normalizeScopes(input []string) ([]string, error) {
	scopes := slices.Clone(input)
	var portability, other bool
	for _, scope := range scopes {
		if len(scope) == 0 || len(scope) > 512 {
			return nil, errors.New("Google OAuth scopes must be nonempty scope strings")
		}
		for _, ch := range scope {
			if ch < 0x21 || ch > 0x7e || ch == '"' || ch == '\\' {
				return nil, errors.New("Google OAuth scopes cannot contain whitespace or invalid characters")
			}
		}
		if strings.HasPrefix(scope, portabilityScopePrefix) {
			portability = true
		} else {
			other = true
		}
	}
	if portability && other {
		return nil, errors.New("Data Portability OAuth scopes cannot be mixed with other scopes; configure separate grants")
	}
	slices.Sort(scopes)
	return slices.Compact(scopes), nil
}

func loadClient(path string, scopes []string) (*oauth2.Config, error) {
	info, err := os.Stat(path)
	if err != nil || !info.Mode().IsRegular() || info.Size() > maxClientBytes {
		return nil, errors.New("cannot read the Google OAuth client JSON file; configure an accessible Desktop app client file (at most 1 MiB)")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.New("cannot open the Google OAuth Desktop app client JSON file")
	}
	defer f.Close()
	data, err := io.ReadAll(io.LimitReader(f, maxClientBytes+1))
	if err != nil || len(data) > maxClientBytes {
		return nil, errors.New("cannot read the Google OAuth Desktop app client JSON file within the 1 MiB limit")
	}
	defer clear(data)
	var envelope map[string]json.RawMessage
	if json.Unmarshal(data, &envelope) != nil || len(envelope["installed"]) == 0 || envelope["web"] != nil || envelope["type"] != nil {
		return nil, errors.New("Google OAuth credentials must be a downloaded Desktop app client JSON with an installed section; web and service-account credentials are unsupported")
	}
	config, err := google.ConfigFromJSON(data, scopes...)
	if err != nil || config == nil || config.ClientID == "" || config.ClientSecret == "" || !strings.HasSuffix(config.ClientID, ".apps.googleusercontent.com") || strings.ContainsAny(config.ClientID, " \t\r\n") {
		return nil, errors.New("invalid Google OAuth Desktop app client JSON; download a new Desktop app client from Google Cloud")
	}
	if (config.Endpoint.AuthURL != googleAuthURL && config.Endpoint.AuthURL != "https://accounts.google.com/o/oauth2/auth") || config.Endpoint.TokenURL != googleTokenURL {
		return nil, errors.New("Google OAuth client JSON contains unsupported authorization/token endpoints; download an unmodified Desktop app client from Google Cloud")
	}
	config.Endpoint = oauth2.Endpoint{AuthURL: googleAuthURL, TokenURL: googleTokenURL, AuthStyle: oauth2.AuthStyleInParams}
	// The downloaded redirect list is never a network destination. Login always
	// replaces it with its own ephemeral IPv4 loopback callback.
	config.RedirectURL = ""
	return config, nil
}

func validateGrant(token *oauth2.Token, required []string) error {
	value := token.Extra("scope")
	if value == nil {
		// OAuth permits omitted scope when it is identical to the request.
		return nil
	}
	text, ok := value.(string)
	if !ok {
		return ErrAuthorizationRequired
	}
	granted := strings.Fields(text)
	if _, err := normalizeScopes(granted); err != nil {
		return ErrAuthorizationRequired
	}
	for _, scope := range required {
		if !slices.Contains(granted, scope) {
			return ErrAuthorizationRequired
		}
	}
	return nil
}
