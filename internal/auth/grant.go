package auth

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

// The library exposes Token.Extra("scope") but represents explicit JSON null
// the same as an omitted scope. Check grant metadata at the HTTP boundary so a
// malformed/empty explicitly supplied scope cannot bypass the grant check.
// Exchange, PKCE, token parsing and refresh remain entirely library-owned.
type grantTransport struct {
	base   http.RoundTripper
	scopes []string
}

func (t *grantTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	response, err := t.base.RoundTrip(request)
	if err != nil || response.StatusCode < 200 || response.StatusCode > 299 {
		return response, err
	}
	data, err := io.ReadAll(io.LimitReader(response.Body, (1<<20)+1))
	_ = response.Body.Close()
	if err != nil || len(data) > 1<<20 {
		clear(data)
		if temporary := temporaryNetwork(err); temporary != nil {
			return nil, temporary
		}
		return nil, errors.New("cannot read the bounded Google OAuth token response")
	}
	response.Body = &tokenResponseBody{Reader: bytes.NewReader(data), data: data}
	var scope string
	var present, valid bool
	media, _, _ := mime.ParseMediaType(response.Header.Get("Content-Type"))
	if media == "application/x-www-form-urlencoded" || media == "text/plain" {
		values, parseErr := url.ParseQuery(string(data))
		present = values.Has("scope")
		scope = values.Get("scope")
		valid = parseErr == nil && len(values["scope"]) == 1
	} else {
		var metadata struct {
			Scope json.RawMessage `json:"scope"`
		}
		if json.Unmarshal(data, &metadata) == nil && len(metadata.Scope) > 0 {
			present = true
			valid = json.Unmarshal(metadata.Scope, &scope) == nil
		}
	}
	if present && (!valid || validateGrant((&oauth2.Token{}).WithExtra(map[string]any{"scope": scope}), t.scopes) != nil) {
		_ = response.Body.Close()
		return nil, ErrAuthorizationRequired
	}
	return response, nil
}

type tokenResponseBody struct {
	*bytes.Reader
	data []byte
}

func (b *tokenResponseBody) Close() error { clear(b.data); b.data = nil; return nil }
