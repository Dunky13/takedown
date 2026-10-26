package storage

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"io"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

func testSigner(t *testing.T) ssh.Signer {
	t.Helper()
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	signer, err := ssh.NewSignerFromKey(key)
	if err != nil {
		t.Fatal(err)
	}
	return signer
}

func sftpServer(t *testing.T) (address, knownPath string) {
	t.Helper()
	signer := testSigner(t)
	config := &ssh.ServerConfig{PasswordCallback: func(_ ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
		if string(password) != "test-only-password" {
			return nil, errors.New("wrong test password")
		}
		return nil, nil
	}}
	config.AddHostKey(signer)
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	knownPath = filepath.Join(t.TempDir(), "known_hosts")
	if err := os.WriteFile(knownPath, []byte(knownhosts.Line([]string{knownhosts.Normalize(listener.Addr().String())}, signer.PublicKey())+"\n"), 0600); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	var mu sync.Mutex
	var conns []net.Conn
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			mu.Lock()
			conns = append(conns, conn)
			mu.Unlock()
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer conn.Close()
				server, channels, requests, err := ssh.NewServerConn(conn, config)
				if err != nil {
					return
				}
				defer server.Close()
				go ssh.DiscardRequests(requests)
				for incoming := range channels {
					if incoming.ChannelType() != "session" {
						_ = incoming.Reject(ssh.UnknownChannelType, "session required")
						continue
					}
					channel, requests, err := incoming.Accept()
					if err != nil {
						continue
					}
					for request := range requests {
						var subsystem struct{ Name string }
						if request.Type != "subsystem" || ssh.Unmarshal(request.Payload, &subsystem) != nil || subsystem.Name != "sftp" {
							_ = request.Reply(false, nil)
							continue
						}
						_ = request.Reply(true, nil)
						service, err := sftp.NewServer(channel)
						if err == nil {
							_ = service.Serve()
							_ = service.Close()
						}
						_ = channel.Close()
						break
					}
				}
			}()
		}
	}()
	t.Cleanup(func() {
		_ = listener.Close()
		<-done
		mu.Lock()
		for _, conn := range conns {
			_ = conn.Close()
		}
		mu.Unlock()
		wg.Wait()
	})
	return listener.Addr().String(), knownPath
}

func TestSFTPRealConnectionResumeIntegrityAndConfinement(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("in-process SFTP fixture uses a Unix filesystem namespace")
	}
	address, knownPath := sftpServer(t)
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	destination := (&url.URL{Scheme: "sftp", Host: address, User: url.User("test"), Path: filepath.ToSlash(root) + "/exports"}).String()
	backend, err := Open(context.Background(), destination, Options{Password: "test-only-password", KnownHostsPath: knownPath})
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := backend.Close(); err != nil {
			t.Error(err)
		}
	}()
	reader, ok := backend.(ReadBackend)
	if !ok {
		t.Fatal("SFTP missing integrity reader")
	}
	writeArchive(t, backend, "batch/part.zip.part", 0, "prefix-stale-data")
	writeArchive(t, backend, "batch/part.zip.part", 6, "verified")
	if got := readArchive(t, reader, "batch/part.zip.part"); got != "prefixverified" {
		t.Fatalf("bad resume: %q", got)
	}
	if err := backend.Rename(context.Background(), "batch/part.zip.part", "batch/part.zip"); err != nil {
		t.Fatal(err)
	}
	r, err := reader.OpenRead(context.Background(), "batch/part.zip")
	if err != nil {
		t.Fatal(err)
	}
	buf := make([]byte, 8)
	if n, err := r.ReadAt(buf, 6); n != 8 || err != nil || string(buf) != "verified" {
		t.Fatalf("remote ReadAt: %d %q %v", n, buf, err)
	}
	if _, err := r.Seek(6, io.SeekStart); err != nil {
		t.Fatal(err)
	}
	if data, err := io.ReadAll(r); err != nil || string(data) != "verified" {
		t.Fatalf("remote Seek/Read: %q %v", data, err)
	}
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	info, err := backend.Stat(context.Background(), "batch/part.zip")
	if err != nil || info.Size() != 14 {
		t.Fatalf("remote Stat: %v %v", info, err)
	}
	if _, err := backend.FreeSpace(context.Background()); err != nil && !errors.Is(err, ErrSpaceUnknown) {
		t.Fatal(err)
	}
	outside := filepath.Join(root, "outside.zip")
	if err := os.WriteFile(outside, []byte("safe"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "exports", "link.zip")); err == nil {
		if w, err := backend.OpenWriteResume(context.Background(), "link.zip", 0); !errors.Is(err, ErrUnsafePath) {
			if w != nil {
				w.Close()
			}
			t.Fatalf("remote symlink accepted: %v", err)
		}
	}
	if _, err := backend.Stat(context.Background(), "../outside.zip"); !errors.Is(err, ErrUnsafePath) {
		t.Fatalf("remote traversal accepted: %v", err)
	}
	if data, err := os.ReadFile(outside); err != nil || string(data) != "safe" {
		t.Fatalf("outside file changed: %q %v", data, err)
	}
}

func TestSFTPRejectsUnknownAndChangedHostKeys(t *testing.T) {
	address, _ := sftpServer(t)
	destination := (&url.URL{Scheme: "sftp", Host: address, User: url.User("test"), Path: "/exports"}).String()
	for _, changed := range []bool{false, true} {
		knownPath := filepath.Join(t.TempDir(), "known_hosts")
		var content string
		if changed {
			content = knownhosts.Line([]string{knownhosts.Normalize(address)}, testSigner(t).PublicKey()) + "\n"
		}
		if err := os.WriteFile(knownPath, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
		_, err := Open(context.Background(), destination, Options{Password: "test-only-password", KnownHostsPath: knownPath})
		var keyError *knownhosts.KeyError
		if !errors.As(err, &keyError) {
			t.Fatalf("changed=%t: expected host-key rejection, got %v", changed, err)
		}
	}
}

func TestSFTPCancellationInterruptsHandshake(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	accepted := make(chan net.Conn, 1)
	go func() {
		conn, err := listener.Accept()
		if err == nil {
			accepted <- conn
		}
	}()
	knownPath := filepath.Join(t.TempDir(), "known_hosts")
	if err := os.WriteFile(knownPath, nil, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	result := make(chan error, 1)
	go func() {
		_, err := Open(ctx, "sftp://test@"+listener.Addr().String()+"/exports", Options{Password: "test-only-password", KnownHostsPath: knownPath})
		result <- err
	}()
	select {
	case conn := <-accepted:
		defer conn.Close()
	case <-time.After(3 * time.Second):
		t.Fatal("SSH client did not connect")
	}
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancel: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("cancellation did not interrupt handshake")
	}
}
