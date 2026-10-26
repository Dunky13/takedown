package auth

import (
	"context"
	"errors"
	"time"
)

// AuthorizationStatus describes local setup only. It never claims that Google
// still accepts a stored grant, or that the login hint is a verified identity.
type AuthorizationStatus struct {
	Configured           bool      `json:"configured"`
	TokenStored          bool      `json:"token_stored"`
	RefreshAvailable     bool      `json:"refresh_available"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at,omitempty"`
	Message              string    `json:"message"`
}

// InspectAuthorization reads configuration and the keyring without refreshing
// a token, opening a browser, writing secrets, or making a network request.
func (m *Manager) InspectAuthorization(ctx context.Context) (AuthorizationStatus, error) {
	_, release, err := m.acquire(ctx)
	if err != nil {
		return AuthorizationStatus{}, err
	}
	defer release()
	status := AuthorizationStatus{}
	if err = m.ensureConfig(); err != nil {
		return status, err
	}
	status.Configured = true
	token, err := m.vault.load()
	if errors.Is(err, ErrAuthorizationRequired) {
		status.Message = "No usable saved authorization; run auth or Authorize in browser."
		return status, nil
	}
	if err != nil {
		return status, err
	}
	status.TokenStored = true
	status.RefreshAvailable = token.RefreshToken != ""
	status.AccessTokenExpiresAt = token.Expiry
	status.Message = "Authorization is stored in the OS keyring; Google grant validity and account identity have not been checked."
	if !token.Valid() {
		status.Message += " The access token has expired; a valid refresh grant is required for background work."
	}
	return status, nil
}
