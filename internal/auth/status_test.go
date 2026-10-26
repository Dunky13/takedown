package auth

import (
	"context"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestInspectAuthorizationNeverRefreshesOrLeaksCredentials(t *testing.T) {
	store := &fakeSecrets{}
	m := testManager(t, testOptions(t), store, func(*http.Request) (*http.Response, error) { t.Fatal("local status contacted Google"); return nil, nil })
	token := validToken()
	token.Expiry = time.Now().Add(-time.Hour)
	seedToken(t, m, token)
	status, err := m.InspectAuthorization(context.Background())
	if err != nil || !status.Configured || !status.TokenStored || !status.RefreshAvailable || !status.AccessTokenExpiresAt.Equal(token.Expiry) {
		t.Fatalf("status: %+v, %v", status, err)
	}
	if !strings.Contains(status.Message, "have not been checked") || !strings.Contains(status.Message, "expired") {
		t.Fatalf("status overstates authorization: %+v", status)
	}
	if strings.Contains(status.Message, token.AccessToken) || strings.Contains(status.Message, token.RefreshToken) {
		t.Fatal("status exposed a token")
	}
}
