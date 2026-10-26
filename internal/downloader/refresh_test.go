package downloader

import (
	"context"
	"testing"
)

type rejectionAwareLinks struct {
	linkFunc
	rejected string
	calls    int
}

func (s *rejectionAwareLinks) RefreshRejectedLink(_ context.Context, _ string, _ int, rejected string) (string, error) {
	s.rejected = rejected
	s.calls++
	return "https://example.test/renewed.zip?signature=new", nil
}

func TestRefreshSuppliesRejectedSignatureAndReusesRenewedLink(t *testing.T) {
	links := &rejectionAwareLinks{linkFunc: func(context.Context, string, int) (string, error) {
		t.Fatal("rejection-aware source used the legacy refresh contract")
		return "", nil
	}}
	f := newFixture(t, zipBytes(t, 100), "archive.zip", "", func(o *Options) { o.Links = links })
	const expired = "https://example.test/archive.zip?signature=expired"
	fresh, err := f.engine.refresh(context.Background(), f.part, expired)
	if err != nil || fresh == expired || links.rejected != expired {
		t.Fatalf("rejected signature lost: %q, %q, %v", fresh, links.rejected, err)
	}
	again, err := f.engine.refresh(context.Background(), f.part, expired)
	if err != nil || again != fresh || links.calls != 1 {
		t.Fatalf("concurrent old-signature caller refreshed unnecessarily: %v, calls=%d", err, links.calls)
	}
}
